FROM circleci/golang:1.8

COPY Makefile scripts/ ./
COPY Gopkg.lock Gopkg.toml ./
RUN make setup
COPY cmd/ format/ ./
RUN make build
