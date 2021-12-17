# Sagas

A simplified implementation (demo) of an order fulfilment flow as a choreographed [Saga](https://microservices.io/patterns/data/saga.html).

## Prerequisites
Docker Desktop (with kubernetes enabled) and skaffold   
.Net 6    
Deno   
Go   

## Run
``skaffold run``  

and when everything is up and running start an order flow 
with: ```deno run -A sales-app-sim.ts``` while monitoring with ```deno run -A monitoring-sim.ts``` 
(also check the logs for notification service..)  



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

**Order Fulfillment flow (as a choreographed saga)**

1. sales-app-sim emits order.placed event
2. orderService creates an order and emits order.created event  
3. warehouseService checks if order is in stock and reserves the item and finally emits item.reserved event or item.notinstock event
4. paymentService takes money from the customer and emits payment.completed event
5. warehouseService ships the order and emits order.shipped event

(notification service communicates the order status to the customer for several steps(events) in this flow)


**Rollback example with compensating transactions**

todo: describe...∏

