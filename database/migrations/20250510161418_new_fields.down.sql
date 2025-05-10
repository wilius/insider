alter table notifications.message
    drop column if exists provider_message_id,
    drop column if exists provider;