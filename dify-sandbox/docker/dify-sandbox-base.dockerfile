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
