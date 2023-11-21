# goim

## 分层

### logic
业务逻辑层，可多节点部署

### comet
用户连接层，可多节点部署

### job
任务推送层，通过 redis 订阅发布功能实现，暂时不支持多节点部署

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
GO111MODULE=on go build -o target/logic cmd/logic/main.go
GO111MODULE=on go build -o target/comet cmd/comet/main.go
GO111MODULE=on go build -o target/job cmd/job/main.go
target/logic -conf=target/logic.toml
target/comet -conf=target/comet.toml
target/job -conf=target/job.toml
```
