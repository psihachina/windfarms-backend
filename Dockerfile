FROM golang:1.15.6-alpine AS builder

RUN go version
RUN apk add git

COPY ./ /github.com/psihachina/windfarms-backend
WORKDIR /github.com/psihachina/windfarms-backend

RUN go mod download && go get -u ./...
RUN CGO_ENABLE=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates libc-dev gfortran gcc make
RUN wget ftp://ftp.cpc.ncep.noaa.gov/wd51we/wgrib2/wgrib2.tgz
RUN tar -xzvf wgrib2.tgz 
RUN cd grib2 && export FC=gfortran && export CC=gcc && make && cd wgrib2 && cp wgrib2 /bin/

COPY --from=0 /github.com/psihachina/windfarms-backend/.bin/app .
COPY --from=0 /github.com/psihachina/windfarms-backend/scripts ./scripts/
COPY --from=0 /github.com/psihachina/windfarms-backend/configs/ ./configs/

EXPOSE 8000

CMD ["./app"]
