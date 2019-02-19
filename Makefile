VERSION !=	date +%Y%m%d-%H%M%S
SOURCES !=	find . -type f -name \*.go

plain: ${SOURCES}
	go build

server: plain
	./plain

deploy:
	gcloud app deploy

clean:
	-rm -f ./plain
