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

LABEL name="{{NAME}} Image" \
      vendor="ESSENTIAL KAOS" \
      maintainer="Anton Novojilov" \
      license="Apache-2.0" \
      version="{{DOCKER_CHANGE_DATE}}"

COPY --from=builder /go/src/github.com/essentialkaos/{{SHORT_NAME}}/{{SHORT_NAME}} \
                    /usr/bin/

# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["{{SHORT_NAME}}"]

################################################################################
