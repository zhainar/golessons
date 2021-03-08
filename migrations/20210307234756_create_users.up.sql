create table users (
    id SERIAL primary key not null,
    email varchar not null unique,
    encrypted_password varchar not null
);