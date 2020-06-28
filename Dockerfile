FROM golang:1.14-alpine

ENV GOPATH /go
ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GIT_TERMINAL_PROMPT=1
ENV GIT_SSL_NO_VERIFY=true

WORKDIR $GOPATH/src/github.com/aminueza/docker-gluster-exporter

COPY . .

RUN  \
     apk add --update --no-cache git && \
     go build . && cp docker-gluester-exporter /go/bin/docker-gluester-exporter


FROM gluster/gluster-centos

COPY --from=0 /go/bin/docker-gluester-exporter /usr/bin/docker-gluester-exporter

WORKDIR /app

COPY configureGluster.sh .
RUN chmod +x configureGluster.sh

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod 777 /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]