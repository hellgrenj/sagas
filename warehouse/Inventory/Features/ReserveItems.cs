using System.Data;
using Dapper;
using FluentValidation;
using MediatR;
using warehouse.Events.Inbound;
using warehouse.Events.Outbound;
using warehouse.Inventory.Model;

namespace warehouse.Inventory.Features;

public record ReserveItemsCommand(OrderCreatedEvent OrderCreatedEvent) : IRequest<Unit>;

public class ReserveItemsCommandValidator : AbstractValidator<ReserveItemsCommand>
{
    public ReserveItemsCommandValidator()
    {
        RuleFor(x => x.OrderCreatedEvent.Order.Item).NotEmpty();
        RuleFor(x => x.OrderCreatedEvent.Order.Quantity).GreaterThan(0);
        RuleFor(x => x.OrderCreatedEvent.Order.Price).GreaterThan(0);
    }
}

public class ReserveItemsHandler : IRequestHandler<ReserveItemsCommand, Unit>
{
    private static Random random = new();
    private readonly ILogger<ReserveItemsHandler> _logger;
    private readonly ReserveItemsCommandValidator _validator;
    private readonly IDbConnection _connection;

    private readonly IMediator _mediator;
    public ReserveItemsHandler(ILogger<ReserveItemsHandler> logger, ReserveItemsCommandValidator validator, IDbConnection connection, IMediator mediator)
    {
        _logger = logger;
        _validator = validator;
        _connection = connection;
        _mediator = mediator;
    }

    public async Task<Unit> Handle(ReserveItemsCommand cmd, CancellationToken cancellationToken)
    {
        _logger.LogInformation($"Reserving {cmd.OrderCreatedEvent.Order.Quantity} {cmd.OrderCreatedEvent.Order.Item}");



        var validationResult = _validator.Validate(cmd);
        if (validationResult.IsValid)
        {
            var reservation = new Reservation()
            {
                CorrelationId = cmd.OrderCreatedEvent.CorrelationId,
                OrderId = cmd.OrderCreatedEvent.Order.Id,
                Item = cmd.OrderCreatedEvent.Order.Item,
                Price = cmd.OrderCreatedEvent.Order.Price,
                Quantity = cmd.OrderCreatedEvent.Order.Quantity
            };
            if (ItemsInStock(reservation))
            {
                var sql = "INSERT INTO reservations (correlationid, orderid, item, quantity, price) Values (@correlationid,@orderid, @item, @quantity, @price) RETURNING Id;";
                reservation.Id = (int)await _connection.ExecuteScalarAsync(sql,
                new { correlationid = reservation.CorrelationId, orderid = reservation.OrderId, item = reservation.Item, quantity = reservation.Quantity, price = reservation.Price });
                _logger.LogInformation($"Reservation created with id {reservation.Id}");
                await _mediator.Publish(new ItemsReservedEvent { CorrelationId = reservation.CorrelationId, Reservation = reservation });
            }
            else
            {
                await _mediator.Publish(new ItemsNotInStockEvent { CorrelationId = reservation.CorrelationId, Order = cmd.OrderCreatedEvent.Order });
            }

        }
        else
        {
            _logger.LogError($"ReserveItems validation failed \n {validationResult.ToString()}");
        }

        return Unit.Value;

    }
    private bool ItemsInStock(Reservation reservation)
    {
        var percentage = 50;
        _logger.LogInformation($"{percentage}% chance the items are in stock in this simulation");
        var randomNumber = random.Next(100);
        if (randomNumber < percentage)
        {
            _logger.LogInformation($"Items are in stock, proceeding with the reservation of {reservation.Quantity} {reservation.Item}");
            return true;
        }
        else
        {
            _logger.LogInformation("Items are not in stock");
            return false;
        }
    }

}