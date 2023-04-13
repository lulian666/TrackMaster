FROM registry.jellow.site/iftech/golang:1.18

WORKDIR /app
COPY . .
RUN make build

EXPOSE 8000

#CMD ["/app/TrackMaster"]

ENTRYPOINT ./entrypoint.sh