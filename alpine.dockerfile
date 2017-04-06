FROM alpine
MAINTAINER Ronmi Ren <ronmi.ren@gmail.com>

ADD ./md-slider /usr/local/bin/
CMD ["md-slider"]
