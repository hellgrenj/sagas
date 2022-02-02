import { connect } from "https://deno.land/x/amqp@v0.17.0/mod.ts";

console.log(
  "Simulating a sales app placing an order every 10 second (close with ctrl + c)",
);

const simulateOrderPlaced = async () => {
  const connection = await connect();
  const channel = await connection.openChannel();

  const exchange = "order.topics";
  await channel.declareExchange({ exchange: exchange, type: "topic" });

  const order = {
    item: "iPhone 13",
    price: 2000,
    quantity: 3,
  };
  const msg = {
    order,
    correlationId: crypto.randomUUID(),
    messageId: crypto.randomUUID(),
    name: "order placed",
  };

  console.log("placing order", order);

  await channel.publish(
    { exchange: "order.topics", routingKey: "order.placed" },
    { contentType: "application/json" },
    new TextEncoder().encode(JSON.stringify(msg)),
  );

  await connection.close();
};

simulateOrderPlaced();
setInterval(async () => {
  await simulateOrderPlaced();
}, 10000);
