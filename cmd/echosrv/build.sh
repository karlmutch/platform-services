#!/bin/bash -x

if ( find /project -maxdepth 0 -empty | read v );
then
  echo "source code must be mounted into the /project directory"
  exit 990
fi

set -e
set -o pipefail

[ -e internal/gen/echosrv ] || mkdir -p internal/gen/echosrv
#[ -e vendor/github.com/karlmutch/platform-services/internal ] || mkdir -p vendor/github.com/karlmutch/platform-services/internal
#[ -e vendor/github.com/karlmutch/platform-services/internal/gen ] || ln -s `pwd`/internal/gen vendor/github.com/karlmutch/platform-services/internal/gen
protoc -Icmd/echosrv -I/usr/include/google --plugin=`which protoc-gen-go` --go_out=./internal/gen/echosrv --go_opt=paths=source_relative --plugin=`which protoc-gen-go-grpc` --go-grpc_out=./internal/gen/echosrv cmd/echosrv/echosrv.proto --go-grpc_opt=paths=source_relative

if [ "$1" == "gen" ]; then
    exit 0
fi

export HASH=`git rev-parse HEAD`

go mod vendor
go mod tidy

mkdir -p cmd/echosrv/bin
go build -asmflags -trimpath -ldflags "-X github.com/karlmutch/platform-services/internal/version.GitHash=$HASH" -o cmd/echosrv/bin/echosrv cmd/echosrv/*.go
go build -asmflags -trimpath -ldflags "-X github.com/karlmutch/platform-services/internal/version.GitHash=$HASH" -race -o cmd/echosrv/bin/echosrv-race cmd/echosrv/*.go
if ! [ -z "${TRAVIS_TAG}" ]; then
    if ! [ -z "${GITHUB_TOKEN}" ]; then
        github-release release --user karlmutch --repo platform-services --tag ${TRAVIS_TAG} --pre-release || true
        github-release upload --user karlmutch --repo platform-services  --tag ${TRAVIS_TAG} --name echosrv --file cmd/echosrv/bin/echosrv
    fi
fi
