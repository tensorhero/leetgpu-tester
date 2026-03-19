# Stage 1: Build the Go tester binary
FROM golang:1.24-bookworm AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /leetgpu-tester .

# Stage 2: Python runtime with Numba (CUDA simulator)
FROM python:3.11-slim-bookworm
RUN pip install --no-cache-dir numba>=0.59 numpy>=1.24
ENV NUMBA_ENABLE_CUDASIM=1
RUN useradd -m -u 1000 tester
COPY --from=builder /leetgpu-tester /usr/local/bin/leetgpu-tester
USER tester
WORKDIR /workspace
ENTRYPOINT ["leetgpu-tester"]
