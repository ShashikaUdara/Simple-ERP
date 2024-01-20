-- Create the nothing database schema
CREATE SCHEMA IF NOT EXISTS erp;

-- Switch to the nothing schema
USE erp;

-- Create the user table
CREATE TABLE user (
    id          INT          AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    password    VARCHAR(255) NOT NULL
);