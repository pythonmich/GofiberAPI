server:
	go run main.go

migrate-up:
	migrate -path database/migrations -database ${PQURL} -verbose up

migrate-down:
	migrate -path database/migrations -database ${PQURL} -verbose down

migrate-up1:
	migrate -path database/migrations -database ${PQURL} -verbose up 1

migrate-down1:
	migrate -path database/migrations -database ${PQURL} -verbose down 1

migrate-up2:
	migrate -path database/migrations -database ${PQURL} -verbose up 3

migrate-down2:
	migrate -path database/migrations -database ${PQURL} -verbose down 3


migrate-dirty:
	migrate -path database/migrations -database ${PQURL} force 3

#migrate create -ext sql -dir database/migrations -seq session_schema

.PHONY: server migrate-down migrate-up migrate-dirty migrate-up1 migrate-down1 migrate-down2 migrate-up2