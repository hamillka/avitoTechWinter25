CREATE DATABASE avito_shop_service;

\connect "avito_shop_service";

CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT        NOT NULL,
    coins    INT DEFAULT 1000 CHECK (coins >= 0)
);

CREATE INDEX idx_username_password ON users (username, password);

CREATE TABLE IF NOT EXISTS transactions
(
    id          SERIAL PRIMARY KEY,
    sender_id   INT,
    FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
    receiver_id INT,
    FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE,
    amount      INT CHECK (amount > 0)
);

CREATE INDEX idx_sender_id ON transactions (sender_id);
CREATE INDEX idx_receiver_id ON transactions (receiver_id);

CREATE TABLE IF NOT EXISTS merch
(
    id   SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    cost INT CHECK (cost > 0)
);

CREATE TABLE IF NOT EXISTS inventory
(
    id       SERIAL PRIMARY KEY,
    user_id  INT,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    merch_id INT,
    FOREIGN KEY (merch_id) REFERENCES merch (id) ON DELETE CASCADE,
    amount   INT CHECK (amount > 0)
);

CREATE INDEX idx_inventory_user ON inventory (user_id);

INSERT INTO avito_shop_service.public.users
VALUES (0, 'AvitoShop', 'AvitoPassword', 0),
       (105005, 'initUser', 'initUserPassword', 1000);

INSERT INTO avito_shop_service.public.merch
VALUES (1, 't-shirt', 80),
       (2, 'cup', 20),
       (3, 'book', 50),
       (4, 'pen', 10),
       (5, 'powerbank', 200),
       (6, 'hoody', 300),
       (7, 'umbrella', 200),
       (8, 'socks', 10),
       (9, 'wallet', 50),
       (10, 'pink-hoody', 500);
