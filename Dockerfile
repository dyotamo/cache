FROM golang:1.21.4-bullseye AS builder
WORKDIR /builder
COPY go.mod go.mod
COPY cmd/ cmd/
RUN go mod tidy && CGO_ENABLED=0 go build -o app cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /builder/app app
CMD [ "./app" ]