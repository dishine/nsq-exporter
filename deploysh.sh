#!/bin/sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -o nsq_exporter *.go

scp nsq_exporter root@152.136.148.67:/root/




# nohup ./nsq_exporter -nsqd.addr=http://10.18.7.9:4151 -collect=stats.topics,stats.channels,stats.clients

nohup /usr/local/prometheus_exporter/node_exporter/node_exporter >/data/prometheus_exporter/node_exporter/node_exporter.log 2>&1 &