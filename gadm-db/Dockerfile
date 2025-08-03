FROM postgis/postgis:16-3.5

# Copy the initialization script and template
COPY init-db.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/init-db.sh

# Copy organized HBA configuration files
COPY hba-security.conf /docker-entrypoint-initdb.d/
COPY setup-pg-hba.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/setup-pg-hba.sh

# Copy PostgreSQL memory optimization script
COPY postgres-memory-config.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/postgres-memory-config.sh