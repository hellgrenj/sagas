using System.Data;
using Dapper;
using FluentValidation;
using MediatR;
using order.Events.Inbound;
using order.Events.Outbound;
using order.Model;

namespace order.Features;

public record CancelOrderCommand(ItemsNotInStockEvent itemsNotInStockEvent) : IRequest<Unit>;


public class CancelOrderCommandValidator : AbstractValidator<CancelOrderCommand>
{
    public CancelOrderCommandValidator()
    {
        RuleFor(x => x.itemsNotInStockEvent.CorrelationId).NotEmpty();
    }
}
public class CancelOrderCommandHandler : IRequestHandler<CancelOrderCommand, Unit>
{
    private readonly ILogger<PlaceOrderHandler> _logger;
    private readonly CancelOrderCommandValidator _validator;
    private readonly IDbConnection _connection;

    private readonly IMediator _mediator;
    public CancelOrderCommandHandler(ILogger<PlaceOrderHandler> logger, CancelOrderCommandValidator validator, IDbConnection connection, IMediator mediator)
    {
        _logger = logger;
        _validator = validator;
        _connection = connection;
        _mediator = mediator;
    }
    // TODO should have an audit table where we store all related information for each event (order state change..)
    // could use Postgres json support for this..
    public async Task<Unit> Handle(CancelOrderCommand cmd, CancellationToken cancellationToken)
    {
        var validationResult = _validator.Validate(cmd);
        if (validationResult.IsValid)
        {
            var order = cmd.itemsNotInStockEvent.Order;
            order.State = OrderStates.Cancelled;
            var sql = "UPDATE ordertable SET state=@state WHERE id=@id;";
            await _connection.ExecuteScalarAsync(sql,
            new { state = order.State, id = order.Id });
            _logger.LogInformation($"Cancelled Order with id {order.Id}");
            await _mediator.Publish(new OrderCancelledEvent { CorrelationId = cmd.itemsNotInStockEvent.CorrelationId, 
            Reason = "Items not in stock", Order = order });
        }
        else
        {
            _logger.LogError($"CancelOrderCommand validation failed \n {validationResult.ToString()}");
        }
        return Unit.Value;
    }

}