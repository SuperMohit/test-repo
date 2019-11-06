FROM golang:1.12 as builder

LABEL maintainer="Yugen Maintainer"

# Copy local code to the container image.
RUN mkdir -p src
COPY src src
WORKDIR src

RUN go get -d -v
# Run unit tests
RUN CGO_ENABLED=0 GOOS=linux go test 

# Build the code inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o yugen-go


# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM photon:3.0

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/yugen-go yugen-go

EXPOSE 8080

# Run the web service on container startup.
CMD ["/yugen-go"]