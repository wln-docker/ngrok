# 使用Docker搭建Ngrok服务器

## 准备工作
* 公网服务器一台并且安装docker
* 域名一枚,做Ngrok服务器域名，如：ngrok.wlniao.com 
* wlniao/ngrok Docker镜像
* 拉取镜像 docker pull wlniao/ngrok

## 启动一个容器生成ngrok客户端,服务器端和CA证书
```linux
docker run --rm -it -e DOMAIN="ngrok.wlniao.com" -v /data/ngrok:/wln wlniao/ngrok /bin/sh /build.sh
```
当看到build ok !的时候,就可以在我们挂载的宿主目录/data/ngrok/bin下看到生成的客户端和服务端

```
bin/ngrokd                  服务端
bin/ngrok                   linux客户端
bin/darwin_amd64/ngrok      osx客户端
bin/windows_amd64/ngrok.exe windows客户端
```

## 启动Ngrok server
直接挂载刚刚的/data/ngrok到容器即可启动服务

## 服务端参数说明
```
 -domain 访问ngrok是所设置的服务地址生成证书时那个
 -httpAddr http协议端口 默认为80
 -httpsAddr https协议端口 默认为443 （可配置https证书）
 -tunnelAddr 通道端口 默认4443
```

## 客户端连接
* 下载我们生成的客户端（官方下载：https://ngrok.com/download）
* 首先创建一个ngrok.cfg配置文件
```
server_addr: "ngrok.wlniao.com:4443"
trust_host_root_certs: false
tunnels:
  web:
    proto:
      http: "3000"
  web2:
    proto:
      http: "8080"
  tcp:
    proto:
      tcp: "3022"
    remote_port: 4444
  tcp2:
    proto:
      tcp: "22"
    remote_port: 4445
```
## 启动客户端命令
```
./ngrok -config ./ngrok.cfg
./ngrok -config ./ngrok.cfg -proto=http -subdomain=test 8080 #启动http转发
./ngrok -config ./ngrok.cfg -proto=tpc 3389  #启动tcp转发 本地3389端口，远程随机暴露大端口
```
## 参数说明
```
trust_host_root_certs #是否信任系统根证书，如果是带自签名证书编译的 ngrok 客户端，这个值应该设置为 false；如果使用 CA 证书，或者用户添加了根证书，这个值应该设置为 true
-proto     #转发协议 不指定默认是 http+https
-subdomain #访问本地时的三级域名 不指定就会随机生成 tcp不支持此参数
8080       #本地服务的端口号
-config    #指定配置位置，默认为$HOME/.ngrok2/ngrok.yml
```