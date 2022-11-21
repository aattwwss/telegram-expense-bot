create table expenditure_bot.transaction_type
(
    id         serial
        primary key,
    name       varchar(10)                    not null
        unique,
    multiplier smallint default '1'::smallint not null
);


create table expenditure_bot.category
(
    id                  serial
        primary key,
    name                varchar(50) not null,
    transaction_type_id integer     not null
        references expenditure_bot.transaction_type,
    unique (name, transaction_type_id)
);


create table expenditure_bot.telegram_user
(
    id         bigint      not null
        primary key
        unique,
    is_bot     boolean     not null,
    first_name varchar(64) not null,
    last_name  varchar(64),
    username   varchar(32)
);

create table expenditure_bot.currency
(
    code        varchar(3)           not null
        constraint currency_pk
            primary key
        constraint check_code
            check ((code)::text ~ '^[A-Z]{3}$'::text),
    denominator smallint default 100 not null
);

comment on constraint check_code on expenditure_bot.currency is '3 letter uppercase';

alter table expenditure_bot.currency
    owner to postgres;

create table expenditure_bot.transaction
(
    id          serial
        primary key,
    datetime    timestamp with time zone not null,
    category_id integer                  not null
        references expenditure_bot.category,
    description varchar(255),
    user_id     bigint                   not null
        references expenditure_bot.telegram_user,
    amount      bigint  default 0        not null,
    currency    char(3) default 'SGD'::bpchar
        constraint transaction_currency_fk
            references expenditure_bot.currency
);

comment on column expenditure_bot.transaction.amount is 'Normalised to the lowest denominator';
