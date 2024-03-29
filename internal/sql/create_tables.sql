CREATE TABLE orders (
  order_id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  order_data jsonb
);

SELECT * FROM orders;