# # Start from golang base image
# FROM golang:alpine 

# # Install git.
# # Git is required for fetching the dependencies.
# RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# # Setup folders
# RUN mkdir /app
# WORKDIR /app

# # Copy the source from the current directory to the working Directory inside the container
# COPY . .
# COPY .env .

# # Download all the dependencies
# RUN go get -d -v ./...

# # Install the package
# RUN go install -v ./...

# # Build the Go app
# RUN go build -o /build

# # Expose port 8080 to the outside world
# EXPOSE 8080

# # Run the executable
# CMD [ "/build" ]



# Our builder image used to build the Go binary
FROM golang:alpine as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Our production image used to run our app
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache git make musl-dev go
COPY --from=builder /app/main .

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin
EXPOSE 8080
CMD ["./main"]