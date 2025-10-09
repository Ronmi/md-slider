FROM ronmi/mingo:latest
LABEL maintainer="Ronmi Ren <ronmi.ren@gmail.com>"
COPY md-slider /md-slider
COPY data /data

VOLUME ["/data"]
EXPOSE 8000
WORKDIR /data
CMD ["/md-slider"]
