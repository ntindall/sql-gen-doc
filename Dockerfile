FROM circleci/golang:1.8
USER root
WORKDIR /go/src/github.com/ntindall/sql-gen-doc

COPY Makefile ./
ADD scripts scripts
COPY Gopkg.lock Gopkg.toml ./
RUN make setup
COPY cmd/ format/ ./
RUN make build
