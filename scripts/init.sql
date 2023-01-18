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

INSERT INTO currency (code, denominator)
VALUES ('SGD', 100);

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
    timezone varchar(30) default 'Asia/Singapore'::character varying not null,
    create_time    timestamp with time zone not null default NOW(),
    update_time    timestamp with time zone not null default NOW()
);

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
    description varchar(20)             not null,
    user_id     bigint                   not null
        references app_user,
    amount      bigint  default 0        not null,
    currency    char(3) default 'SGD'::bpchar
        constraint transaction_currency_fk
            references currency
);

comment on column transaction.amount is 'Normalised to the lowest denominator';

create table message_context
(
    id            serial primary key,
    chat_id       int not null,
    message_id    int not null,
    message       text not null,
    created_at    timestamp with time zone not null
);

-- Take note of the newline in the reply text string
insert into transaction_type (id, name, multiplier, display_order, reply_text)
values (1,'🔴 Spent', -1, 2, 'You spent <b>%s</b> on <b>%s</b>
    ');

insert into category (name, transaction_type_id, display_order)
values ('Child', 1, 1),
       ('Education', 1, 2),
       ('Entertainment', 1, 3),
       ('Food', 1, 4),
       ('Gifts', 1, 5),
       ('Health', 1, 6),
       ('Housing', 1, 7),
       ('Insurance', 1, 8),
       ('Shopping', 1, 9),
       ('Taxes', 1, 10),
       ('Transport', 1, 11),
       ('Other', 1, 99);
