CREATE TABLE reservations (
    id serial PRIMARY KEY,
    correlationid TEXT NOT NULL UNIQUE,
    orderid INTEGER NOT NULL UNIQUE,
    item TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    price INTEGER NOT NULL
);