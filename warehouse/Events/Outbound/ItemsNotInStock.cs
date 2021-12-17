using MediatR;
using warehouse.Infrastructure;
using warehouse.Model;

namespace warehouse.Events.Outbound;

public class ItemsNotInStockEvent : Event, INotification
{
    public ItemsNotInStockEvent()
    {
        this.Name = "items not in stock";
        this.MessageId = Guid.NewGuid().ToString();
    }
    public Order Order { get; set; }
}


public class ItemsNotInStockEventHandler : INotificationHandler<ItemsNotInStockEvent>
{
    private readonly IRabbitEmitter _rabbitEmitter;
    private readonly ILogger<ItemsNotInStockEventHandler> _logger;
    public ItemsNotInStockEventHandler(IRabbitEmitter rabbitEmitter, ILogger<ItemsNotInStockEventHandler> logger)
    {
        _rabbitEmitter = rabbitEmitter;
        _logger = logger;
    }

    public Task Handle(ItemsNotInStockEvent e, CancellationToken cancellationToken)
    {
        _rabbitEmitter.PublishEvent(e);
        _logger.LogInformation($"ItemsNotInStockEvent published with correlationId: {e.CorrelationId}");
        return Task.CompletedTask;
    }
}