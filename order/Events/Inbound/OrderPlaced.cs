using order.DTOs;

namespace order.Events.Inbound;

public class OrderPlacedEvent : Event
{
    public Order Order { get; set; }
}