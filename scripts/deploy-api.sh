## Author: Peter Nagy <https://peternagy.ie>
## Since: June, 2017
## Description: Deploy file to build production container

# Create deploy space
cd .. && rm -rf deploy/api

# Compile the source
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./deploy/api api/src/runner.go

# Move binary into container
cp config/server-config.local.json deploy/

# Create production image
cd deploy
sudo docker build --no-cache=true -t paste-api-scratch -f Dockerfile.scratch .
