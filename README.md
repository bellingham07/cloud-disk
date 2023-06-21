# Cloud Disk    
> 轻量级云盘系统，基于go-zero ,xorm实现

使用到的命令
```text
# 创建API服务
goctl api new core
# 启动服务
go run core.go -f etc/core-api.yaml
# 使用API文件生成代码
goctl api go -api core.api -dir . -style go_zero
# 格式化代码
goctl api format --dir core.api 
# 阿里云oss
go get github.com/aliyun/aliyun-oss-go-sdk/oss
```