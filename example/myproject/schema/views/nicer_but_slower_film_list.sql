create view nicer_but_slower_film_list as
select
    film.film_id as FID,
    film.title,
    film.description,
    category.name as category,
    film.rental_rate as price,
    film.length,
    film.rating,
    string_agg(
        upper(substring(actor.first_name, 1, 1)) || 
        lower(substring(actor.first_name, 2, length(actor.first_name))) ||
        ' ' ||
        upper(substring(actor.last_name, 1, 1)) || 
        lower(substring(actor.last_name, 2, length(actor.last_name)))
        , ', '
    ) as actors
from category
    left join film_category on category.category_id = film_category.category_id
    left join film on film.film_id = film_category.film_id
    inner join film_actor on film.film_id = film_actor.film_id
    inner join actor on film_actor.actor_id = actor.actor_id
group by film.film_id;