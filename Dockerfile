FROM golang:1.14-buster as goBuilder
WORKDIR /project
COPY . .
RUN rm -rf .env
RUN make

FROM debian:buster
WORKDIR /project
COPY --from=goBuilder /project/var/canbot /usr/local/bin/
RUN apt-get update
RUN apt-get install -y ca-certificates
CMD ["canbot"]
