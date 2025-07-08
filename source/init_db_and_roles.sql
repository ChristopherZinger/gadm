BEGIN ISOLATION LEVEL SERIALIZABLE;
-- database: gadm
-- schemas: common, gadm, postgis

CREATE ROLE admin_group;
CREATE ROLE  admin 
    WITH LOGIN PASSWORD 'secret'
    IN ROLE admin_group;

CREATE ROLE app_group;
CREATE ROLE app_user
    WITH LOGIN PASSWORD 'secret'
    IN ROLE app_group;


-- create schemas as admin
CREATE SCHEMA gadm AUTHORIZATION admin; 
CREATE SCHEMA postgis  AUTHORIZATION admin; 
CREATE SCHEMA common AUTHORIZATION admin; 

ALTER DATABASE gadm SET search_path TO common, gadm, postgis, public;

--  grant all schemas to app_group
GRANT ALL ON SCHEMA gadm TO app_group ;
GRANT ALL ON SCHEMA postgis TO app_group;
GRANT ALL ON SCHEMA common TO app_group;

-- db access for admin and app_user
GRANT ALL ON DATABASE gadm to admin_group;
GRANT ALL ON DATABASE gadm to app_group;

-- schema access to app_user
ALTER DEFAULT PRIVILEGES FOR ROLE admin GRANT ALL ON SCHEMAS TO app_group;

-- allow app_user on gadm schema created by admin
ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA gadm GRANT SELECT ON TABLES TO app_group;
ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA gadm GRANT SELECT ON SEQUENCES TO app_group;

-- TODO: allow app_user on common schema created by admin
ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA common GRANT ALL ON TABLES TO app_group;
ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA common GRANT ALL ON SEQUENCES TO app_group;
ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA common GRANT ALL ON ROUTINES TO app_group;
ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA common GRANT USAGE ON TYPES TO app_group;


-- allow app_user on postgis schema created by admin
ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA postgis 
    GRANT EXECUTE ON ROUTINES TO app_group;

ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA postgis 
    GRANT USAGE ON TYPES TO app_group;

-- -- allow app_user on admin's schemas
-- ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA postgis 
--     GRANT USAGE ON SCHEMAS TO app_group;


CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA postgis;

-- select * from gadm.test;
-- select * from common.test;


COMMIT;


BEGIN ISOLATION LEVEL SERIALIZABLE;
-- TESTING

-- Set search path for current session
SET search_path TO common, gadm, postgis, public;

set role admin;
select user;

create table gadm.test (
    id serial primary key, 
    name text, 
    geom geometry(MultiPolygon, 4326));

insert into gadm.test (name) values ('test'), ('test2'), ('test3');

create table common.test (id serial primary key, name text);
insert into common.test (name) values ('test'), ('test2'), ('test3');

set role app_user;
select * from gadm.test;
select * from common.test;

insert into common.test (name) values ('created by app_user');

ROLLBACK;
