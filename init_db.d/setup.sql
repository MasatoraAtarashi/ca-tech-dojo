drop table if exists users;
drop table if exists api_responses;

create table users (
                       id serial primary key,
                       username varchar(255),
                       firstName varchar(255),
                       lastName varchar(255),
                       email varchar(255),
                       password varchar(255),
                       phone varchar(255),
                       userStatus integer
);

create table api_responses (
                               code   integer,
                               type varchar(255),
                               message varchar(255)
);
