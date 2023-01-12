SET SEARCH_PATH TO "shop";

CREATE SEQUENCE "warehouse_id_seq";
CREATE TABLE "warehouses" (
                            "id" INT DEFAULT nextval('warehouse_id_seq'),
                            "name" VARCHAR (100) NOT NULL,
                            PRIMARY KEY (id)
);
