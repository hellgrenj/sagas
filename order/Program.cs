using System.Data;
using FluentValidation;
using MediatR;
using model.Repositories;
using Npgsql;
using order.Infrastructure;
using order.Model;
using order.Repositories;

var connectionString = Environment.GetEnvironmentVariable("POSTGRES_CONNECTION_STRING");
if (string.IsNullOrEmpty(connectionString))
    connectionString = "Host=orderdb;Database=order;Username=orderusr;Password=orderpwd";


IHost host = Host.CreateDefaultBuilder(args)
    .ConfigureServices(services =>
    {
        services.AddHostedService<RabbitListener>();
        services.AddMediatR(typeof(Program));
        services.AddTransient<IDbConnection>((sp) => new NpgsqlConnection(connectionString));
        services.AddTransient<IRepository<Order>, OrderRepository>();
        services.AddValidatorsFromAssemblyContaining<Program>();
        services.AddSingleton<IRabbitEmitter, RabbitEmitter>();
        services.AddSingleton<IRabbitConnectionHandler, RabbitConnectionHandler>();
    })
    .Build();

var rabbitEmitter = host.Services.GetService<IRabbitEmitter>();
rabbitEmitter.Open();
var applicationLifetime = host.Services.GetService<IHostApplicationLifetime>();
applicationLifetime.ApplicationStopping.Register(() => rabbitEmitter.Close());
await host.RunAsync();
