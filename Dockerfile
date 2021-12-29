# Building the binary of the App
FROM golang:1.15 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /go/src/ethereum-auth

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest

LABEL owner_team OPS

HEALTHCHECK --start-period=30s --retries=1 --interval=4s CMD curl --max-time 3 --connect-timeout 2 -sSf http://127.0.0.1:3030/health || exit 1

WORKDIR /app

# Create the `public` dir and copy all the assets into it
# RUN mkdir ./static
# COPY ./static ./static

# `boilerplate` should be replaced here as well
COPY --from=build /go/src/ethereum-auth/app .

# Exposes port 3030 because our program listens on that port
EXPOSE 3030

CMD ["./app"]