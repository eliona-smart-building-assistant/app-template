FROM golang:1.27.3-alpine3.23 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG GIT_COMMIT=unknown
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X 'app-name/internal/api.BuildTimestamp=$(date -u +%Y-%m-%dT%H:%M:%SZ)' -X 'app-name/internal/api.GitCommit=${GIT_COMMIT}'" \
    -o /out/app ./cmd/app-name

FROM alpine:3.23.8

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H -u 1000 app \
    && mkdir /data && chown app:app /data

COPY --from=build /out/app /app

USER 1000
ENV DATA_DIR=/data
VOLUME /data
EXPOSE 3000 5000
ENTRYPOINT ["/app"]
