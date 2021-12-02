GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build

build:
	rm -rf target/
	mkdir target/
	cp cmd/logic/logic-example.toml target/logic.toml
	cp cmd/comet/comet-example.toml target/comet.toml
	cp cmd/job/job-example.toml target/job.toml
	$(GOBUILD) -o target/logic cmd/logic/main.go
	$(GOBUILD) -o target/comet cmd/comet/main.go
	$(GOBUILD) -o target/job cmd/job/main.go

clean:
	rm -rf target/

run:
	nohup target/logic -conf=target/logic.toml 2>&1 > target/logic.log &
	nohup target/comet -conf=target/comet.toml 2>&1 > target/comet.log &
	nohup target/job -conf=target/job.toml 2>&1 > target/job.log &
