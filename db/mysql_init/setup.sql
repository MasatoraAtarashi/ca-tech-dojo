drop table if exists users cascade;
drop table if exists characters cascade;
drop table if exists user_characters;

create table users (
                       id serial primary key,
                       token varchar(255) not null,
                       name varchar(255) not null
);

create table characters (
                            id serial primary key,
                            name varchar(255) not null,
                            weight integer not null
);

create table user_characters (
                                 id serial primary key,
                                 user_id integer  references users(id),
                                 character_id integer  references characters(id)
);

insert into characters values (1, 'Satan', 10);
insert into characters values (2, 'Bahamut', 20);
insert into characters values (3, 'Goblin', 70);
