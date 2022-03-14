ARG BUILDER_IMAGE=golang:alpine
ARG DISTROLESS_IMAGE=gcr.io/distroless/static

# Stage I : Build the builder image
FROM ${BUILDER_IMAGE} AS builder

# Install GIT and CA-CERTIFICATES
RUN apk update && \
    apk add --no-cache git ca-certificates && \
    update-ca-certificates

# Set the workdir and copy over the sources
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

# Add custom CA
ADD ./fortigate.crt /usr/local/share/ca-certificates/fortigate.crt
RUN update-ca-certificates

# Fetch dependencies
#RUN go get -d -v
RUN go mod download
RUN go mod verify

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/server .

# Build the runtime image
# using static non-root image
# user:group = nobody:nobody, uid:gid = 65534:65534
FROM ${DISTROLESS_IMAGE}

# Copy the binary
COPY --from=builder /go/bin/server /go/bin/server

COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Run the server
CMD ["/go/bin/server"]
