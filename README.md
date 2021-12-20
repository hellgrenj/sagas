# Sagas

A simplified implementation (demo) of an order fulfilment flow as a choreographed [Saga](https://microservices.io/patterns/data/saga.html).

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

The databases are not exposed on localhost so you need to ``kubectl port-forward <podname> 5462:5462`` (for Postgres) or ``kubectl port-forward <podname> 27017:27017`` (for MongoDB) to access them locally (with DBeaver and mongosh or what ever tool you prefer).

**Order Fulfillment flow (as a choreographed saga)**

1. sales-app-sim emits order.placed event
2. orderService creates an order and emits order.created event  
3. warehouseService checks if order is in stock and reserves the item and finally emits item.reserved event or item.notinstock event
4. paymentService takes money from the customer and emits payment.completed event
5. warehouseService ships the order and emits order.shipped event

(notification service communicates the order status to the customer for several steps(events) in this flow)


**Rollback example with compensating transactions**

For every order simulated there is a 50% risk of the items not being in stock. When this happens the warehouse service will emit an item.notinstock event. The order service will listen of this event and cancel the order and emit an order.cancelled event. Notification service listens for all events and informs the customers about important steps.


