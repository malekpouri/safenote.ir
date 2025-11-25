# Build WASM Stage
FROM golang:1.24-alpine AS wasm-builder
WORKDIR /wasm
COPY wasm/main.go .
RUN GOOS=js GOARCH=wasm go build -o safenote.wasm main.go
RUN cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .

# Build Frontend Stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .
# Copy WASM artifacts to static folder so they are served
COPY --from=wasm-builder /wasm/safenote.wasm ./static/safenote.wasm
COPY --from=wasm-builder /wasm/wasm_exec.js ./static/wasm_exec.js

RUN npm run build

# Run Stage
FROM node:20-alpine

WORKDIR /app

COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/node_modules ./node_modules

EXPOSE 3000

CMD ["node", "build"]
