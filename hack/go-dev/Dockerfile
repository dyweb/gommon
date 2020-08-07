ARG BASE_IMAGE=ubuntu:18.04
FROM $BASE_IMAGE

LABEL maintainer="contact@dongyue.io"

# https://github.com/dyweb/gommon/issues/98
RUN \
    export DEBIAN_FRONTEND=noninteractive \
    && apt-get update \
    && apt-get upgrade -y --no-install-recommends \
    && apt-get install -y --no-install-recommends \
        bash \
        build-essential \
        ca-certificates \
        curl \
        wget \
        git-core \
        ssh-client \
        man \
        vim \
        zip \
        unzip \
        tmux \
        netcat \
        telnet \
        tree \
        strace \
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
# TODO: do we need chmod -R 777 like https://github.com/docker-library/golang/blob/master/Dockerfile-debian.template
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin"

ARG BUILD_GO_VERSION=1.11.2

# TODO: Deprecate glide and dep, we are all go mod now

ENV GO_VERSION=$BUILD_GO_VERSION \
    GLIDE_VERSION=v0.13.2

RUN \
    curl -L https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz | tar -C /usr/local -xz \
    && curl -sSL https://github.com/Masterminds/glide/releases/download/$GLIDE_VERSION/glide-$GLIDE_VERSION-linux-amd64.tar.gz \
       | tar -vxz -C /usr/local/bin --strip=1 \
    && rm /usr/local/bin/README.md /usr/local/bin/LICENSE

# TODO: install dep may have problem when go mod is enabled ...
# dep releases are way behind master, so we install from source
RUN go get -u -v github.com/golang/dep/cmd/dep \
    && go get -u -v golang.org/x/tools/cmd/goimports
# TODO: should clean up the mod cache

WORKDIR $GOPATH