all:
	gulp css scripts concat-css
	go run main.go api.go analytics.go parsegame.go

startPm2:
	sudo pm2 start -x ./lol-counter -n "go"

dev:
	gulp css scripts concat-css
	gin -a 8080 run

matchups:
	cd scripts/
	go run *.go

scrape:
	go run scrape.go api.go analytics.go parsegame.go

scrapewithlog:
	go run scrape.go api.go analytics.go parsegame.go > "scrape.$(date +%Y-%m-%d_%H:%M:%S).log" 2>&1
