FROM --platform=$BUILDPLATFORM golang:1.25.1 AS deps
WORKDIR /src

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

FROM deps AS builder
ARG TARGET

COPY . .

RUN mkdir -p /out && go build -trimpath -ldflags "-s -w" -o /out/main "$TARGET"

FROM gcr.io/distroless/static-debian12:nonroot AS runtime
ARG BUILD_DATE

ENV BUILD_DATE=$BUILD_DATE

COPY --from=builder /out/main /main

ENTRYPOINT ["/main"]
