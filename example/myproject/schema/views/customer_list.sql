create view customer_list as
select
    cu.customer_id as id,
    cu.first_name || ' ' || cu.last_name as name,
    a.address,
    a.postal_code as "zip code",
    a.phone,
    city.city,
    country.country,
    case
        when cu.active then 'active'
        else ''
    end as notes,
    cu.store_id as SID
from customer as cu
    join address on cu.address_id = address.address_id
    join city on address.city_id = city.city_id
    join country on city.country_id = country.country_id;