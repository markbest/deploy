## deploy
Golang实现的代码部署系统

## 实现流程
打包文件 -> 上传服务器 -> 解压 -> 执行后续命令

## 使用方法
- 进入conf目录下将conf.toml.example重命名为conf.toml并完成配置信息
- go run deploy.go

## TODO
- [ ] 上传服务器过程显示实时进度
- [x] 大文件分割上传