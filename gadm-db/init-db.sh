#!/usr/bin/env bash
set -e

# Check for required environment variables
REQUIRED_VARS=(
  POSTGRES_USER
  POSTGRES_DB
  GADM_DB_NAME
  ADMIN_GROUP
  ADMIN_USER
  ADMIN_PASSWORD
  APP_GROUP
  APP_USER
  APP_PASSWORD
  POSTGRES_SUPERUSER_PASSWORD
)

missing_vars=()
for var in "${REQUIRED_VARS[@]}"; do
  if [[ -z "${!var}" ]]; then
    missing_vars+=("$var")
  fi
done

if (( ${#missing_vars[@]} )); then
  echo "Error: The following required environment variables are missing:"
  for var in "${missing_vars[@]}"; do
    echo "  $var"
  done
  exit 1
fi

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE "$GADM_DB_NAME";
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$GADM_DB_NAME" <<-EOSQL2
    -- Declare PostgreSQL variables from environment
    \set admin_group '${ADMIN_GROUP}'
    \set admin_user '${ADMIN_USER}'
    \set admin_password '${ADMIN_PASSWORD}'
    \set app_group '${APP_GROUP}'
    \set app_user '${APP_USER}'
    \set app_password '${APP_PASSWORD}'

    BEGIN ISOLATION LEVEL SERIALIZABLE;

    CREATE ROLE :admin_group;
    CREATE ROLE :admin_user
        WITH LOGIN PASSWORD :'admin_password'
        IN ROLE :admin_group;

    CREATE ROLE :app_group;
    CREATE ROLE :app_user
        WITH LOGIN PASSWORD :'app_password'
        IN ROLE :app_group;

    CREATE SCHEMA gadm AUTHORIZATION :admin_user; 
    CREATE SCHEMA postgis  AUTHORIZATION :admin_user; 
    CREATE SCHEMA common AUTHORIZATION :admin_user; 

    ALTER DATABASE gadm SET search_path TO common, gadm, postgis, public;

    --  grant all schemas to {APP_GROUP}
    GRANT ALL ON SCHEMA gadm TO :app_group ;
    GRANT ALL ON SCHEMA postgis TO :app_group;
    GRANT ALL ON SCHEMA common TO :app_group;

    -- db access for admin and app_user
    GRANT ALL ON DATABASE gadm to :admin_group;
    GRANT ALL ON DATABASE gadm to :app_group;

    -- schema access to app_user
    ALTER DEFAULT PRIVILEGES FOR ROLE :admin_user GRANT ALL ON SCHEMAS TO :app_group;

    -- allow app_user on gadm schema created by admin
    ALTER DEFAULT PRIVILEGES FOR ROLE :admin_user IN SCHEMA gadm GRANT SELECT ON TABLES TO :app_group;
    ALTER DEFAULT PRIVILEGES FOR ROLE :admin_user IN SCHEMA gadm GRANT SELECT ON SEQUENCES TO :app_group;

    -- TODO: allow app_user on common schema created by admin
    ALTER DEFAULT PRIVILEGES FOR ROLE :admin_user IN SCHEMA common GRANT ALL ON TABLES TO :app_group;
    ALTER DEFAULT PRIVILEGES FOR ROLE :admin_user IN SCHEMA common GRANT ALL ON SEQUENCES TO :app_group;
    ALTER DEFAULT PRIVILEGES FOR ROLE :admin_user IN SCHEMA common GRANT ALL ON ROUTINES TO :app_group;
    ALTER DEFAULT PRIVILEGES FOR ROLE :admin_user IN SCHEMA common GRANT USAGE ON TYPES TO :app_group;

    -- allow app_user on postgis schema created by admin
    ALTER DEFAULT PRIVILEGES FOR ROLE :admin_user IN SCHEMA postgis 
        GRANT EXECUTE ON ROUTINES TO :app_group;
    ALTER DEFAULT PRIVILEGES FOR ROLE :admin_user IN SCHEMA postgis 
        GRANT USAGE ON TYPES TO :app_group;

    CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA postgis;
    COMMIT;
EOSQL2

# Secure the postgres superuser
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$GADM_DB_NAME" <<-EOSQL3
    -- Declare postgres password variable for this session
    \set postgres_password '${POSTGRES_SUPERUSER_PASSWORD}'
    
    -- Set a strong password for postgres superuser
    ALTER ROLE postgres WITH PASSWORD :'postgres_password';
    
    -- Limit postgres user to only 1 connection at a time
    ALTER ROLE postgres CONNECTION LIMIT 1;
    
    -- Connect to gadm database and drop postgres database
    DROP DATABASE IF EXISTS postgres;
EOSQL3


