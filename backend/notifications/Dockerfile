FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY configs/ ./configs
COPY internal/ ./internal
COPY pkg/ ./pkg

RUN go build -o /notifications ./cmd/server/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /notifications /notifications

EXPOSE 8080

USER nonroot:nonroot

COPY configs/ ./configs

ENTRYPOINT ["/notifications"]