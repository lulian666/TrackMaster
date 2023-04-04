FROM golang:1.20-alpine@sha256:48f336ef8366b9d6246293e3047259d0f614ee167db1869bdbc343d6e09aed8a

# Create folder /app and non-privileged user as root
RUN mkdir /app && \
    adduser -S app-user

# Copy the project to the /app folder
COPY . /app

# Set the current folder as /app
WORKDIR /app

# Build the app
RUN go get ./... && \
    go build -o ./TrackMaster

# Change owner
RUN chown -R app-user /app

# Use the non-privileged user for next actions
USER app-user

EXPOSE 8000 80

# Set the entrypoint
ENTRYPOINT ./entrypoint.sh
