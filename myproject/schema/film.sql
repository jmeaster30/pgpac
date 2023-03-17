create extension citext;
create type test;
create type film_rating as enum ('G', 'PG', 'PG-13', 'R', 'NC-17');
create table film (
    film_id bigint GENERATED ALWAYS AS IDENTITY,
    title varchar(255),
    description text,
    release_year int,
    language_id int,
    original_language_id int,
    rental_duration int,
    rental_rate numeric(4, 2),
    length int,
    replacement_cost numeric(5, 2),
    rating film_rating default 'G',
    last_update timestamp
)