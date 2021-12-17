FROM mcr.microsoft.com/dotnet/sdk:6.0 AS build-env
WORKDIR /app
COPY ./order.csproj .
RUN dotnet restore
COPY . .
RUN dotnet publish -c Release -o out

FROM mcr.microsoft.com/dotnet/aspnet:6.0

# create and use non-root user.. 
RUN adduser \
  --disabled-password \
  --home /apirunner \
  --gecos '' apirunner \
  && chown -R apirunner /apirunner
USER apirunner

WORKDIR /app
COPY --from=build-env /app/out ./
ENTRYPOINT ["dotnet", "./order.dll"]