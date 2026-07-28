importdb:
	docker cp .\backupDb.sql postgres-db:/backupDb.sql
	docker exec postgres-db psql -U root -d kienDb -f /backupDb.sql
exportdb:
	docker exec postgres-db sh -c "pg_dump -U root -d kienDb > /tmp/backupDb.sql"
	docker cp postgres-db:/tmp/backupDb.sql .\backupDb.sql
server:
	go run .