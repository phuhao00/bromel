shadowsocks 翻墙



1.https://github.com/shadowsocks/go-shadowsocks2

ubantu

```shell

```



sever

```shell
go-shadowsocks2 -s 'ss://AEAD_CHACHA20_POLY1305:your-password@:8488' -verbose
```



client



```shell
go-shadowsocks2 -c 'ss://AEAD_CHACHA20_POLY1305:your-password@[server_address]:8488' \
    -verbose -socks :1080 -u -udptun :8053=8.8.8.8:53,:8054=8.8.4.4:53 \
                             -tcptun :8053=8.8.8.8:53,:8054=8.8.4.4:53
```

下载：

#### firefox



#### SwitchyOmega

https://github.com/FelisCatus/SwitchyOmega

https://github.com/FelisCatus/SwitchyOmega/releases/download/v2.5.20/proxy_switchyomega-2.5.20-an+fx.xpi



设置proxy

socks5 ;127.0.0.1;1080