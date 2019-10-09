BINARY = main
GOARCH = amd64

region?=eu-west-1
accountid?=012345678

all: build

build:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY} *.go

deploy:
	zip main.zip main
	aws s3 cp main.zip s3://lambda-scripts-${region}-${accountid}/compliance-enforcer/main.zip

clean:
	rm -f ${BINARY}k
