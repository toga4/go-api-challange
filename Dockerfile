FROM golang:1.17 as builder

WORKDIR /go/src/app
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base:debug
COPY --from=builder /go/bin/app /
ENTRYPOINT ["/app"]
