DOCKER-PROXY-RANCHER
----
This is a nginx proxy that use the rancher.io api for service discovery

## BUILD IT

docker build -t docker-proxy .

## RUN IT

docker run -d -p 80:80 -e IP=[rancherServerAPI]-e USER=[rancherUserAPI] -e PASSWORD=[rancherPasswordAPI] docker-proxy
