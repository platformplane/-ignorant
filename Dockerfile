FROM golang:1-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build github.com/platformplane/scanner/cmd/scanner
RUN CGO_ENABLED=0 GOOS=linux go build github.com/platformplane/scanner/cmd/converter


FROM alpine:3

RUN apk add --no-cache tini curl git
RUN git config --global --add safe.directory '*'

# Trivy CLI
ENV TRIVY_VERSION=0.55.2
ENV TRIVY_CACHE_DIR=/cache/trivy

RUN arch=$(uname -m) && \
    if [ "${arch}" = "x86_64" ]; then \
    arch="64bit"; \
    elif [ "${arch}" = "aarch64" ]; then \
    arch="ARM64"; \
    fi && \
    curl -fsSL https://github.com/aquasecurity/trivy/releases/download/v${TRIVY_VERSION}/trivy_${TRIVY_VERSION}_Linux-${arch}.tar.gz | tar -zxf - -C /usr/local/bin/ trivy

# Gitleaks CLI
ENV GITLEAKS_VERSION=8.19.3
RUN arch=$(uname -m) && \
    if [ "${arch}" = "x86_64" ]; then \
    arch="x64"; \
    elif [ "${arch}" = "aarch64" ]; then \
    arch="arm64"; \
    fi && \
    curl -fsSL https://github.com/gitleaks/gitleaks/releases/download/v${GITLEAKS_VERSION}/gitleaks_${GITLEAKS_VERSION}_linux_${arch}.tar.gz | tar -zxf - -C /usr/local/bin/ gitleaks

COPY --from=build /src/scanner /usr/local/bin/scanner
COPY --from=build /src/converter /usr/local/bin/converter

VOLUME /cache

WORKDIR /src

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/local/bin/scanner"]