using warehouse.Model;

namespace warehouse.Events.Inbound;

public class PaymentCompleted : Event
{
    public int OrderId { get; set; }
}
