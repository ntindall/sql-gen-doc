FROM cimg/go:1.20.6
USER root
WORKDIR /go/src/github.com/ntindall/sql-gen-doc

RUN apt-get update \
   && apt-get install -y --force-yes --no-install-recommends\
   default-mysql-client \
   default-libmysqlclient-dev \

COPY Makefile ./
ADD scripts/logs.sh ./scripts/
COPY go.mod go.sum ./
COPY ./ ./
RUN make setup
RUN make build
