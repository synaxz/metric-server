# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.22.3 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy local code to the container image.
COPY . .

# Build the binary.
# -o specifies the output file
# By default, it will be placed in the current directory, making it easy to run.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o metricmemory

# Use the official lightweight Scratch image.
# It's essentially a blank slate with no operating system.
FROM scratch

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/metricmemory /metricmemory

EXPOSE 8080

# Run the web service on container startup.
ENTRYPOINT ["/metricmemory"]

