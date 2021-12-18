

using MediatR;
using warehouse.Events.Inbound;
using warehouse.Events.Outbound;

namespace warehouse.Inventory.Features;

public record ShipOrderCommand(PaymentCompleted paymentCompletedEvent) : IRequest<Unit>;


public class ShipOrderCommandHandler : IRequestHandler<ShipOrderCommand, Unit>
{
    private readonly ILogger<ShipOrderCommand> _logger;

    private readonly IMediator _mediator;
    public ShipOrderCommandHandler(ILogger<ShipOrderCommand> logger, IMediator mediator)
    {
        _logger = logger;
        _mediator = mediator;
    }

    public async Task<Unit> Handle(ShipOrderCommand cmd, CancellationToken cancellationToken)
    {
        _logger.LogInformation($"Shipping order {cmd.paymentCompletedEvent.OrderId}");
        
        // do the thing...

        await _mediator.Publish(new OrderShippedEvent { CorrelationId = cmd.paymentCompletedEvent.CorrelationId, OrderId = cmd.paymentCompletedEvent.OrderId });

        return Unit.Value;
    }
}






