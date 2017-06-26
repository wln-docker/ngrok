# 使用Docker搭建Ngrok服务器

## 准备工作
* 公网服务器一台并且安装docker
* 域名一枚,做Ngrok服务器域名，如：wlniao.cn
* 生成自己的CA根证书及Ngrok域名证书，并覆盖cert目录中的ca.cer、server.crt及server.key文件
* 拉取镜像 docker pull wlniao/ngrok

## 首先启动一个容器生成ngrok客户端
```linux
docker run --rm -it -e DOMAIN="wlniao.cn" -v /docker/ngrok:/wln wlniao/ngrok /bin/sh /start.sh
```
当看到build ok !的时候,就可以在我们挂载的宿主目录/docker/ngrok下看到生成的客户端和服务端

```
ngrok                     linux客户端
ngrokd                    linux服务端
windows_amd64/ngrok.exe   windows客户端
windows_amd64/ngrokd.exe  windows服务端
```

## 启动Ngrok server
直接挂载刚刚的/docker/ngrok到容器/wln目录即可启动服务
```linux
docker run -d -p 4443:4443 -e DOMAIN="wlniao.cn" -v /docker/ngrok:/wln wlniao/ngrok
```

## 服务端参数说明
```
 -domain 访问ngrok是所设置的服务地址生成证书时那个
 -httpAddr http协议端口 默认为80
 -httpsAddr https协议端口 默认为443 （可配置https证书）
 -tunnelAddr 通道端口 默认4443（管理用）
```

## 客户端连接
* 下载我们生成的客户端（官方下载：https://ngrok.com/download）
* 首先创建一个ngrok.cfg配置文件
* trust_host_root_certs #是否信任系统根证书，如果是带自签名证书编译的 ngrok 客户端，这个值应该设置为 false；如果使用 CA 证书，或者用户添加了根证书，这个值应该设置为 true
* proto下协议为http时，lable为子域名
```
server_addr: "wlniao.cn:4443"
trust_host_root_certs: false
tunnels:
  test:
    proto:
      http: 80
```

## 启动客户端命令
```
ngrok -config ngrok.cfg start test            #启动指定转发配置
ngrok -config ngrok.cfg start-all             #启动全部转发配置
ngrok -config ngrok.cfg 4040                  #启动web转发，随机子域名
ngrok -config=ngrok.cfg -subdomain=test 4040  #启动http转发，子域名：test
ngrok -config ngrok.cfg -proto=tcp 127.0.0.3:3389   #启动tcp转发 本地3389端口，远程随机暴露大端口
```
## 启动命令参数说明
```
-config    #指定配置位置
-proto     #转发协议 不指定默认是 http+https
-subdomain #访问本地时的三级域名 不指定就会随机生成 tcp不支持此参数
8080       #本地服务的端口号
```