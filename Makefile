.PHONY:
.SILENT:

build-image:
	docker build -t gdrive-bot-im .

start-container:
	docker run --name gdrive-bot -p 8080:8080 --env-file .env gdrive-bot-im
