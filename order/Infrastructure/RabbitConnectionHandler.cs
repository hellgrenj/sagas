using RabbitMQ.Client;

namespace order.Infrastructure;


public interface IRabbitConnectionHandler
{
    IConnection TryOpen(ConnectionFactory factory, ILogger logger, int attempt = 0);
    void Close(ILogger logger, IConnection connection = null, IModel channel = null);

}
public class RabbitConnectionHandler : IRabbitConnectionHandler
{
    public IConnection TryOpen(ConnectionFactory factory, ILogger logger, int attempt = 0)
    {

        try
        {
            return factory.CreateConnection();
        }
        catch (Exception e)
        {
            attempt++;
            if (attempt < 10)
            {
                logger.LogWarning($"Failed to connect to RabbitMQ attempt {attempt} retrying in 5 seconds");
                System.Threading.Thread.Sleep(5000);
                return TryOpen(factory, logger, attempt);
            }
            else
            {
                logger.LogError($"Failed to connect after 10 attemps, exception {e.ToString()}");
                throw;
            }
        }
    }
    public void Close(ILogger logger, IConnection connection = null, IModel channel = null)
    {
        if (channel != null)
        {
            logger.LogInformation("closing channel");
            channel.Close(200, "Goodbye");

        }
        if (connection != null)
        {
            logger.LogInformation("closing connection");
            connection.Close();
        }
    }
}