FROM node:16 as web-builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
COPY ./frontend .
RUN DISABLE_ESLINT_PLUGIN='true' npm run build

FROM golang AS go-builder
WORKDIR /app
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web-builder /app/dist ./frontend/dist
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w -extldflags '-static'" -o pandora-fuclaude-plus-helper ./cmd/server/main.go


FROM alpine AS runner
WORKDIR /app
COPY --from=go-builder /app/pandora-fuclaude-plus-helper ./pandora-fuclaude-plus-helper
RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true \
RUN mkdir -p /data
RUN ["chmod", "+x", "/app/pandora-fuclaude-plus-helper"]
ENTRYPOINT ["/app/pandora-fuclaude-plus-helper"]

