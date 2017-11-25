FROM circleci/golang:1.8

RUN curl -L \
 https://github.com/docker/compose/releases/download/1.11.2/docker-compose-`uname -s`-`uname -m` > ~/docker-compose
RUN chmod +x ~/docker-compose
RUN sudo mv ~/docker-compose /usr/local/bin/docker-compose
