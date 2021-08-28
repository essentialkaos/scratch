## BUILDER #####################################################################

FROM golang:alpine as builder

WORKDIR /go/src/github.com/essentialkaos/{{SHORT_NAME}}

COPY . .

ENV GO111MODULE=auto

RUN apk add --no-cache git=~2.32 make=4.3-r0 upx=3.96-r1 && \
    make deps && \
    make all && \
    upx {{SHORT_NAME}}

## FINAL IMAGE #################################################################

FROM essentialkaos/alpine:3.13

LABEL org.opencontainers.image.title="{{NAME}}" \
      org.opencontainers.image.description="{{DESC}}" \
      org.opencontainers.image.vendor="ESSENTIAL KAOS" \
      org.opencontainers.image.authors="Anton Novojilov" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.url="https://kaos.sh/{{SHORT_NAME}}" \
      org.opencontainers.image.source="https://github.com/essentialkaos/{{SHORT_NAME}}"

COPY --from=builder /go/src/github.com/essentialkaos/{{SHORT_NAME}}/{{SHORT_NAME}} \
                    /usr/bin/

# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["{{SHORT_NAME}}"]

################################################################################
