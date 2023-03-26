create view sales_by_film_category as
select
   category.name as category,
    sum(payment.amount) as total_sales
from payment
    inner join rental on payment.rental_id = rental.rental_id
    inner join inventory on rental.inventory_id = inventory.inventory_id
    inner join film on inventory.film_id = film.film_id
    inner join film_category on film.film_id = film_category.film_id
    inner join category on film_category.category_id = category.category_id
group by category.name
order by total_sales desc;
