FROM golang:1.7.5-alpine

RUN apk add --no-cache gcc libc-dev bash git mysql-client curl openssh-client
