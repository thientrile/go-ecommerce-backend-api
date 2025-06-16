FROM golang:alpine AS builder

WORKDIR /build
COPY . .
# # Set environment
ENV ENV=docker

RUN go mod download

RUN go build -o crm.shopdev.com ./cmd/server

FROM scratch

COPY ./configs /configs
COPY --from=builder /build/crm.shopdev.com /

ENTRYPOINT ["/crm.shopdev.com", "configs/docker.yaml"]