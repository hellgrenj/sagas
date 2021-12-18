using MediatR;
using warehouse.Infrastructure;
using warehouse.Inventory.Model;

namespace warehouse.Events.Outbound;

public class OrderShippedEvent : Event, INotification
{
    public OrderShippedEvent()
    {
        this.Name = "order shipped";
        this.MessageId = Guid.NewGuid().ToString();
    }
    public int OrderId { get; set; }
}


public class OrderShippedEventHandler : INotificationHandler<OrderShippedEvent>
{
    private readonly IRabbitEmitter _rabbitEmitter;
    private readonly ILogger<OrderShippedEventHandler> _logger;
    public OrderShippedEventHandler(IRabbitEmitter rabbitEmitter, ILogger<OrderShippedEventHandler> logger)
    {
        _rabbitEmitter = rabbitEmitter;
        _logger = logger;
    }

    public Task Handle(OrderShippedEvent e, CancellationToken cancellationToken)
    {
        _rabbitEmitter.PublishEvent(e);
        _logger.LogInformation($"OrderShippedEvent published with correlationId: {e.CorrelationId}");
        return Task.CompletedTask;
    }
}