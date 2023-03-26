create table film_category (
    film_id bigint not null references film(film_id),
    category_id bigint not null references film(film_id),
    last_update timestamp not null default CURRENT_TIMESTAMP,

    primary key (film_id, category_id)
);