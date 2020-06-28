FROM golang:1.14-alpine

ENV GOPATH /go
ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GIT_TERMINAL_PROMPT=1
ENV GIT_SSL_NO_VERIFY=true

WORKDIR $GOPATH/src/github.com/aminueza/docker-gluster-prometheus

COPY . .

RUN  \
     apk add --update --no-cache git && \
     go build . && cp docker-gluster-prometheus /go/bin/docker-gluster-prometheus


FROM gluster/gluster-centos

COPY --from=0 /go/bin/docker-gluster-prometheus /usr/bin/docker-gluster-prometheus

WORKDIR /app

COPY configureGluster.sh .
RUN chmod +x configureGluster.sh

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod 777 /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]