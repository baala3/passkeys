CREATE DATABASE mydb;

\c mydb;

CREATE SCHEMA myschema;

CREATE ROLE myuser
WITH LOGIN PASSWORD 'mypassword';

GRANT ALL PRIVILEGES ON SCHEMA myschema TO myuser;
