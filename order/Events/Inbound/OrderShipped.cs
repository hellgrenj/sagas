using order.Model;

namespace order.Events.Inbound;

public class OrderShippedevent : Event
{
    public int OrderId { get; set; }
}