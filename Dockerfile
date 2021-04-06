FROM circleci/golang:1.15
USER root
WORKDIR /go/src/github.com/ntindall/sql-gen-doc

RUN apt-get update \
   && apt-get install -y --force-yes --no-install-recommends\
   mysql-client \
   libmysqlclient-dev

COPY Makefile ./
ADD scripts/logs.sh ./scripts/
COPY go.mod go.sum ./
RUN make setup
COPY ./ ./
RUN make build
