namespace order.Model;
public class Order
{
    public int Id { get; set; }
    public string Item { get; set; }
    public int Price { get; set; }
    public int Quantity { get; set; }

    public string State { get; set; }
}
public class OrderStates
{
    public const string PaymentPending = "PaymentPending";
    public const string Cancelled = "Cancelled";
    public const string PaymentConfirmed = "PaymentConfirmed";
    public const string Deliverd = "Delivered";
}
