using FluentValidation;
using MediatR;
using model.Repositories;
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
    private readonly IRepository<Order> _orderRepository;
    private readonly IMediator _mediator;
    public CancelOrderCommandHandler(ILogger<PlaceOrderHandler> logger, CancelOrderCommandValidator validator, IRepository<Order> orderRepository, IMediator mediator)
    {
        _logger = logger;
        _validator = validator;
        _orderRepository = orderRepository;
        _mediator = mediator;
    }
    public async Task<Unit> Handle(CancelOrderCommand cmd, CancellationToken cancellationToken)
    {
        var validationResult = _validator.Validate(cmd);
        if (validationResult.IsValid)
        {
            var order = await _orderRepository.GetAsync(cmd.itemsNotInStockEvent.Order.Id.Value);
            order.ChangeState(OrderStates.Cancelled);
            await _orderRepository.SaveAsync(order);
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