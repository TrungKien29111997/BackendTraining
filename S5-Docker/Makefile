include .env
export

CONN_STRING = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
MIGRATION_DIRS = internal/db/migrations

importdb:
	docker cp .\backupDb.sql postgres-db:/backupDb.sql
	docker exec postgres-db psql -U root -d kienDb -f /backupDb.sql
exportdb:
	docker exec postgres-db sh -c "pg_dump -U root -d kienDb > /tmp/backupDb.sql"
	docker cp postgres-db:/tmp/backupDb.sql .\backupDb.sql
server:
	go run .

# Create migration (make migrate_create NAME=profiles)
migrate_create:
	migrate create -ext sql -dir $(MIGRATION_DIRS) -seq $(NAME)
migrate_up:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" up
migrate_down:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down 1
mirgate_force:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" force $(VERSION)
# Delete all migration (not recommend)
migrate_drop:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" drop
migrate_goto:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" goto $(VERSION)

.PHONY: importdb exportdb server migrate_create migrate_up migrate_down