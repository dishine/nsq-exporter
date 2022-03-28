#!/bin/sh

# 打包
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -o nsq_exporter *.go

# 传送
scp nsq_exporter root@152.136.148.67:/root/

# 授权
chmod +x nsq_exporter

nohup ./nsq_exporter \
-nsqd.addr=http://10.18.7.9:4151 \
-collect=stats.topics,stats.channels,stats.clients \
>/data/prometheus_exporter/nsq_exporter/nsq_exporter.log 2>&1 &

