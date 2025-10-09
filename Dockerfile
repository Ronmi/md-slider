FROM golang AS builder

WORKDIR /src
COPY . .
RUN mkdir -p data ; go build -o md-slider

FROM ronmi/mingo:latest
LABEL maintainer="Ronmi Ren <ronmi.ren@gmail.com>"
COPY --from=builder /src/md-slider /md-slider
COPY --from=builder /src/data /data
WORKDIR /data
VOLUME ["/data"]
EXPOSE 8000
CMD ["/md-slider"]
