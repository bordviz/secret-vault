FROM golang:1.22.1

RUN mkdir /vault
WORKDIR /vault

COPY . .
RUN chmod a+x docker/*.sh