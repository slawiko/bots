FROM golang:1.18-alpine3.15 as build

ADD . /build
WORKDIR /build

RUN apk add --no-cache --update git && \
    rm -rf /var/cache/apk/*
    
RUN cd bot && go build -o /dist/main

FROM alpine:3.15

ENV USER=botuser \
    GROUP=botgroup \
    APP_DIR=/bot

COPY --from=build /dist /bot

RUN addgroup -S ${GROUP} && \
    adduser -S ${USER} ${GROUP} && \
    chown ${USER}:${GROUP} /bot

USER ${USER}

ENTRYPOINT [ "/bot/main" ]
