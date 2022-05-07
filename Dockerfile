FROM golang:1.18-alpine3.15 as build

ENV SRC_DIR=/build \
    DIST_DIR=/dist

ADD . ${SRC_DIR}
WORKDIR ${SRC_DIR}

RUN apk add --no-cache --update git && \
    rm -rf /var/cache/apk/*
    
RUN cd bot && go build -o ${DIST_DIR}/main

FROM alpine:3.15

ENV USER=botuser \
    GROUP=botgroup \
    APP_DIR=/bot

COPY --from=build ${DIST_DIR} ${APP_DIR}

RUN addgroup -S ${GROUP} && \
    adduser -S ${USER} ${GROUP} && \
    chown ${USER}:${GROUP} ${APP_DIR}

USER ${USER}

ENTRYPOINT [ "/bot/main" ]
