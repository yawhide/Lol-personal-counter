all:
	go run main.go api.go analytics.go parsegame.go

dev:
	gin -a 8080 run

matchups:
	cd scripts/
	go run *.go
