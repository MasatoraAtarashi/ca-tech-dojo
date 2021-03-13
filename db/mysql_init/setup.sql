drop table if exists users;

create table users (
                       id serial primary key,
                       token varchar(255) not null,
                       name varchar(255) not null
);
