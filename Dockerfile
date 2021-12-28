FROM golang:1.17.5

MAINTAINER karlmutch@gmail.com

ENV LANG C.UTF-8

RUN apt-get -y update

RUN apt-get -y install git software-properties-common wget openssl ssh curl jq apt-utils unzip python3-pip && \
    apt-get clean && \
    apt-get autoremove && \
    pip install awscli --upgrade

# Protobuf version
ENV PROTOBUF_VERSION="3.19.1"
ENV PROTOBUF_ZIP=protoc-${PROTOBUF_VERSION}-linux-x86_64.zip
ENV PROTOBUF_URL=https://github.com/google/protobuf/releases/download/v${PROTOBUF_VERSION}/${PROTOBUF_ZIP}

ARG USER
ENV USER ${USER}
ARG USER_ID
ENV USER_ID ${USER_ID}
ARG USER_GROUP_ID
ENV USER_GROUP_ID ${USER_GROUP_ID}

RUN groupadd -f -g ${USER_GROUP_ID} ${USER}
RUN useradd -g ${USER_GROUP_ID} -u ${USER_ID} -ms /bin/bash ${USER}

RUN wget ${PROTOBUF_URL} && \
    unzip ${PROTOBUF_ZIP} -d /usr && \
    chmod +x /usr/bin/protoc && \
    find /usr/include/google -type d -print0 | xargs -0 chmod ugo+rx && \
    chmod -R +r /usr/include/google

USER ${USER}
WORKDIR /home/${USER}

VOLUME /project
WORKDIR /project/src/github.com/karlmutch/platform-services

ENV GOPATH=/project
ENV PATH="${GOPATH}/bin:/home/${USER}/bin:${PATH}:/usr/bin"

RUN \
    mkdir -p /home/${USER}/bin && \
    GOBIN=/home/${USER}/bin go install --mod=readonly google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \
    GOBIN=/home/${USER}/bin go install --mod=readonly google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
    GOBIN=/home/${USER}/bin go install --mod=readonly github.com/go-swagger/go-swagger/cmd/swagger@v0.28.0 && \
    GOBIN=/home/${USER}/bin go install github.com/github-release/github-release@v0.10.0


RUN \
    GOBIN=/home/${USER}/bin go install --mod=readonly github.com/karlmutch/duat/cmd/semver@0.16.0 && \
    GOBIN=/home/${USER}/bin go install --mod=readonly github.com/karlmutch/duat/cmd/stencil@0.16.0

CMD /bin/bash -C ./all-build.sh
