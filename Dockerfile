############################
# STEP 1 build executable binary
############################
FROM golang:1.21-alpine as builder

ARG DEPLOY_USER
ARG DEPLOY_TOKEN

# For ssh
ARG DEPLOY_KEY

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates curl openssh-client

RUN mkdir /root/.ssh/

# Copy SSH Key
RUN echo "${DEPLOY_KEY}" > /root/.ssh/id_rsa 

# make sure your domain is accepted
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN chmod 400 ${HOME}/.ssh/id_rsa

# Fetch dependencies.
WORKDIR /code
RUN if [ -n "$DEPLOY_TOKEN" ]; then \
    git config --global url."https://${DEPLOY_USER}:${DEPLOY_TOKEN}@github.com/".insteadOf "https://github.com/"; \
    else \
    git config --global url."git@github.com:".insteadOf "https://github.com/"; \
    fi

# Fetch dependencies.
WORKDIR /code
COPY go.mod go.sum /code/
RUN go mod download

# Copy source code
ADD . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPRIVATE=github.com/nawaltni/api/gen/go/nawalt go build  -o /go/bin/tracker

############################
# STEP 2 build a small image
############################
FROM alpine
# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
RUN true
COPY --from=builder /go/bin/tracker /go/bin/tracker

# Create appuser
RUN adduser -D -g '' appuser
WORKDIR /code
COPY ./db /code/db
COPY config.toml ./
RUN chown appuser /code


# Use an unprivileged user.
USER appuser
# Run the parser binary.
ENTRYPOINT ["/go/bin/tracker"]