create table language (
    language_id bigint generated by default as identity primary key,
    name citext not null,
    last_update timestamp not null default current_timestamp
);
