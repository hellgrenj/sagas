using System.Data;
using Dapper;
using MediatR;
using Npgsql;

namespace order.Infrastructure;

public record TryMarkMessageAsProcessedCommand(string messageId) : IRequest<bool>;

public class TryMarkMessageAsProcessedHandler : IRequestHandler<TryMarkMessageAsProcessedCommand, bool>
{
    private readonly ILogger<TryMarkMessageAsProcessedHandler> _logger;
    private readonly IDbConnection _connection;

    public TryMarkMessageAsProcessedHandler(ILogger<TryMarkMessageAsProcessedHandler> logger, IDbConnection connection)
    {
        _logger = logger;
        _connection = connection;
    }
    public async Task<bool> Handle(TryMarkMessageAsProcessedCommand cmd, CancellationToken cancellationToken)
    {
        var messageAlreadyProcesed = false;
        try
        {
            var sql = "INSERT INTO processedmessages (messageid) Values (@messageid);";
            await _connection.ExecuteScalarAsync(sql,
            new { messageid = cmd.messageId });
            _logger.LogInformation($"message with messageId {cmd.messageId} marked as processed");
        }
        catch (PostgresException ex)
        {
            if (ex.SqlState == "23505")
            {
                _logger.LogInformation($"message with messageId {cmd.messageId} is already marked as processed");
                messageAlreadyProcesed = true;
            }
            else
            {
                throw;
            }
        }
        return messageAlreadyProcesed;
    }


}