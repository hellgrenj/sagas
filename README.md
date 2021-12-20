# Sagas

A simplified implementation (demo) of an order fulfilment flow as a choreographed [Saga](https://microservices.io/patterns/data/saga.html). 

### Purpose
Demo-environment for a workshop in event-driven architecture and the Saga pattern. 

## Prerequisites
Docker Desktop (with kubernetes enabled) and skaffold   
.NET 6    
Deno   
Go   

## Run
``skaffold run``  

and when everything is up and running start an order flow 
with: ```deno run -A sales-app-sim.ts``` while monitoring with ```deno run -A monitoring-sim.ts```   
(also check the logs for the notification service ``kubectl logs <podname> -f``)  



## Scenario  

**Architecture**  

```
          sales-app-sim 
                |
                |
              rabbit----------- monitoring-sim
                |
                |
                |
-------------------------------------------          
   |         |         |           |               
 order   warehouse   payment  notification
service   service    service    service 

```
The Order service is built with .NET 6 and Postgres (db migrations with Roundhouse)  
The Warehouse service is built with .NET 6 and Postgres (db migrations with Roundhouse)  
The Payment service is built with Go and MongoDB  
The Notification service is built with Go  
sales-app-sim and monitoring-sim are Deno scripts. 

The RabbitMQ management UI is exposed on localhost:15672 (guest/guest)  

The databases are not exposed on localhost so you need to ``kubectl port-forward <podname> 5432:5432`` (for Postgres) or ``kubectl port-forward <podname> 27017:27017`` (for MongoDB) to access them locally (with DBeaver and mongosh or what ever tool you prefer).

**Order Fulfillment flow (as a choreographed saga)**

1. sales-app-sim emits ``order.placed`` event
2. orderService creates an order and emits ``order.created`` event  
3. warehouseService checks if the order is in stock and reserves the items and finally emits ``order.items.reserved`` event or ``order.items.notinstock`` event
4. paymentService takes money from the customer and emits ``order.payment.completed`` event
5. warehouseService ships the order and emits ``order.shipped`` event  
6. orderService completes the order and emits an ``order.completed`` event

(notification service communicates the order status to the customer for several steps(events) in this flow)


**Rollback example with compensating transactions**

For every order simulated there is a 50% risk of the items not being in stock. When this happens the warehouse service will emit an ``order.items.notinstock`` event. The order service will listen for this event and cancel the order and emit an ``order.cancelled`` event. Notification service listens for all events and informs the customers about important steps.

**Exercise**

Are there more rollback scenarios? How would we handle the payment failing? How would we handle the shipping failing?

