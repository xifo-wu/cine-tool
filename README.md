# CINE TOOL

## 功能

- [ ] 软链接
- [x] AList 重定向服务


## Next TODO

- 尝试优化 302 重定向的加载速度
- 媒体文件夹重命名
- 保存分享的文件支持同步生成软链
- 听说 EMBY 自身的目录监控是不是不好使，可能监听不到变化。想办法优化

## TIPS

1. AList 302 重定向因为跨域的原因不支持 EMBY WEB 浏览器访问，想解决跨域问题，可以搜索 chrome cors proxy 插件解决。

## DEV

##### 生成 CD2 GRPC 代码

```bash
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
core/pb/CloudDrive.proto
```

