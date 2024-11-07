# Dockerfile to test the init feature

FROM golang AS build
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /app/kairos-init .

FROM ubuntu:24.04
COPY --from=build /app/kairos-init /kairos-init
RUN /kairos-init -l debug -f all
RUN rm /kairos-init