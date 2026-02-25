FROM golang:1.26

RUN go install github.com/go-delve/delve/cmd/dlv@latest

ARG UID=1000
ARG GID=1000

RUN groupadd -g $GID user && \
  useradd -l -m -u $UID -g $GID -s /usr/bin/bash user

RUN chown -R user:user /go
USER user

WORKDIR /app

CMD [ "sleep", "infinity" ]
