# gingonic-tournament
## Project Structure

yasin-cicd/
├── cmd/
│   └── app/
│       └── main.go
├── bin/
│   ├── app
├── internal/
│   ├── crud/
│   │   ├── leaderboard.go
│   │   ├── tournament.go
│   │   ├── users.go
│   ├── db/
│   │   ├── postgres.go
│   │   ├── redis.go
│   ├── router/
│   │   ├── leaderboard.go
│   │   ├── tournament.go
│   │   ├── users.go
├── model/
│   ├── leaderboard.go
│   ├── tournament.go
│   ├── user.go
├── service/
│   ├── leaderboard.go
│   ├── tournament.go
│   ├── users.go
├── validation/
│   ├── leaderboard.go
│   ├── tournament.go
│   ├── users.go
├── build/
│   ├── Dockerfile
├── deployments/
│   ├── docker-compose.yaml
├── scripts/
│   ├── build.sh
│   ├── test.sh
├── vendor/
├── test/
│   ├── db_test.go
├── go.mod
├── go.sum
├── .env
├── .gitlab-ci.yml


## Swagger Documentation
swag init -g cmd/app/main.go -o docs
http://localhost:8080/swagger/index.html

## Package management go commands:

go mod init yasin-cicd-app
go mod tidy

## Testing

go test -v ./...
docker compose --profile test up --build

## create local docker compose
docker compose --profile default up --build
docker compose --profile test up --build

### get in into the container
docker ps -a
docker exec -it <id-of-running-container> /bin/bash
exit

### remove leftover images and useless db's:
docker compose --profile default down --remove-orphans
docker rmi $(docker images -q --filter "dangling=true") && docker rm -f $(docker ps -a -q --filter "status=exited")
docker rm -f $(docker ps -a -q --filter "status=exited")
docker ps -a
docker rm -f <id of postgres and redis if you run local>

## API Documentation

next level: https://swagger.io/docs/

http://10.0.1.10:8080/swagger/index.html

Base URL
http://localhost:8080/
http://10.0.1.10:8080/

## example curl commands:

curl -s http://localhost:8080/health && break || sleep 5

curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Alice","money":100, "level":1 }' 
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Frank","money":200, "level":1}' \
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Yasin","money":250, "level":3}' \
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Batu","money":150, "level":2}' \
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Meriç","money":199, "level":3}' \
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Kassandra","money":500, "level":4}'



curl http://localhost:8080/leaderboard

curl -X PUT http://localhost:8080/user -H "Content-Type: application/json" -d '{"name":"Alice","money":200}'
curl -X GET http://localhost:8080/user/Alice

curl -X DELETE http://localhost:8080/user/Alice
curl http://localhost:8080/user/Alice
curl http://localhost:8080/leaderboard

### To see ip adresses of containers: 

docker inspect -f '{{.Name}} - IP: {{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}} - Ports: {{.NetworkSettings.Ports}}' $(docker ps -q)

output: 
/yasin-cicd-app-1 - IP: 10.0.1.1010.0.2.10 - Ports: map[8080/tcp:[{0.0.0.0 8080} {:: 8080}]]
/yasin-cicd-postgres-1 - IP: 10.0.2.11 - Ports: map[5432/tcp:[{0.0.0.0 5432} {:: 5432}]]
/yasin-cicd-redis-1 - IP: 10.0.2.12 - Ports: map[6379/tcp:[{0.0.0.0 6379} {:: 6379}]]
/gitlab-runner - IP: 172.17.0.2 - Ports: map[]

app container can still communicate with the postgres and redis containers by using service names as hostnames, 
which Docker Compose automatically sets up as DNS entries for the containers in the same network.
