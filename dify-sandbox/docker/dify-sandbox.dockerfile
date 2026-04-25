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

