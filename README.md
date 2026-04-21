# pdf-service

A lightweight HTTP service that applies text watermarks to PDFs using a pdflatex pipeline. Used by Terminal and other StairSupplies applications via `PandocService.php`.

## API

### `GET /health`

Liveness check. No authentication required.

**Response `200 OK`**
```json
{ "status": "ok" }
```

---

### `POST /watermark`

Accepts a PDF as the request body and returns a watermarked PDF.

**Request**

| Header | Required | Default | Description |
|---|---|---|---|
| `Authorization` | When `PDF_SERVICE_API_KEY` is set | — | `Bearer <key>` |
| `X-Watermark-Text` | Yes | — | Text to stamp on each page |
| `X-Watermark-Color` | No | `red` | xcolor color name (e.g. `gray`, `blue`) |
| `X-Watermark-Opacity` | No | `0.5` | Float, 0.0–1.0 |
| `X-Watermark-Size` | No | `60` | Font size in pt |
| `X-Watermark-Position` | No | `top-centre` | `top-left`, `top-centre`, `top-right`, `centre`, `bottom-left`, `bottom-centre`, `bottom-right` |
| `X-Watermark-Angle` | No | `0` | Rotation in degrees |
| `X-Watermark-Bold` | No | `false` | Set to `true` for bold weight |

Body: raw PDF bytes (`Content-Type: application/pdf`). Maximum 120 MB.

**Response**

- `200 OK` — watermarked PDF bytes (`Content-Type: application/pdf`)
- `400 Bad Request` — missing or invalid header
- `401 Unauthorized` — missing or wrong bearer token (only when `PDF_SERVICE_API_KEY` is set)
- `413 Request Entity Too Large` — body exceeds 120 MB
- `500 Internal Server Error` — pdflatex failed; body contains pdflatex output for debugging

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PDF_SERVICE_API_KEY` | _(unset)_ | Bearer token for auth. When unset, all requests pass through unauthenticated. |
| `PORT` | `8080` | Port the HTTP server listens on. |
| `PDFLATEX_TIMEOUT` | `55` | Per-request pdflatex timeout in seconds (covers both passes). |
| `WRITE_TIMEOUT` | `30` | HTTP server write timeout in seconds. Set higher (e.g. `120`) when processing large PDFs. |
| `APP_ENV` | _(unset)_ | Set to `production` for structured JSON logging; omit for human-readable console output. |

## Running Locally

Requires Go 1.23+ and `pdflatex` (TeX Live).

```bash
make run
```

The server starts on port 8080 by default. Test with:

```bash
curl http://localhost:8080/health

curl -X POST http://localhost:8080/watermark -H "X-Watermark-Text: DRAFT" --data-binary @input.pdf -o output.pdf
```

## Docker

```bash
make docker-build        # build image tagged pdf-service:local
make docker-smoke        # build, start, verify /health, stop
```

Or manually:

```bash
docker run -d --rm -p 8080:8080 -e PDF_SERVICE_API_KEY=secret pdf-service:local
```

## Tests

```bash
make test
```

The `TestWatermark_ContentType` integration test requires `pdflatex` to be installed and is automatically skipped when it is not.

## CI / CD

The `docker-publish.yml` workflow triggers on `v*` tag pushes. It runs `go test ./...` first and only builds and pushes the image to `ghcr.io/stairsupplies/pdf-service` if tests pass.

```bash
git tag v1.2.0
git push origin v1.2.0
```

## terminal-docker Integration

In `docker-compose.yml` the service runs at host port `3006`. Configure via `.env`:

```dotenv
PANDOC_SERVICE_URL=http://pdf-service:8080
PANDOC_SERVICE_API_KEY=your-secret-key
PANDOC_SERVICE_TIMEOUT=60
```

The matching `PDF_SERVICE_API_KEY` must be set on the `pdf-service` container for auth to be enforced.
