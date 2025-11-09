# Build stage
FROM golang:tip-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o buffman .

# Final stage
FROM alpine:latest

# Arguments for user and group IDs
ARG USERNAME=buffman
ARG USER_UID=1000
ARG USER_GID=1000

# Install dependencies and tools
RUN apk add --no-cache \
    ca-certificates \
    wget \
    unzip \
    libc6-compat \
    bash \
    shadow \
    python3 \
    py3-pip \
    build-base \
    python3-dev \
    libffi-dev \
    openssl-dev

# Install Python packages RIGHT HERE - before user creation
RUN pip install protobuf grpcio-tools --break-system-packages

# Create group and user with matching UID/GID (keep your original commands)
RUN groupadd -g $USER_GID $USERNAME && \
    useradd -u $USER_UID -g $USER_GID -m $USERNAME

# Download and install flatc binary
RUN wget -O /tmp/flatc.zip https://github.com/google/flatbuffers/releases/download/v25.2.10/Linux.flatc.binary.g++-13.zip && \
    unzip /tmp/flatc.zip -d /tmp && \
    mv /tmp/flatc /usr/local/bin/flatc && \
    chmod +x /usr/local/bin/flatc && \
    rm -rf /tmp/flatc.zip

RUN wget -O /tmp/nano.tar.gz https://jpa.kapsi.fi/nanopb/download/nanopb-0.4.9.1.tar.gz && \ 
    tar -xzf /tmp/nano.tar.gz -C /tmp && \
    ln -s /tmp/nanopb/generator/nanopb_generator.py /usr/local/bin/nanopb

# Copy the binary from builder stage
COPY --from=builder /app/buffman /usr/local/bin/buffman

# Set working directory and permissions
WORKDIR /buffman
RUN chown -R $USER_UID:$USER_GID /buffman

# Switch to non-root user
USER $USERNAME

# Set the entrypoint
ENTRYPOINT ["buffman"]
