using System.Data;
using Dapper;
using model.Repositories;
using order.Model;

namespace order.Repositories;

public class OrderRepository : IRepository<Order>
{
    private readonly IDbConnection _connection;

    public OrderRepository(IDbConnection connection)
    {
        _connection = connection;
    }
    public async Task<Order> GetAsync(int id)
    {
        return await _connection.QueryFirstAsync<Order>("SELECT * FROM ordertable WHERE id = @id", new { id = id });
    }

    // TODO update audit table (use JSON support in postgres for this!)
    public async Task<int> SaveAsync(Order order)
    {
        if (order.Id.HasValue)
        {
            var sql = "UPDATE ordertable SET item=@item, quantity=@quantity, price=@price, state=@state WHERE id=@id RETURNING Id;";
            return (int)await _connection.ExecuteScalarAsync(sql,
            new { item = order.Item, quantity = order.Quantity, price = order.Price, state = order.State, id = order.Id.Value });
        }
        else
        {
            var sql = "INSERT INTO ordertable (correlationid, item, quantity, price, state) Values (@correlationid, @item, @quantity, @price,  @state) RETURNING Id;";
            return (int)await _connection.ExecuteScalarAsync(sql,
            new { correlationid = order.CorrelationId, item = order.Item, quantity = order.Quantity, price = order.Price, state = order.State });
        }
    }
}