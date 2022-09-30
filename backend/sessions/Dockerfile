FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY configs/ ./configs
COPY internal/ ./internal
COPY pkg/ ./pkg

RUN go build -o /sessions ./cmd/server/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /sessions /sessions

EXPOSE 8060

USER nonroot:nonroot

COPY configs/ ./configs

ENTRYPOINT ["/sessions"]