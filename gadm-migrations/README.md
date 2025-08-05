# GADM API Migration System

## How to create new migration

```
$ cd <PROJECT_HOME>/gadm-migrations
$ goose create test sql
```

## How to apply migrations

### Apply UP migrations (forward)

```bash
$ cd <PROJECT_HOME>
$ docker compose -f docker-compose.test.yml build gadm-migrations
$ docker compose -f docker-compose.test.yml up gadm-migrations
```

### Apply DOWN migrations (rollback)

```bash
$ cd <PROJECT_HOME>
$ docker compose -f docker-compose.test.yml build gadm-migrations
$ docker compose -f docker-compose.test.yml run --rm gadm-migrations down
```
