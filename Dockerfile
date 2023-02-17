FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Set environment variables
ENV DB_HOST=$DB_HOST \
    DB_PORT=$DB_PORT \
    DB_USER=$DB_USER \
    DB_PASSWORD=$DB_PASSWORD \
    DB_NAME=$DB_NAME \
    SALT=$SALT \
    JWT_SECRET=$JWT_SECRET \
    REDIS_SECRET=$REDIS_SECRET \
    REDIS_HOST=$REDIS_HOST \
    REDIS_PASSWORD=$REDIS_PASSWORD


# Build the Go API
RUN go build -o main .

# Expose port 8080 for the API
EXPOSE 8080

# Run the Go API when the container launches
CMD ["./main"]
