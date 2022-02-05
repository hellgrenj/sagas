using System.Data;
using Dapper;
using FluentValidation;
using MediatR;
using model.Repositories;
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
    private readonly IRepository<Order> _orderRepository;

    private readonly IMediator _mediator;
    public CompleteOrderCommandHandler(ILogger<PlaceOrderHandler> logger, CompleteOrderCommandValidator validator, IRepository<Order> orderRepository, IMediator mediator)
    {
        _logger = logger;
        _validator = validator;
        _orderRepository = orderRepository;
        _mediator = mediator;
    }

    public async Task<Unit> Handle(CompleteOrderCommand cmd, CancellationToken cancellationToken)
    {
        var validationResult = _validator.Validate(cmd);
        if (validationResult.IsValid)
        {
            var order = await _orderRepository.GetAsync(cmd.OrderShippedevent.OrderId);
            order.ChangeState(OrderStates.Completed);
            order.Id = await _orderRepository.SaveAsync(order);
            _logger.LogInformation($"Completed Order with id {order.Id.Value}");
            await _mediator.Publish(new OrderCompletedEvent { CorrelationId = cmd.OrderShippedevent.CorrelationId, OrderId = order.Id.Value });
        }
        else
        {
            _logger.LogError($"CompleteOrderCommand validation failed \n {validationResult.ToString()}");
        }
        return Unit.Value;
    }

}