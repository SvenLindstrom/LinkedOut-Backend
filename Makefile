
run:
	go run main.go
run-auth:
	go run main.go auth

start_test:
	sudo docker compose -f compose.test.yaml up

stop_test:
	sudo docker compose -f compose.test.yaml down

start:
	sudo docker compose -f compose.dev.yaml up

stop:
	sudo docker compose -f compose.dev.yaml down

start_prod:
	sudo docker compose up

stop_prod:
	sudo docker compose down
