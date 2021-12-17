using System.Data;
using FluentValidation;
using MediatR;
using Npgsql;
using warehouse.Infrastructure;

var connectionString = Environment.GetEnvironmentVariable("POSTGRES_CONNECTION_STRING");
if (string.IsNullOrEmpty(connectionString))
    connectionString = "Host=warehousedb;Database=warehouse;Username=warehouseusr;Password=warehousepwd";


IHost host = Host.CreateDefaultBuilder(args)
    .ConfigureServices(services =>
    {
        services.AddHostedService<RabbitListener>();
        services.AddMediatR(typeof(Program));
        services.AddTransient<IDbConnection>((sp) => new NpgsqlConnection(connectionString));
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
