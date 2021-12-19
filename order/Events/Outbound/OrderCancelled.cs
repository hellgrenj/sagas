using MediatR;
using order.Infrastructure;
using order.Model;

namespace order.Events.Outbound;

public class OrderCancelledEvent : Event, INotification
{
    public OrderCancelledEvent()
    {
        this.Name = "order cancelled";
        this.MessageId = Guid.NewGuid().ToString();
    }

    public string Reason { get; set; }
    public Order Order { get; set; }

}
public class OrderCancelledEventHandler : INotificationHandler<OrderCancelledEvent>
{
    private readonly IRabbitEmitter _rabbitEmitter;
    private readonly ILogger<OrderCancelledEventHandler> _logger;
    public OrderCancelledEventHandler(IRabbitEmitter rabbitEmitter, ILogger<OrderCancelledEventHandler> logger)
    {
        _rabbitEmitter = rabbitEmitter;
        _logger = logger;
    }

    public Task Handle(OrderCancelledEvent e, CancellationToken cancellationToken)
    {
        _rabbitEmitter.PublishEvent(e);
        _logger.LogInformation($"OrderCancelledEvent published with correlationId: {e.CorrelationId}");
        return Task.CompletedTask;
    }
}