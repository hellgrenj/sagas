using order.Model;

namespace order.Events.Inbound;

public class ItemsNotInStockEvent : Event
{
    public Order Order { get; set; }
}