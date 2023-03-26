create view sales_by_store as
select
    city.city || ', ' || country.country as store,
    staff.first_name || ' ' || staff.last_name as manager,
    sum(payment.amount) as total_sales
from payment
    inner join rental on rental.rental_id = payment.rental_id
    inner join inventory on rental.inventory_id = inventory.inventory_id
    inner join store on inventory.store_id = store.store_id
    inner join address on store.address_id = address.address_id
    inner join city on address.city_id = city.city_id
    inner join country on city.country_id = country.country_id
    inner join staff on store.manager_staff_id = staff.staff_id
group by store.store_id
order by country.country, city.city;
