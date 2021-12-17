namespace warehouse.Inventory.Model;

public class Reservation
{
    public int Id { get; set; }
    public string CorrelationId { get; set; }

    public int OrderId { get; set; }
    public string Item { get; set; }
    public int Price { get; set; }
    public int Quantity { get; set; }
}