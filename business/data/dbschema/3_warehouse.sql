SET SEARCH_PATH TO shop;

CREATE SEQUENCE warehouse_id_seq;
CREATE TABLE warehouses
(
    PRIMARY KEY (id),
    id   INT DEFAULT nextval('warehouse_id_seq'),
    name VARCHAR(100) NOT NULL
);
