# Zookeeper Monitor

Zookeeper Monitor 针对 Hyperf 框架服务治理组件，Zookeeper驱动组件服务注册信息监控。

### 依赖类库

```bash
$ go get -v github.com/larspensjo/config
$ go get -v github.com/samuel/go-zookeeper/zk
$ go get -v github.com/json-iterator/go
```

### 交叉编译

1. Mac 下编译 Linux 和 Windows 64位可执行程序

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/zookeeper-monitor-linux-x86_64 src/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/zookeeper-monitor-windows-x86_64.exe src/main.go
```

2. Linux 下编译 Mac 和 Windows 64位可执行程序

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/zookeeper-monitor-darwin-amd64 src/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/zookeeper-monitor-windows-x86_64.exe src/main.go
```

3. Windows 下编译 Mac 和 Linux 64位可执行程序

```bash
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o bin/zookeeper-monitor-darwin-amd64 src/main.go

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o bin/zookeeper-monitor-linux-x86_64 src/main.go
```