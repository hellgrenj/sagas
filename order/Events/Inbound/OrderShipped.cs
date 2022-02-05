using order.DTOs;

namespace order.Events.Inbound;

public class OrderShippedevent : Event
{
    public int OrderId { get; set; }
}