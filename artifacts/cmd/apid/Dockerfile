
# syntax=docker/dockerfile:1

# Super basic docker setup for now as it works well enough...
# The sizing is out of whack but thats what we get to ignore doing anything reasonable with pulling only intentional things.
FROM golang:1.17-alpine as base


# this name is cute, it exists for no other reason than that.
WORKDIR /appdemon
COPY go.mod ./
COPY go.sum ./

COPY . ./
# Get all the deps so that things dont break 
RUN go mod download



RUN go build -o appd ./internal/apid/cmd

# Pretend like we are smart and create an actual image to publish
FROM scratch
# decide to throw away ~all our gains by just blindly copying stuff. Dont do this in the future.
COPY --from=base ./ /
EXPOSE 8080

ENTRYPOINT [ "/appdemon/appd" ]