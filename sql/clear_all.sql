DELETE FROM value;
DELETE FROM vault;

ALTER SEQUENCE vault_id_seq RESTART WITH 1;
ALTER SEQUENCE value_id_seq RESTART WITH 1;