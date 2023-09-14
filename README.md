# projectWithNats

for start this project 
1) start nats-streaming in docker: docker run --name nats --rm -p <port:port> -p <port:port> nats --http_port <port>
2) start postgresql in docker: docker run --name postgresql -p <port:port> -e POSTGRES_USER=<user> -e POSTGRES_PASSWORD=<password> -e POSTGRES_DB=<nameDB> -d postgres:13.3
3) Create table with code from data.sql
4) located in the root directory of the project start in terminal: go run cmd/main.go
