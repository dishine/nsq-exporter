#!/bin/sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -o nsq_exporter *.go

scp nsq_exporter root@152.136.148.67:/root/

chmod +x nsq_exporter

# nohup ./nsq_exporter -nsqd.addr=http://10.18.7.9:4151 -collect=stats.topics,stats.channels,stats.clients

nohup /usr/local/prometheus_exporter/beanstalkd_exporter/beanstalkd_exporter-1.0.5.linux-amd64 \
-beanstalkd.address=10.18.7.37:11300 \
-web.listen-address=:9118 \
>/data/prometheus_exporter/beanstalkd_exporter/beanstalkd_exporter.log 2>&1 &



./beanstalkd_exporter-1.0.5.linux-amd64 -beanstalkd.address=10.18.7.37:11300 -web.listen-address=:9118