# 构建基础镜像
```dockerfile
FROM python:3.15.0a8-slim

ARG DEBIAN_MIRROR=mirrors.tuna.tsinghua.edu.cn
ARG DEBIAN_MIRROR_PATH=/etc/apt/sources.list.d/debian.sources

ENV PIP_INDEX_URL=https://mirrors.cloud.tencent.com/pypi/simple/
ENV PIP_TRUSTED_HOST=mirrors.cloud.tencent.com
ENV PIP_DEFAULT_TIMEOUT=100
ENV GOPROXY=https://goproxy.cn,direct

RUN sed -i "s/deb.debian.org/${DEBIAN_MIRROR}/g" ${DEBIAN_MIRROR_PATH}

RUN apt update \
        && apt-get install -y --no-install-recommends \
        xz-utils \
        pkg-config \ 
	gcc \
	libseccomp-dev \
	build-essential \
        curl \
        strace \
        nano \
        git

RUN pip3 install --no-cache-dir httpx==0.27.2 requests==2.32.3 jinja2==3.1.6 PySocks httpx[socks]
```
在项目根目录执行 构建命令
```shell
# docker build -f ./docker/dify-sandbox-base.dockerfile -t dify-sandbox:base .                                                          
```
查看构建的镜像
```shell
# docker images|grep base
WARNING: This output is designed for human readability. For machine-readable output, please use --format.
dify-sandbox:base                         31b0f0fc0ffe        636MB          172MB
```

为啥使用基础镜像？安装依赖的时间太久了。基于基础镜像构建下次就不用下载这些依赖了

# 构建运行的容器镜像
```dockerfile
FROM dify-sandbox:base

ARG REPO_NAME=dify-sandbox
ARG RESOURCE_PATH=resource

WORKDIR /app

RUN git clone https://github.com/langgenius/${REPO_NAME}.git

WORKDIR ${REPO_NAME}

COPY ./docker ${RESOURCE_PATH}

RUN rm -rf /usr/local/go && tar -C /usr/local -xzf ./${RESOURCE_PATH}/install-pkg/go1.25.6.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

RUN tar -xvf ./${RESOURCE_PATH}/install-pkg/node-v24.14.1-linux-x64.tar.xz -C /opt \
	&& mkdir -p /usr/local/bin \
	&& ln -s /opt/node-v24.14.1-linux-x64/bin/node /usr/local/bin/node

RUN chmod +x ./build/build_amd64.sh \
	&& ./build/build_amd64.sh 

RUN rm -f ./conf/config.yaml && cp ./resource/config/config.yaml  ./conf


# CMD ["tail", "-f", "/dev/null"]
ENTRYPOINT ["./main"]
```
执行构建
```shell
# docker build -f ./docker/dify-sandbox.dockerfile -t dify-sandbox:source .
```

启动服务
```shell
# docker run -d -p 8194:8194 --name dify-sandbox dify-sandbox:source
```

查看服务日志
```shell
# docker logs -f dify-sandbox
```

# 运行测试脚本
测试脚本内容
```bash
#!/bin/bash
set -eu

SCRIPT_DIR="$(dirname "$(realpath "$0")")"

source ${SCRIPT_DIR}/test_env.sh

# 健康检查接口
echo "健康检查接口:/health"
curl ${REQUEST_HOST}/health

# Python3 示例
echo ""
echo "Python3 运行代码:/v1/sandbox/run"
curl -X POST ${REQUEST_HOST}/v1/sandbox/run \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: ${X_API_KEY}" \
  -d '{
    "language": "python3",
    "code": "import json\nperson = {\"name\": \"John\", \"age\": 30, \"city\": \"New York\"}\njson_str = json.dumps(person)\nprint(json_str)",
    "preload": "",
    "enable_network": false
  }'

echo ""
echo "Node.js 运行代码:/v1/sandbox/run"
# Node.js 示例
curl -X POST ${REQUEST_HOST}/v1/sandbox/run \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: ${X_API_KEY}" \
  -d '{
    "language": "nodejs",
    "code": "const person = {name: \"John\", age: 30, city: \"New York\"};\nconst jsonString = JSON.stringify(person);\nconsole.log(jsonString);",
    "preload": "",
    "enable_network": false
  }'

echo ""

```
```shell
$ bash test-shell/test_dify_sandbox_run.sh 
健康检查接口:/health
"ok"
Python3 运行代码:/v1/sandbox/run
{"code":0,"message":"success","data":{"error":"error: operation not permitted\n","stdout":"{\"name\": \"John\", \"age\": 30, \"city\": \"New York\"}\n"}}
Node.js 运行代码:/v1/sandbox/run
{"code":0,"message":"success","data":{"error":"","stdout":"{\"name\":\"John\",\"age\":30,\"city\":\"New York\"}\n"}}

```
操作系统信息
```bash
$ cat /etc/os-release 
PRETTY_NAME="Ubuntu 24.04.4 LTS"
NAME="Ubuntu"
VERSION_ID="24.04"
VERSION="24.04.4 LTS (Noble Numbat)"
VERSION_CODENAME=noble
ID=ubuntu
ID_LIKE=debian
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=noble
LOGO=ubuntu-logo

$ uname -a
Linux scd-X99 6.17.0-22-generic #22~24.04.1-Ubuntu SMP PREEMPT_DYNAMIC Thu Mar 26 15:25:54 UTC 2 x86_64 x86_64 x86_64 GNU/Linux
```
这里的返回结果为啥 `"error":"error: operation not permitted\n"`?
参考我之前文章的分析 [https://blog.csdn.net/modelmd/article/details/160346298](https://blog.csdn.net/modelmd/article/details/160346298), 不同的操作系统可能原因不一样

我的本地是由于`SYSCALL=epoll_pwait` 系统调用被限制了
只需要在参数调用的时候改成 `"enable_network": true` 就可以允许`SYS_EPOLL_PWAIT = 281`  epoll_pwait 系统调用了
也可以直接修改源码，官方的源码不是很安全。修改参数 `"enable_network": true` 之后再次运行

```bash
$ bash test-shell/test_dify_sandbox_run.sh 
健康检查接口:/health
"ok"
Python3 运行代码:/v1/sandbox/run
{"code":0,"message":"success","data":{"error":"","stdout":"{\"name\": \"John\", \"age\": 30, \"city\": \"New York\"}\n"}}
Node.js 运行代码:/v1/sandbox/run
{"code":0,"message":"success","data":{"error":"","stdout":"{\"name\":\"John\",\"age\":30,\"city\":\"New York\"}\n"}}
```
