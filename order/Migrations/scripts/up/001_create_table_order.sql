CREATE TABLE ordertable (
    id serial PRIMARY KEY,
    correlationid TEXT NOT NULL UNIQUE,
    item TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    price INTEGER NOT NULL,
    state TEXT NOT NULL
);