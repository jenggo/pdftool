FROM oven/bun:alpine AS frontend

WORKDIR /app

COPY ui/package.json ui/bun.lock ./
RUN bun install

COPY ui/ ./
RUN bun run check
RUN bun run build

FROM golang:alpine AS backend

RUN apk add musl-dev
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s

WORKDIR /src

COPY go.* ./
RUN go mod download -x

COPY . .
RUN golangci-lint run --timeout=2m
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o pdftool -trimpath


FROM chainguard/static

COPY --from=backend /src/pdftool .
COPY --from=frontend /app/build ui/

EXPOSE 2804

ENTRYPOINT ["./pdftool"]

