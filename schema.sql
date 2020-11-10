-- schema.sql
-- Since we might run the import many times we'll drop if exists

DROP DATABASE IF EXISTS mb_scheduler_db;

CREATE DATABASE mb_scheduler_db;

\c mb_scheduler_db

CREATE TABLE IF NOT EXISTS schedule_rt (
    id SERIAL PRIMARY KEY,
    runtime INTEGER,
    fullname VARCHAR,
    username VARCHAR,
    password VARCHAR,
    classtime VARCHAR,
    weekday VARCHAR,
    date VARCHAR
)