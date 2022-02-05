using order.DTOs;

namespace order.Events.Inbound;

public class ItemsNotInStockEvent : Event
{
    public Order Order { get; set; }
}