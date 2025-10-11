
run:
	go run main.go
run-auth:
	go run main.go auth


start:
	sudo docker compose -f compose.dev.yaml up

stop:
	sudo docker compose -f compose.dev.yaml down

start_prod:
	sudo docker compose up

stop_prod:
	sudo docker compose down
