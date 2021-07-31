server:
	go run main.go

postgres:
	docker run --name postgres-finance -p 5432:5432 -e POSTGRES_USER=python_mich -e POSTGRES_PASSWORD=Musyimi7. -d postgres:latest

create-db:
	docker exec -it postgres-finance createdb --username=python_mich --owner=python_mich go_financeapi

drop-db:
	docker exec -it postgres-finance dropdb go_financeapi


migrate-up:
	migrate -path database/migrations -database ${PQURL} -verbose up

migrate-down:
	migrate -path database/migrations -database ${PQURL} -verbose down

migrate-up1:
	migrate -path database/migrations -database ${PQURL} -verbose up 1

migrate-down1:
	migrate -path database/migrations -database ${PQURL} -verbose down 1

migrate-up2:
	migrate -path database/migrations -database ${PQURL} -verbose up 2

migrate-down2:
	migrate -path database/migrations -database ${PQURL} -verbose down 2

migrate-up3:
	migrate -path database/migrations -database ${PQURL} -verbose up 3

migrate-down3:
	migrate -path database/migrations -database ${PQURL} -verbose down 3

migrate-up4:
	migrate -path database/migrations -database ${PQURL} -verbose up 4

migrate-down4:
	migrate -path database/migrations -database ${PQURL} -verbose down 4

migrate-up5:
	migrate -path database/migrations -database ${PQURL} -verbose up 5

migrate-down5:
	migrate -path database/migrations -database ${PQURL} -verbose down 5

migrate-up6:
	migrate -path database/migrations -database ${PQURL} -verbose up 6

migrate-down6:
	migrate -path database/migrations -database ${PQURL} -verbose down 6

migrate-dirty:
	migrate -path database/migrations -database ${PQURL} force 7

#migrate create -ext sql -dir database/migrations -seq session_schema

.PHONY: server postgres create-db drop-db migrate-down migrate-up migrate-dirty migrate-up1 migrate-down1 migrate-down2 migrate-up2 migrate-down3 migrate-up3