using System.Text;
using System.Text.Json;
using warehouse.Events;
using RabbitMQ.Client;
using warehouse.Events.Inbound;
using warehouse.Events.Outbound;

namespace warehouse.Infrastructure;

public interface IRabbitEmitter
{
    void PublishEvent(Event e);
    void Open();
    void Close();
}

public class RabbitEmitter : IRabbitEmitter
{

    private IConnection connection; // per app (thread safe..)
    private readonly ILogger<RabbitEmitter> _logger;
    private readonly IRabbitConnectionHandler _rabbitConnectionHandler;
    public RabbitEmitter(ILogger<RabbitEmitter> logger, IRabbitConnectionHandler rabbitConnectionHandler)
    {
        _logger = logger;
        _rabbitConnectionHandler = rabbitConnectionHandler;
    }

    public void PublishEvent(Event e)
    {
        var exchange = "order.topics";
        _logger.LogInformation($"publishing event: {e} on exchange {exchange}");
        using var channel = connection.CreateModel(); // per call
        channel.ExchangeDeclare(exchange, type: "topic");
        var properties = channel.CreateBasicProperties();
        properties.Persistent = true;

        var (routingKey, message) = e switch
        {
            ItemsReservedEvent _ => ("order.items.reserved", Encoding.UTF8.GetBytes(JsonSerializer.Serialize((ItemsReservedEvent)e,
            new JsonSerializerOptions() { PropertyNamingPolicy = JsonNamingPolicy.CamelCase }))),
            ItemsNotInStockEvent _ => ("order.items.notinstock", Encoding.UTF8.GetBytes(JsonSerializer.Serialize((ItemsNotInStockEvent)e,
            new JsonSerializerOptions() { PropertyNamingPolicy = JsonNamingPolicy.CamelCase }))),
            OrderShippedEvent _ => ("order.shipped", Encoding.UTF8.GetBytes(JsonSerializer.Serialize((OrderShippedEvent)e,
           new JsonSerializerOptions() { PropertyNamingPolicy = JsonNamingPolicy.CamelCase }))),

            _ => throw new Exception("unknown event")
        };

        channel.BasicPublish(exchange,
                                routingKey,
                                basicProperties: null,
                                body: message);
    }
    public void Open()
    {
        _logger.LogInformation("connecting to rabbitmq");
        var factory = new ConnectionFactory() { HostName = "rabbit" };
        factory.AutomaticRecoveryEnabled = true;
        factory.NetworkRecoveryInterval = TimeSpan.FromSeconds(10);
        connection = _rabbitConnectionHandler.TryOpen(factory, _logger);
        _logger.LogInformation("connected to rabbitmq");
    }
    public void Close()
    {
        _rabbitConnectionHandler.Close(_logger, connection);
    }

}