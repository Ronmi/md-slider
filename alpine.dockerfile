FROM alpine
MAINTAINER Ronmi Ren <ronmi.ren@gmail.com>

ADD ./md-slider /usr/local/bin/
WORKDIR /data
VOLUME ["/data"]
USER nobody
CMD ["md-slider"]
