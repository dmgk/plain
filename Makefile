APPID=		plain-im
VERSION!=	date +%Y%m%d-%H%M%S
SOURCES!=	find . -type f -name \*.go

plain: ${SOURCES}
	go build ./...

server: plain
	goapp serve app

deploy:
	# goapp deploy -application ${APPID} -version ${VERSION} app
	gcloud app deploy app

clean:
	-rm -f ./plain
