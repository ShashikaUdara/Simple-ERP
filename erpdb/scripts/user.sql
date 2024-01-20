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

-- Create the user_session table
CREATE TABLE user_session (
    id INT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    session_token VARCHAR(255) NOT NULL,
    ip_address VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255) NOT NULL,
    creation_time TIMESTAMP NOT NULL,
    last_activity_time TIMESTAMP NOT NULL
);