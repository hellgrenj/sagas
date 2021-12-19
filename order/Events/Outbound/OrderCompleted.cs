using MediatR;
using order.Infrastructure;
using order.Model;

namespace order.Events.Outbound;

public class OrderCompletedEvent : Event, INotification
{
    public OrderCompletedEvent()
    {
        this.Name = "order completed";
        this.MessageId = Guid.NewGuid().ToString();
    }
    public int OrderId { get; set; }
}
public class OrderCompletedEventHandler : INotificationHandler<OrderCompletedEvent>
{
    private readonly IRabbitEmitter _rabbitEmitter;
    private readonly ILogger<OrderCompletedEventHandler> _logger;
    public OrderCompletedEventHandler(IRabbitEmitter rabbitEmitter, ILogger<OrderCompletedEventHandler> logger)
    {
        _rabbitEmitter = rabbitEmitter;
        _logger = logger;
    }

    public Task Handle(OrderCompletedEvent e, CancellationToken cancellationToken)
    {
        _rabbitEmitter.PublishEvent(e);
        _logger.LogInformation($"OrderCompletedEvent published with correlationId: {e.CorrelationId}");
        return Task.CompletedTask;
    }
}