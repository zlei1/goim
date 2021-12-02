# goim

## 分层

### logic
提供 HTTP 接口

### comet
websocket 接入层

### job
队列 消费logic产生的消息

## 外部依赖
 * redis
 * etcd v3

## 启动
```
make build
make run

或者

rm -rf target/
mkdir target/
cp cmd/logic/logic-example.toml target/logic.toml
cp cmd/comet/comet-example.toml target/comet.toml
cp cmd/job/job-example.toml target/job.toml
GO111MODULE=on go -o target/logic cmd/logic/main.go
GO111MODULE=on go -o target/comet cmd/comet/main.go
GO111MODULE=on go -o target/job cmd/job/main.go
target/logic -conf=target/logic.toml
target/comet -conf=target/comet.toml
target/job -conf=target/job.toml
```
