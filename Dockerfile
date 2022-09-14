FROM golang:1.18-buster as build

WORKDIR ./src

COPY go.mod go.sum main.go ./
COPY provenance ./provenance

RUN go mod download && go build -o /out/generator

FROM ubuntu:18.04
COPY --from=build /out/generator /generator
RUN mkdir /artifacts/
# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/generator"]
