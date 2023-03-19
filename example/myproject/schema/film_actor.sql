create table film_actor (
    actor_id bigint not null references actor(actor_id),
    film_id bigint not null,
    last_update timestamp not null default CURRENT_TIMESTAMP,

    primary key (actor_id, film_id),
    foreign key (film_id) references film(film_id)
);

create index idx_film_actor_film on film_actor(film_id);