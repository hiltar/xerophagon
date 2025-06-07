FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY main.go ./
COPY static/ ./static/
COPY templates/ ./templates/

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o xerophagon .

FROM scratch
COPY --from=builder /app/xerophagon /xerophagon

USER 1000:1000
WORKDIR /
VOLUME ["/data"]

CMD ["/xerophagon"]
