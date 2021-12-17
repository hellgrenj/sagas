CREATE TABLE processedmessages (
    id serial PRIMARY KEY,
    messageid TEXT NOT NULL UNIQUE
);