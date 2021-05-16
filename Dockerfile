FROM golang:1.15.6-alpine AS builder

RUN go version
RUN apk add git

COPY ./ /github.com/psihachina/windfarms-backend
WORKDIR /github.com/psihachina/windfarms-backend

RUN go mod download && go get -u ./...
RUN CGO_ENABLE=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

FROM ubuntu:latest

WORKDIR /root/

RUN apt-get update && apt-get install -y \
                ca-certificates \
                curl \
                wget \
                build-essential \
                bzip2 \
                tar \
                amqp-tools \
                openssh-client \
                gfortran \
                --no-install-recommends && rm -r /var/lib/apt/lists/* && \
                wget ftp://ftp.cpc.ncep.noaa.gov/wd51we/wgrib2/wgrib2.tgz.v2.0.4 -O /tmp/wgrib2.tgz && \
                mkdir -p /usr/local/grib2/ && \
                cd /tmp/ && \
                tar -xf /tmp/wgrib2.tgz && \
                rm -r /tmp/wgrib2.tgz && \
                mv /tmp/grib2/ /usr/local/grib2/ &&\
                cd /usr/local/grib2/grib2 && \
                make && \
                ln -s /usr/local/grib2/grib2/wgrib2/wgrib2 /usr/local/bin/wgrib2 && \
                apt-get -y autoremove build-essential

COPY --from=0 /github.com/psihachina/windfarms-backend/.bin/app .
COPY --from=0 /github.com/psihachina/windfarms-backend/scripts .
COPY --from=0 /github.com/psihachina/windfarms-backend/configs/ ./configs/

EXPOSE 8000

CMD ["./app"]
