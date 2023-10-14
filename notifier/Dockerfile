# Use an official Go runtime as a parent image
FROM golang:1.21.3

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
ADD . /app

# Install any needed dependencies
RUN go mod tidy && \
    go mod download && \
    go build -o notifier .

# Make port 8080 available to the world outside this container
EXPOSE 8080

# Run the binary when the container launches
CMD ["./notifier"]
