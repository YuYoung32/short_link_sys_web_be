use short_link_sys;
create table short_link_sys.links
(
    short_link  varchar(255) collate utf8mb3_bin not null
        primary key,
    long_link   longtext                         null,
    create_time bigint                           null,
    update_time bigint                           null,
    comment     longtext                         null
);

create table short_link_sys.visits
(
    short_link varchar(255) collate utf8mb3_bin null,
    ip         longtext                         null,
    region     longtext                         null,
    visit_time bigint                           null
);

create view link_visits as
select links.short_link, links.long_link, links.comment, visits.ip, visits.region, visits.visit_time
from links
         inner join visits on links.short_link = visits.short_link