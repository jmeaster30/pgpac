create extension citext;
create type film_rating as enum ('G', 'PG', 'PG-13', 'R', 'NC-17');
create table film (
    film_id bigint GENERATED ALWAYS AS IDENTITY,
    title varchar(255),
    description citext,
    release_year int,
    language_id int references language(language_id),
    original_language_id int references language(language_id),
    rental_duration int not null,
    rental_rate numeric(4, 2),
    length int,
    replacement_cost numeric(5, 2),
    rating film_rating default 'G',
    last_update timestamp
)