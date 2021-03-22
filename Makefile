install: 
	go get -u github.com/dgrijalva/jwt-go
	go get -u github.com/gin-gonic/gin
	go get -u github.com/jmoiron/sqlx
	go get -u github.com/joho/godotenv
	go get -u github.com/lib/pq
	go get -u github.com/sirupsen/logrus
	go get -u github.com/spf13/viper
rundb:
	docker run --name=windfarms-db -e POSTGRES_PASSWORD='2210' -p 5436:5432 -d --rm postgres
stopdb:
	migrate -path ./schema -database "postgres://postgres:2210@localhost:5436/postgres?sslmode=disable" down
run:
	migrate -path ./schema -database "postgres://postgres:2210@localhost:5436/postgres?sslmode=disable" up
	go run cmd/main.go