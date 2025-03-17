FROM golang:1.24-alpine3.21 AS builder

RUN apk add --no-cache make git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make bin/agentruntime

FROM alpine:3.21 AS runner

RUN apk add --no-cache ca-certificates

WORKDIR /app
RUN mkdir /app/agents && chown -R 1000:1000 /app
USER 1000:1000
COPY --from=builder --chown=1000:1000 /app/bin/agentruntime /app/agentruntime
COPY --from=builder --chown=1000:1000 /app/examples /app/examples

ENTRYPOINT ["./agentruntime", "serve"]