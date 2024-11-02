CREATE TABLE IF NOT EXISTS shopping_items(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    store_id INT
);