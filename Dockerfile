FROM golang:1.14 AS build
WORKDIR /go/src/github.com/duyet/applause-btn
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

FROM alpine:latest
WORKDIR /app
RUN mkdir ./public
COPY ./public ./public
COPY --from=build /go/src/github.com/duyet/applause-btn/app .
EXPOSE 3000
CMD ["./app"]