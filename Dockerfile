# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x && go mod verify

ARG TARGETOS
ARG TARGETARCH

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o app-binary .

FROM gcr.io/distroless/static AS final

WORKDIR /usr/src/app

USER nobody:nobody

COPY --from=build /usr/src/app/app-binary /usr/bin/server

EXPOSE 8080

ENTRYPOINT [ "/usr/bin/server" ]
