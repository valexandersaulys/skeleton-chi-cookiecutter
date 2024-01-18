FROM golang:alpine AS builder

RUN apk update && \
    apk add --no-cache make gcc g++

WORKDIR /src
COPY go.mod /src
COPY go.sum /src
RUN go mod download

# Build time environmental values
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

COPY . /src
RUN make build

FROM alpine:latest AS runtime
COPY --from=builder /src/_build/{{cookiecutter.application_name}}.bin /usr/bin/{{cookiecutter.application_name}}.bin

# Runtime environmental values -- these don't affect build
ENV CSRF_PROTECTION_KEY="1234"
ENV RUNTIME_ENV=LOCAL
ENV ADD_DUMMIES=false

ENTRYPOINT ["/usr/bin/output.bin", "--log-level", "INFO", "--log-path", "/tmp/chi.logs", "--port", "2999"]
