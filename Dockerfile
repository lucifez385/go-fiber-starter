FROM golang:1.20.2-alpine as deps

RUN apk update && \
        apk upgrade && \
        apk add --no-cache \
        ca-certificates \
        tini \
        tzdata
ENV TERM=xterm
ENV TZ=Asia/Bangkok

RUN mkdir -p /app
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod tidy
RUN go mod download
COPY . .

FROM deps AS dev
RUN go install github.com/cosmtrek/air@latest
CMD ["air", "-c", ".air.toml"]

# FROM deps AS lint
# RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.44.2 && \
#     golangci-lint run --out-format html --issues-exit-code 0 > lint-report.html

# FROM deps AS test
# RUN CGO_ENABLED=0 go test ./... -v -cover -coverprofile=coverage.out

FROM deps AS build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.14 as release
RUN apk update && \
        apk upgrade && \
        apk add --no-cache \
        ca-certificates \
        tini \
        tzdata
ENV TERM=xterm
ENV TZ=Asia/Bangkok

WORKDIR /app
RUN mkdir -p tmp
COPY --from=build /app/app ./
ENTRYPOINT ["./app"]
