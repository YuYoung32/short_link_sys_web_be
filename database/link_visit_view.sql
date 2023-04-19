create view link_visits as
select links.short_link, links.long_link, links.comment, visits.ip, visits.region, visits.visit_time
from links
         inner join visits on links.short_link = visits.short_link