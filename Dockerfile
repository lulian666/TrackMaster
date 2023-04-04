FROM golang:1.20-alpine@sha256:48f336ef8366b9d6246293e3047259d0f614ee167db1869bdbc343d6e09aed8a AS builder

RUN apk add --update --no-cache git make
WORKDIR /app
COPY . .
RUN make build



FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/TrackMaster /app/TrackMaster
COPY --from=builder /app/entrypoint.sh /app/entrypoint.sh

EXPOSE 8000

# Set the entrypoint
#CMD ["/app/TrackMaster"]
ENTRYPOINT ./entrypoint.sh
