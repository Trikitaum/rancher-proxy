FROM golang
MAINTAINER Cesar gonzalez <cesargonz1984@gmail.com>

RUN apt-get update
RUN apt-get install -y nginx
RUN rm /etc/nginx/sites-enabled/default
ADD . /go/src/cgonzalez/docker-proxy
RUN go install cgonzalez/docker-proxy
ADD run.sh /run.sh
WORKDIR /etc/nginx
# Append "daemon off;" to the beginning of the configuration
RUN echo "daemon off;" >> /etc/nginx/nginx.conf
EXPOSE 80
# Define default command.

CMD /run.sh
