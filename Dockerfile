# Stage 1 — compile the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o pdf-service .

# Stage 2 — slim runtime with the TeX Live packages required for pdfpages and tikz.
# NOTE: texlive-latex-extra pulls in pgf/tikz (required for overlay watermarking) and
#       pdfpages. lmodern provides Latin Modern fonts — the production-quality successor
#       to Computer Modern with full scalability at arbitrary sizes. Expect a final image
#       size of several hundred MB — this is inherent to any TeX Live installation.
FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        texlive-latex-base \
        texlive-latex-extra \
        texlive-fonts-recommended \
        lmodern \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/pdf-service /usr/local/bin/pdf-service

EXPOSE 8080

ENV APP_ENV=production

USER nobody

CMD ["/usr/local/bin/pdf-service"]
