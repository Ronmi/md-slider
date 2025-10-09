FROM node AS node_builder
WORKDIR /src
COPY . .
WORKDIR /src/js
ENV CI=true
RUN corepack enable pnpm && pnpm install && pnpm run tsc

FROM golang AS builder

WORKDIR /src
COPY --from=node_builder /src /src
RUN mkdir -p data ; go build -o md-slider

FROM ronmi/mingo:latest
LABEL maintainer="Ronmi Ren <ronmi.ren@gmail.com>"
COPY --from=builder /src/md-slider /md-slider
COPY --from=builder /src/data /data
WORKDIR /data
VOLUME ["/data"]
EXPOSE 8000
CMD ["/md-slider"]
