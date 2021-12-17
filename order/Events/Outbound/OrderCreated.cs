using MediatR;
using order.Infrastructure;
using order.Model;

namespace order.Events.Outbound;

public class OrderCreatedEvent : Event, INotification
{
    public OrderCreatedEvent()
    {
        this.Name = "order created";
        this.MessageId = Guid.NewGuid().ToString();
    }
    public Order Order { get; set; }
}
public class OrderCreatedEventHandler : INotificationHandler<OrderCreatedEvent>
{
    private readonly IRabbitEmitter _rabbitEmitter;
    private readonly ILogger<OrderCreatedEventHandler> _logger;
    public OrderCreatedEventHandler(IRabbitEmitter rabbitEmitter, ILogger<OrderCreatedEventHandler> logger)
    {
        _rabbitEmitter = rabbitEmitter;
        _logger = logger;
    }

    public Task Handle(OrderCreatedEvent e, CancellationToken cancellationToken)
    {
        _rabbitEmitter.PublishEvent(e);
        _logger.LogInformation($"OrderCreatedEvent published with correlationId: {e.CorrelationId}");
        return Task.CompletedTask;
    }
}