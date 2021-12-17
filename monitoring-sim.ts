import { connect } from "https://deno.land/x/amqp@v0.17.0/mod.ts";

console.log("listening for all topics on exchange order.topics".toUpperCase());

const queueName = "monitoring_queue";

const connection = await connect({ hostname: "127.0.0.1" });

const channel = await connection.openChannel();

const exchange = "order.topics";
await channel.declareExchange({ exchange, type: "topic" });

await channel.declareQueue({ queue: queueName });
await channel.bindQueue({
  routingKey: "#", // all topics
  queue: queueName,
  exchange, // exchange
});

await channel.consume(
  { queue: queueName },
  async (args, _, data) => {
    const ev = JSON.parse(new TextDecoder().decode(data));
    console.log("received event ", ev.name);
    console.log(ev);
    await channel.ack({ deliveryTag: args.deliveryTag });
  },
);
