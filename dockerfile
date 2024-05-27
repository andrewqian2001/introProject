# Use the official Golang image as the base image so we dont need to install go manually in the container
FROM golang:1.22.3

# Set the working directory inside the container, any copy/run command will be executed relative to this container
WORKDIR /app

# Copy the go.mod and go.sum files into the working directory, these files define our dependencies
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application, the -o main flag specifies that the output binary should be named main. The . indicates that the source files are in the current directory 
RUN go build -o main .

# Expose the port the application runs on
EXPOSE 8090

# Run the executable
CMD ["./main"]


# use command: docker build -t todo-app .
# This command builds the Docker image and tags it as todo-app.
# Run the Docker Container: docker run -d -p 8090:8090 -v $(pwd)/app.log:/introproject/app.log todo-app
# This command runs the Docker container in detached mode (-d), 
# maps port 8090 of the container to port 8090 on the host (-p 8090:8090), and mounts the app.log file from the current directory to the container (-v $(pwd)/app.log:/introproject/app.log).
#


