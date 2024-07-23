CREATE TABLE IF NOT EXISTS message
(
    id      serial       primary key,
    createdTime  timestamp    not null,
    fromUserId int not null,
    ChatId  int not null,
    Text text not null,
    Attachments  text[]
);

CREATE TABLE IF NOT EXISTS user_last_seen
(
    id      serial       primary key,
    createdTime  timestamp    not null,
    fromUserId int not null,
    ChatId  int not null
);

CREATE TABLE IF NOT EXISTS last_message
(
    id      serial       primary key,
    ChatId  int not null,
    createdTime  timestamp    not null
);

CREATE TABLE IF NOT EXISTS chats
(
    id      serial       primary key,
    Users  integer[],
    createdTime  timestamp    not null
);

CREATE TABLE IF NOT EXISTS users
(
    id      serial       primary key,
    login    varchar(50)   UNIQUE   not null,
    ChatIds  integer[]
);

CREATE TABLE IF NOT EXISTS devices
(
    id      serial       primary key,
    deviceId    text      not null,
    userId  integer not null,
    lastSeen  timestamp
);

CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_users_login ON public.users USING btree (UPPER((login)::text));