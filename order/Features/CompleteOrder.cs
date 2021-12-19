using System.Data;
using Dapper;
using FluentValidation;
using MediatR;
using order.Events.Inbound;
using order.Events.Outbound;
using order.Model;

namespace order.Features;

public record CompleteOrderCommand(OrderShippedevent OrderShippedevent) : IRequest<Unit>;


public class CompleteOrderCommandValidator : AbstractValidator<CompleteOrderCommand>
{
    public CompleteOrderCommandValidator()
    {
        RuleFor(x => x.OrderShippedevent.CorrelationId).NotEmpty();
    }
}
public class CompleteOrderCommandHandler : IRequestHandler<CompleteOrderCommand, Unit>
{
    private readonly ILogger<PlaceOrderHandler> _logger;
    private readonly CompleteOrderCommandValidator _validator;
    private readonly IDbConnection _connection;

    private readonly IMediator _mediator;
    public CompleteOrderCommandHandler(ILogger<PlaceOrderHandler> logger, CompleteOrderCommandValidator validator, IDbConnection connection, IMediator mediator)
    {
        _logger = logger;
        _validator = validator;
        _connection = connection;
        _mediator = mediator;
    }
    // TODO should have an audit table where we store all related information for each event (order state change..)
    // could use Postgres json support for this..
    public async Task<Unit> Handle(CompleteOrderCommand cmd, CancellationToken cancellationToken)
    {
        var validationResult = _validator.Validate(cmd);
        if (validationResult.IsValid)
        {
            var orderId = cmd.OrderShippedevent.OrderId;
            
            var sql = "UPDATE ordertable SET state=@state WHERE id=@id;";
            await _connection.ExecuteScalarAsync(sql,
            new { state = OrderStates.Completed, id = orderId });
            _logger.LogInformation($"Completed Order with id {orderId}");
            await _mediator.Publish(new OrderCompletedEvent { CorrelationId = cmd.OrderShippedevent.CorrelationId, OrderId = orderId });
        }
        else
        {
            _logger.LogError($"CompleteOrderCommand validation failed \n {validationResult.ToString()}");
        }
        return Unit.Value;
    }

}