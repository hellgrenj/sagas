using warehouse.Model;

namespace warehouse.Events.Inbound;

public class OrderCreatedEvent : Event
{
    public Order Order { get; set; }
}
