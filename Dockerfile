FROM golang:1.20 as builder

WORKDIR /workspace

# Copy Go Module manifest & dependency files
COPY go.mod go.mod
#COPY go.sum go.sum

# Install deps
RUN go mod download

# Copy source files
COPY main.go main.go
#COPY pkg/ pkg/

# Build the thing
RUN CGO_ENABLED=0 go build -o job main.go

# Use distroless as minimal base image to package the binary
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/job .
USER 65532:65532

ENTRYPOINT ["/job"]
