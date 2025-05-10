alter table notifications.message
    add column if not exists provider_message_id text,
    add column if not exists provider            text;