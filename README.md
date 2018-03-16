## deploy
Golang实现的代码部署系统

## 实现流程
打包本地文件 -> 上传远程服务器 -> 解压压缩包 -> 执行命令

## 使用方法
- 进入conf目录下将conf.toml.example重命名为conf.toml并完成配置信息
- go run deploy.go

## TODO
- [ ] 上传远程服务器过程显示进度
- [ ] 并发切片上传压缩包