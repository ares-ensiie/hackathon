FROM golang:1.6.4
MAINTAINER Ares <ares@ares-ensiie.eu>

ADD ./hackathon-go-client /go/src/git.ares-ensiie.eu/hackathon/hackathon-go-client
ADD ./hackathon-server /go/src/git.ares-ensiie.eu/hackathon/hackathon-server

ADD ./hackathon-docker/launch_client.sh /srv/launch_client.sh
ADD ./hackathon-docker/launch_server.sh /srv/launch_server.sh
ADD ./hackathon-docker/launch.sh /srv/launch.sh

RUN  ls -la /go/src/git.ares-ensiie.eu/hackathon/ && \
  cd /go/src/git.ares-ensiie.eu/hackathon/hackathon-go-client && \
  go get && \
  go build && \
  cd /go/src/git.ares-ensiie.eu/hackathon/hackathon-server && \
  go get && \
  go build



EXPOSE 1337
EXPOSE 1338

WORKDIR /srv

CMD ["/srv/launch.sh"]
