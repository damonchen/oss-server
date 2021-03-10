# oss-server

阿里云oss和腾讯云oss的反向代理

## design

![image-20210310130540188](docs/assets/README/browser.png)



![image-20210310130601585](docs/assets/README/logic-view.png)



## usage

```bash

git clone https://github.com/damonchen/oss-server
cd oss-server/cmd/osv
go build
osv web --cfg config/config.yaml
```

将`config.yaml` 复制到你的服务器，按照配置文件中的说明进行配置

```bash
curl http://localhost:8092/oss/${proxy_name}?path=${xxxx}
```

上面的`${tagtag}`替换为你实际配置中的`name`，` ${xxxx}`替换为你实际在oss上的路径地址。（注意，aliyun上的bucket名称不要放进去）




## current support

- aliyun oss, using provider name: `aliyun`
- tencent oss, using provider name: `tencent`




## TODO

- [ ] huawei

```bash
http://localhost:8092/oss/tagtag/
```
