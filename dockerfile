# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Make the build script executable
RUN chmod +x build.sh

# Run the build script
RUN ./build.sh

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `build.sh`
CMD ["./bin/rabbit"]