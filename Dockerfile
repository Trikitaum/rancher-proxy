FROM golang
MAINTAINER Cesar gonzalez <cesargonz1984@gmail.com>

RUN apt-get update
RUN apt-get install -y nginx python python-dev python-pip

RUN rm /etc/nginx/sites-enabled/default
ADD . /go/src/cgonzalez/docker-proxy
RUN go install cgonzalez/docker-proxy
ENTRYPOINT /go/bin/docker-proxy
