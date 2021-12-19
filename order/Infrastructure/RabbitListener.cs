using System.Text;
using System.Text.Json;
using MediatR;
using order.Events;
using order.Events.Inbound;
using order.Features;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;

namespace order.Infrastructure;

public class RabbitListener : BackgroundService
{
    private IConnection connection;
    private IModel channel; //one channel for the ONE background process
    private readonly ILogger<RabbitListener> _logger;

    private readonly IRabbitConnectionHandler _rabbitConnetionHandler;

    public IServiceScopeFactory _serviceScopeFactory;
    public RabbitListener(ILogger<RabbitListener> logger, IServiceScopeFactory serviceScopeFactory, IRabbitConnectionHandler rabbitConnectionHandler)
    {
        _serviceScopeFactory = serviceScopeFactory;
        _logger = logger;
        _rabbitConnetionHandler = rabbitConnectionHandler;
    }
    protected override Task ExecuteAsync(CancellationToken stoppingToken)
    {
        StartListening();
        return Task.CompletedTask;
    }
    private void StartListening()
    {
        _logger.LogInformation("connecting to rabbitmq");
        var factory = new ConnectionFactory()
        {
            HostName = "rabbit"
        };
        factory.AutomaticRecoveryEnabled = true;
        factory.NetworkRecoveryInterval = TimeSpan.FromSeconds(10);
        connection = _rabbitConnetionHandler.TryOpen(factory, _logger);
        channel = connection.CreateModel();

        var exchange = "order.topics";
        channel.ExchangeDeclare(exchange, type: "topic");

        var queueName = channel.QueueDeclare("orderservice_queue", durable: true, exclusive: false,
                             autoDelete: false, arguments: null).QueueName;
        channel.BasicQos(prefetchSize: 0, prefetchCount: 1, global: false);

        channel.QueueBind(queue: queueName,
                  exchange,
                  routingKey: "order.placed");
        channel.QueueBind(queue: queueName,
                exchange,
                routingKey: "order.items.notinstock");
        channel.QueueBind(queue: queueName,
                exchange,
                routingKey: "order.shipped");

        _logger.LogInformation("waiting for messages");
        var consumer = new EventingBasicConsumer(channel);
        consumer.Received += async (model, ea) =>
        {
            var body = ea.Body.ToArray();
            var json = Encoding.UTF8.GetString(body);
            var routingKey = ea.RoutingKey;
            try
            {
                await ProcessMessage(json, routingKey);
            }
            catch (Exception e)
            {
                _logger.LogError($"Failed to process message with exception {e.Message} stacktrace {e.StackTrace}");
            }
            finally
            {
                channel.BasicAck(deliveryTag: ea.DeliveryTag, multiple: false);
            }
        };

        channel.BasicConsume(queue: queueName,
                             autoAck: false,
                             consumer: consumer);
    }

    private async Task ProcessMessage(string json, string routingKey)
    {
        _logger.LogInformation($"received message {json}");
        var message = JsonSerializer.Deserialize<Event>(json, new JsonSerializerOptions() { PropertyNamingPolicy = JsonNamingPolicy.CamelCase });
        if (message is null || message.MessageId is null)
            throw new InvalidOperationException("a message with a messageId must be provided");

        var alreadyProcessed = await TryMarkMessageAsProcessed(message.MessageId);
        if (alreadyProcessed)
            return;

        using var scope = _serviceScopeFactory.CreateScope();
        var mediator = scope.ServiceProvider.GetService<IMediator>();
        switch (routingKey)
        {
            case "order.placed":
                var orderPlacedEvent = JsonSerializer.Deserialize<OrderPlacedEvent>(json, new JsonSerializerOptions() { PropertyNamingPolicy = JsonNamingPolicy.CamelCase });
                await mediator.Send(new CreateOrderCommand(orderPlacedEvent));
                break;
            case "order.items.notinstock":
                var itemsNotInStockEvent = JsonSerializer.Deserialize<ItemsNotInStockEvent>(json, new JsonSerializerOptions() { PropertyNamingPolicy = JsonNamingPolicy.CamelCase });
                await mediator.Send(new CancelOrderCommand(itemsNotInStockEvent));
                break;
            case "order.shipped":
                var orderShippedevent = JsonSerializer.Deserialize<OrderShippedevent>(json, new JsonSerializerOptions() { PropertyNamingPolicy = JsonNamingPolicy.CamelCase });
                await mediator.Send(new CompleteOrderCommand(orderShippedevent));
                break;
            default:
                break;
        }
    }
    private async Task<bool> TryMarkMessageAsProcessed(string messageId)
    {
        using var scope = _serviceScopeFactory.CreateScope();
        var mediator = scope.ServiceProvider.GetService<IMediator>();
        return await mediator.Send(new TryMarkMessageAsProcessedCommand(messageId));
    }
    private void StopListening()
    {
        _rabbitConnetionHandler.Close(_logger, connection, channel);
    }

    public override void Dispose()
    {
        StopListening();
        base.Dispose();
    }
}

