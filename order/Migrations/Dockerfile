FROM dotnetroundhouse/roundhouse:1.3.1
COPY . .
# dbt = database type (default MS SQL.. aka sqlserver)
ENTRYPOINT ["dotnet", "rh.dll", "/f=./scripts", "/dbt=postgres", "/c=Host=orderdb;Database=order;Username=orderusr;Password=orderpwd"]
