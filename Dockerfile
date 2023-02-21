FROM docker.io/golang:1.19-alpine AS base

RUN apk --update --no-cache add bash build-base

WORKDIR /build
COPY . /build

RUN mkdir -p bin
RUN go build -trimpath -o bin/ ./cmd/dendrite-monolith-server
RUN go build -trimpath -o bin/ ./cmd/create-account
RUN go build -trimpath -o bin/ ./cmd/generate-keys
RUN apk --update --no-cache add curl

FROM alpine:latest
LABEL org.opencontainers.image.title="Dendrite (Monolith)"
LABEL org.opencontainers.image.description="Next-generation Matrix homeserver written in Go"
LABEL org.opencontainers.image.source="https://github.com/matrix-org/dendrite"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.vendor="The Matrix.org Foundation C.I.C."
LABEL org.opencontainers.image.documentation="https://matrix-org.github.io/dendrite/"

COPY --from=base /build/bin/* /usr/bin/

VOLUME /etc/dendrite
WORKDIR /etc/dendrite

ENTRYPOINT ["/usr/bin/dendrite-monolith-server"]
