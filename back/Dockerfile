FROM golang:1.14.9-alpine3.12 as builder

RUN mkdir /app
WORKDIR /app
RUN apk add --update --no-cache ca-certificates git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /bin/seismographd cmd/server/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/seismographd /bin/seismographd
ENTRYPOINT ["/bin/seismographd"]