all:
	gulp css scripts concat-css
	go run main.go api.go analytics.go parsegame.go

dev:
	gulp css scripts concat-css
	gin -a 8080 run

matchups:
	cd scripts/
	go run *.go
