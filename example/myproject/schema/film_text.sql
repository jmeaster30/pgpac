create table film_text (
    film_id bigint not null primary key,
    title varchar(255) not null,
    description citext
);

create index idx_film_text_full_text on film_text(title, description);