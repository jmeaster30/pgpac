create view film_list as
select
    film.film_id as FID,
    film.title,
    film.description,
    category.name as category,
    film.rental_rate as price,
    film.length,
    film.rating,
    string_agg(actor.first_name || ' ' || actor.last_name, ', ') as actors
from category
    left join film_category on category.category_id = film_category.category_id
    left join film on film_category.film_id = film.film_id
    inner join film_actor on film.film_id = film_actor.film_id
    inner join actor on film_actor.actor_id = actor.actor_id
group by film.film_id; -- I think this group by wouldn't work