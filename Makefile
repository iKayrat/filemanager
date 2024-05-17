# note: call scripts from /scripts
PHONY.: run

run:
	go run cmd/fmsvc/main.go



migrateup:
	migrate -path internal/app/db/migration -database "postgresql://root:root@localhost:5431/fmanagerdb?sslmode=disable" -verbose up
migratedown:
	migrate -path internal/app/db/migration -database "postgresql://root:root@localhost:5431/fmanagerdb?sslmode=disable" -verbose down

createmigrate:
	migrate create -ext sql -dir internal/app/db/migration -seq files

postgres:
	docker run --name fmanager -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:alpine
create:
	docker exec -it fmanager createdb --username=root --owner=root fmanagerdb
drop:
	docker exec -it fmanager dropdb --username=root fmanagerdb
