server:
	go run main.go

migrate-up:
	migrate -path database/migrations -database ${PQURL} -verbose up

migrate-down:
	migrate -path database/migrations -database ${PQURL} -verbose down
migrateDirtyVersion:
	migrate -path database/migrations -database ${PQURL} force 1

.PHONY: server migrate-down migrate-up migrateDirtyVersion