FROM golang:1.22-alpine AS build
COPY go.mod go.sum .
RUN go mod download 

COPY *.go .
RUN go build -o sshwitter 

FROM alpine:latest AS run
COPY --from=build /go/sshwitter .

EXPOSE 8080
CMD [ "./sshwitter" ]
