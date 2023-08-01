FROM golang:latest
LABEL authors="Bradlee Johnson"



WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


# creates working directory
RUN mkdir -p /app

# Set working directory
WORKDIR /app

# Copy Local files
ADD . /app


EXPOSE 8000
# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o /panda
#RUN go build -o app ./main.go

# Execute when container is started
CMD ["/panda"]