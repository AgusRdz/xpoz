FROM golang:1.24-alpine

RUN apk add --no-cache git ca-certificates curl && \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
      | sh -s -- -b "$(go env GOPATH)/bin" v1.64.8

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

COPY . .

CMD ["go", "build", "-o", "bin/xpoz", "./cmd/xpoz"]
