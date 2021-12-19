using System.Text;
using System.Text.Json;
using order.Events;
using order.Events.Outbound;
using RabbitMQ.Client;

namespace order.Infrastructure;

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
            OrderCreatedEvent _ => ("order.created", Encoding.UTF8.GetBytes(JsonSerializer.Serialize((OrderCreatedEvent)e,
            new JsonSerializerOptions() { PropertyNamingPolicy = JsonNamingPolicy.CamelCase }))),
            OrderCancelledEvent _ => ("order.cancelled", Encoding.UTF8.GetBytes(JsonSerializer.Serialize((OrderCancelledEvent)e,
            new JsonSerializerOptions() { PropertyNamingPolicy = JsonNamingPolicy.CamelCase }))),
            OrderCompletedEvent _ => ("order.completed", Encoding.UTF8.GetBytes(JsonSerializer.Serialize((OrderCompletedEvent)e,
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