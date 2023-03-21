create view staff_list as
select
    s.staff_id as id,
    s.first_name || ' ' || s.last_name as name,
    a.address as address,
    a.postal_code as "zip code",
    a.phone,
    city.city,
    country.country,
    s.store_id as SID
from staff s
    join address a on s.address_id = a.address_id
    join city on a.city_id = city.city_id
    join country on city.country_id = country.country_id;