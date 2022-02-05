using System.Data;
using Dapper;
using FluentValidation;
using MediatR;
using model.Repositories;
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
    private readonly IRepository<Order> _orderRepository;

    private readonly IMediator _mediator;
    public PlaceOrderHandler(ILogger<PlaceOrderHandler> logger, CreateOrderCommandValidator validator, IRepository<Order> orderRepository, IMediator mediator)
    {
        _logger = logger;
        _validator = validator;
        _orderRepository = orderRepository;
        _mediator = mediator;
    }

    public async Task<Unit> Handle(CreateOrderCommand cmd, CancellationToken cancellationToken)
    {
        var validationResult = _validator.Validate(cmd); 
        if (validationResult.IsValid)
        {
            var order = new Order(cmd.orderPlacedEvent.Order, cmd.orderPlacedEvent.CorrelationId);
            order.ChangeState(OrderStates.Pending);
            order.Id = await _orderRepository.SaveAsync(order);
            _logger.LogInformation($"Order created with id {order.Id.Value}");
            await _mediator.Publish(new OrderCreatedEvent { CorrelationId = cmd.orderPlacedEvent.CorrelationId, Order = order });
        }
        else
        {
            _logger.LogError($"PlaceOrderCommand validation failed \n {validationResult.ToString()}");
        }
        return Unit.Value;

    }

}