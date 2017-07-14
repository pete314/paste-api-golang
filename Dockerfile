FROM golang:latest
RUN apt-get update && apt-get install -y \
    aufs-tools \
    automake \
    build-essential

# Copy the API Source
RUN mkdir /api
ADD ./api /api/
WORKDIR /api

# Copy config
RUN mkdir /config
ADD ./config /config

# Compile the source
RUN make

# Expose the application on port 8080
EXPOSE 8080

# Create an entry point
ENTRYPOINT ["/api/paste-api-golang"]

# Create command entry point for app
CMD ["/api/paste-api-golang", "-conf", "config/server-config.local.json"]
