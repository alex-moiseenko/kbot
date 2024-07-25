FROM golang:1.20 AS builder

WORKDIR /go/src/app
COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=arm64
RUN make TARGETOS=${TARGETOS} TARGETARCH=${TARGETARCH} build

FROM scratch
WORKDIR /
ARG TARGETOS=linux
ARG TARGETARCH=arm64
COPY --from=builder /go/src/app/kbot-${TARGETOS}-${TARGETARCH} /kbot
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./kbot"]