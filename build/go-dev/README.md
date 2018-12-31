# go-dev

go-dev is a base image for building go code

See [#98](https://github.com/dyweb/gommon/issues/98)

- Docker Hub https://hub.docker.com/r/dyweb/go-dev

## Use

Use it as base in you multi stage build, the image is big, you don't want to use it directly unless 
you need a full go environment to run your code (looking at ginkgo)

TODO: example using multi stage build using dep and go mods

Use it as a playground, run the container in interactive mode and remove it after you are done

````bash
docker run --rm -it --entrypoint /bin/bash dyweb/go-dev:1.11.4 
````

## Build

````bash
make build
# need to login to dockerhub and be member
make push
````