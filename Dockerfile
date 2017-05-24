# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go

FROM golang

# Copy the local package files to the containers workspace
ADD . /go/src/rest1

# Build the app command inside the container.
# You may fetch or mange dependencies here
# either manually or with a tool like godep
RUN go get -u  github.com/pressly/chi
RUN go install rest1
# Run the app command by default when the container starts
ENTRYPOINT /go/bin/rest1

# Document that the service listens on port 3000
EXPOSE 3000
