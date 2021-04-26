FROM golang:latest

COPY ./go.mod $GOPATH/dockerBuild/go.mod
COPY ./go.sum $GOPATH/dockerBuild/go.sum
WORKDIR $GOPATH/dockerBuild
RUN go mod download
COPY . .
RUN go build main.go
CMD [ "./main" ]