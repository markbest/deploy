## deploy
代码部署工具

## process
打包文件 -> 上传服务器 -> 解压 -> 执行后续命令

## usage
- 进入conf目录下将conf.toml.example重命名为conf.toml并完成配置信息
- go run deploy.go

## TODO
- [x] 上传服务器过程显示实时进度
- [x] 大文件分割并发上传