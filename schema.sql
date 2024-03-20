CREATE TABLE product(
    id SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL
);

CREATE TABLE shelf(
    id SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL
);

CREATE TABLE orders(
    number SERIAL PRIMARY KEY
);

CREATE TABLE product_shelf(
    product_id INT REFERENCES product(id) ON DELETE CASCADE,
    shelf_id INT REFERENCES shelf(id) ON DELETE CASCADE,
    is_main_shelf BOOLEAN NOT NULL
);

CREATE TABLE product_orders(
    product_id INT REFERENCES product(id) ON DELETE CASCADE,
    order_num INT REFERENCES orders(number) ON DELETE CASCADE,
    count INT
);

INSERT INTO product(title) VALUES('Ноутбук');
INSERT INTO product(title) VALUES('Телевизор');
INSERT INTO product(title) VALUES('Телефон');
INSERT INTO product(title) VALUES('Системный блок');
INSERT INTO product(title) VALUES('Часы');
INSERT INTO product(title) VALUES('Микрофон');

INSERT INTO shelf(title) VALUES('А');
INSERT INTO shelf(title) VALUES('Б');
INSERT INTO shelf(title) VALUES('Ж');
INSERT INTO shelf(title) VALUES('З');
INSERT INTO shelf(title) VALUES('В');

INSERT INTO orders(number) VALUES(10);
INSERT INTO orders(number) VALUES(11);
INSERT INTO orders(number) VALUES(14);
INSERT INTO orders(number) VALUES(15);

INSERT INTO product_orders VALUES (1, 10, 2);
INSERT INTO product_orders VALUES (2, 11, 3);
INSERT INTO product_orders VALUES (1, 14, 3);
INSERT INTO product_orders VALUES (3, 10, 1);
INSERT INTO product_orders VALUES (4, 14, 4);
INSERT INTO product_orders VALUES (5, 15, 1);
INSERT INTO product_orders VALUES (6, 10, 1);

INSERT INTO product_shelf VALUES (1, 1, true);
INSERT INTO product_shelf VALUES (2, 1, true);
INSERT INTO product_shelf VALUES (3, 2, true);
INSERT INTO product_shelf VALUES (4, 3, true);
INSERT INTO product_shelf VALUES (5, 3, true);
INSERT INTO product_shelf VALUES (6, 3, true);
INSERT INTO product_shelf VALUES (3, 4, false);
INSERT INTO product_shelf VALUES (3, 5, false);
INSERT INTO product_shelf VALUES (5, 1, false);


