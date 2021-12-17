using MediatR;
using warehouse.Infrastructure;
using warehouse.Inventory.Model;

namespace warehouse.Events.Outbound;

public class ItemsReservedEvent : Event, INotification
{
    public ItemsReservedEvent()
    {
        this.Name = "items reserved";
        this.MessageId = Guid.NewGuid().ToString();
    }
    public Reservation Reservation { get; set; }
}


public class ItemsReservedEventHandler : INotificationHandler<ItemsReservedEvent>
{
    private readonly IRabbitEmitter _rabbitEmitter;
    private readonly ILogger<ItemsReservedEventHandler> _logger;
    public ItemsReservedEventHandler(IRabbitEmitter rabbitEmitter, ILogger<ItemsReservedEventHandler> logger)
    {
        _rabbitEmitter = rabbitEmitter;
        _logger = logger;
    }

    public Task Handle(ItemsReservedEvent e, CancellationToken cancellationToken)
    {
        _rabbitEmitter.PublishEvent(e);
        _logger.LogInformation($"ItemsReservedEvent published with correlationId: {e.CorrelationId}");
        return Task.CompletedTask;
    }
}