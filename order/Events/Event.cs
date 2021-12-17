namespace order.Events;
public class Event
{
    public string CorrelationId { get; set; }
    public string Name { get; set; }
    public string MessageId { get; set; }
}
