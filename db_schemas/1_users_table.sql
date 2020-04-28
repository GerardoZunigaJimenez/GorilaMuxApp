-- auto-generated definition
create table awesome_user
(
    dbId serial       not null
        constraint users_pkey
            primary key,
    user_id integer not null,
    first_name varchar(50)  not null,
    last_name  varchar(50)  not null,
    email      varchar(75)  not null,
    gender      varchar(50)  not null
);

alter table awesome_user
    owner to golang_user;

