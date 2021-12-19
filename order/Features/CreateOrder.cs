using System.Data;
using Dapper;
using FluentValidation;
using MediatR;
using order.Events.Inbound;
using order.Events.Outbound;
using order.Model;

namespace order.Features;

public record CreateOrderCommand(OrderPlacedEvent orderPlacedEvent) : IRequest<Unit>;


public class CreateOrderCommandValidator : AbstractValidator<CreateOrderCommand>
{
    public CreateOrderCommandValidator()
    {
        RuleFor(x => x.orderPlacedEvent.Order.Item).NotEmpty();
        RuleFor(x => x.orderPlacedEvent.Order.Quantity).GreaterThan(0);
        RuleFor(x => x.orderPlacedEvent.Order.Price).GreaterThan(0);
    }
}

public class PlaceOrderHandler : IRequestHandler<CreateOrderCommand, Unit>
{
    private readonly ILogger<PlaceOrderHandler> _logger;
    private readonly CreateOrderCommandValidator _validator;
    private readonly IDbConnection _connection;

    private readonly IMediator _mediator;
    public PlaceOrderHandler(ILogger<PlaceOrderHandler> logger, CreateOrderCommandValidator validator, IDbConnection connection, IMediator mediator)
    {
        _logger = logger;
        _validator = validator;
        _connection = connection;
        _mediator = mediator;
    }

    public async Task<Unit> Handle(CreateOrderCommand cmd, CancellationToken cancellationToken)
    {
        var validationResult = _validator.Validate(cmd); 
        if (validationResult.IsValid)
        {
            var order = cmd.orderPlacedEvent.Order;
            order.State = OrderStates.Pending;
            var sql = "INSERT INTO ordertable (correlationid, item, quantity, price, state) Values (@correlationid, @item, @quantity, @price,  @state) RETURNING Id;";
            order.Id = (int)await _connection.ExecuteScalarAsync(sql,
            new { correlationid = cmd.orderPlacedEvent.CorrelationId, item = order.Item, quantity = order.Quantity, price = order.Price, state = order.State });
            _logger.LogInformation($"Order created with id {order.Id}");
            await _mediator.Publish(new OrderCreatedEvent { CorrelationId = cmd.orderPlacedEvent.CorrelationId, Order = order });
        }
        else
        {
            _logger.LogError($"PlaceOrderCommand validation failed \n {validationResult.ToString()}");
        }
        return Unit.Value;

    }

}