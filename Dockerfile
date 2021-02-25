# Build Go
FROM golang:1.16.0-alpine3.13 as builder
WORKDIR /build-dir
COPY go.mod go.sum ./
RUN go mod download -x
RUN apk --no-cache add build-base
COPY . .
RUN CGO_ENABLED=1 go build -installsuffix cgo -o app .
RUN find . -name "*.go" -type f -delete && find . -type d -empty -delete \
    && rm go.mod go.sum web-2fa/package.json web-2fa/package-lock.json


# Build css
FROM node:14.16-alpine3.12 as buildernode
WORKDIR /build-dir
COPY web-2fa web-2fa/
COPY templates templates/
RUN cd web-2fa && npm i && npm run build_prod


# Output container
FROM alpine:3.13
WORKDIR /root/
COPY --from=builder /build-dir ./
COPY --from=buildernode /build-dir/web-2fa/css/* ./web-2fa/css/
RUN ls -lsh ./web-2fa && echo "done"
CMD ["./app"]