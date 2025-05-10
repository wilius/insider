create sequence if not exists notifications.seq__message_id;

create table if not exists notifications.message
(
    id           bigint       not null primary key,
    phone_number text         not null,
    message      varchar(120) not null,
    status       text         not null,
    create_date  timestamp without time zone,
    update_date  timestamp without time zone
);