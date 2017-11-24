FROM golang:1.9-stretch
ENV CODE=/go/src/jaumecapdevila/go-distributed-system
ENV WORKDIR=/app
ENV BIN=go-distributed-system

WORKDIR $WORKDIR

RUN apt-get update && \
    apt-get install git make

ADD ./src $CODE
RUN cd $CODE && make install folder=$WORKDIR bin=$BIN

ENTRYPOINT "$WORKDIR/$BIN"
