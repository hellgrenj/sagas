namespace order.Model;
public class Order
{
    public int? Id { get; set; }
    public string Item { get; set; }
    public int Price { get; set; }
    public int Quantity { get; set; }
    public string State { get; private set; }

    public string CorrelationId { get; set; }
    public Order() { }
    public Order(DTOs.Order order, string correlationId)
    {
        Id = order.Id.HasValue ? order.Id.Value : null;
        Item = order.Item;
        Price = order.Price;
        Quantity = order.Quantity;
        CorrelationId = correlationId;
    }
    public void ChangeState(string newState)
    {
        switch (newState)
        {
            case OrderStates.Pending:
                if (State == OrderStates.Completed || State == OrderStates.Cancelled)
                {
                    throw new InvalidOperationException($"Order in state {State} cannot be changed to Pending");
                }
                State = OrderStates.Pending;
                break;
            case OrderStates.Cancelled:
                if (State == OrderStates.Completed)
                {
                    throw new InvalidOperationException($"Order in state {State} cannot be changed to Cancelled");
                }
                State = OrderStates.Cancelled;
                break;
            case OrderStates.Completed:
                if (State == OrderStates.Cancelled)
                {
                    throw new InvalidOperationException($"Order in state {State} cannot be changed to Completed");
                }
                State = OrderStates.Completed;
                break;
            default:
                throw new InvalidOperationException($"{newState} is not a valid state");
        }
    }
}
public class OrderStates
{
    public const string Pending = "Pending";
    public const string Cancelled = "Cancelled";
    public const string Completed = "Completed";
}
