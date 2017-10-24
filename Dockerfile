FROM golang:1.8
MAINTAINER zihua <yuanzihua0@gmail.com>

WORKDIR /go/src/app
COPY game2048 .
RUN chmod +x game2048
EXPOSE 8080
CMD ["./game2048"]