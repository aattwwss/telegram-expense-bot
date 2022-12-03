create table app_user
(
    id       bigint                                                  not null
        constraint telegram_user_pkey
        primary key
        constraint telegram_user_id_key
        unique,
    locale   varchar(10) default 'en'::character varying             not null,
    currency char(3)     default 'SGD'::bpchar                       not null
        constraint user_currency_fk
        references currency,
    timezone varchar(30) default 'Asia/Singapore'::character varying not null
);

create table transaction_type
(
    id            serial
        primary key,
    name          varchar(10)                               not null
        unique,
    multiplier    smallint    default '1'::smallint         not null,
    display_order integer     default 0                     not null,
    reply_text    varchar(64) default ''::character varying not null
);

create table currency
(
    code        varchar(3)           not null
        constraint currency_pk
            primary key
        constraint check_code
            check ((code)::text ~ '^[A-Z]{3}$'::text),
    denominator smallint default 100 not null
);

comment on constraint check_code on currency is '3 letter uppercase';

create table category
(
    id                  serial
        primary key,
    name                varchar(50) not null,
    transaction_type_id integer     not null
        references transaction_type,
    unique (name, transaction_type_id)
);

create table transaction
(
    id          serial
        primary key,
    datetime    timestamp with time zone not null,
    category_id integer                  not null
        references category,
    description varchar(255)             not null,
    user_id     bigint                   not null
        references app_user,
    amount      bigint  default 0        not null,
    currency    char(3) default 'SGD'::bpchar
        constraint transaction_currency_fk
            references currency
);

comment on column transaction.amount is 'Normalised to the lowest denominator';