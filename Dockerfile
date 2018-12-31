# This Dockerfile is a demo of using go-dev to build a go binary using multi stage build
# It is based on
# https://docs.docker.com/v17.09/engine/userguide/eng-image/dockerfile_best-practices/#use-multi-stage-builds
FROM dyweb/go-dev:1.11.4 as builder

LABEL maintainer="contact@dongyue.io"

ARG PROJECT_ROOT=/go/src/github.com/dyweb/gommon/

WORKDIR $PROJECT_ROOT

# Gopkg.toml and Gopkg.lock lists project dependencies
# These layers will only be re-built when Gopkg files are updated
COPY Gopkg.lock Gopkg.toml $PROJECT_ROOT
RUN dep ensure -v -vendor-only

# Copy all project and build it
COPY . $PROJECT_ROOT
RUN make install

# NOTE: use ubuntu instead of alphine
#
# When using alpine I saw standard_init_linux.go:190: exec user process caused "no such file or directory",
# because I didn't compile go with static flag
# https://stackoverflow.com/questions/49535379/binary-compiled-in-prestage-doesnt-work-in-scratch-container
FROM ubuntu:18.04 as runner
LABEL maintainer="contact@dongyue.io"
LABEL github="github.com/dyweb/gommon"
WORKDIR /usr/bin
COPY --from=builder /go/bin/gommon .
ENTRYPOINT ["gommon"]
CMD ["help"]