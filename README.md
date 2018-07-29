[![Build Status](https://travis-ci.org/JimYJ/scanproxy.svg?branch=master)](https://travis-ci.org/JimYJ/scanproxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/JimYJ/scanproxy)](https://goreportcard.com/report/github.com/JimYJ/scanproxy)

# scanproxy
scanproxy is auto scan IP &amp; port,and check that is proxy if port is open...(scanproxy是一个自动扫描端口，并且检测是否是代理服务器的程序)

### Command line parameter:命令行参数
<br>
  -a  country codes, see ISO 3166-1 (default "CN")<br>
  -f  scan mode, fast is only scan common port, default is scan all port<br>
  -i  set scan how many IP segment in same times, it will affect memory footprint (default 10)<br>
  -m  maximum concurrency number (default 200)<br>
<br>
### e.g.

```
scanproxy_linux_amd64 -a JP -f -m 1000 -i 20
``` 
