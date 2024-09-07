FROM --platform=$TARGETPLATFORM golang:1.22.3

WORKDIR /app

ARG TARGET
ARG BUILD_DATE
ENV BUILD_DATE=$BUILD_DATE

COPY go.sum ./
COPY go.mod ./
RUN go mod download

COPY . .

ENV GOARCH=arm
ENV GOOS=linux

RUN go build -o /main $TARGET

ENTRYPOINT ["/main"]
