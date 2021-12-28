#!/bin/bash -x

if ( find /project -maxdepth 0 -empty | read v );
then
  echo "source code must be mounted into the /project directory"
  exit -1
fi

set -e
set -o pipefail

export HASH=`git rev-parse HEAD`

[ -e internal/gen/downstream ] || mkdir -p internal/gen/downstream
#[ -e vendor/github.com/karlmutch/platform-services/internal ] || mkdir -p vendor/github.com/karlmutch/platform-services/internal
#[ -e vendor/github.com/karlmutch/platform-services/internal/gen ] || ln -s `pwd`/internal/gen vendor/github.com/karlmutch/platform-services/internal/gen
protoc -Icmd/downstream -I/usr/include/google --plugin=`which protoc-gen-go` --go_out=./internal/gen/downstream --go_opt=paths=source_relative --plugin=`which protoc-gen-go-grpc` --go-grpc_out=./internal/gen/downstream cmd/downstream/downstream.proto --go-grpc_opt=paths=source_relative

if [ "$1" == "gen" ]; then
    exit 0
fi

go mod vendor
go mod tidy

export SEMVER=`semver`
TAG_PARTS=$(echo $SEMVER | sed "s/-/\n-/g" | sed "s/\./\n\./g" | sed "s/+/\n+/g")
PATCH=""
for part in $TAG_PARTS
do
    start=`echo "$part" | cut -c1-1`
    if [ "$start" == "+" ]; then
        break
    fi
    if [ "$start" == "-" ]; then
        PATCH+=$part
    fi
done
flags='-X github.com/karlmutch/platform-services/internal/version.GitHash="$HASH" -X github.com/karlmutch/platform-services/internal/version.SemVer="$SEMVER"'
flags="$(eval echo $flags)"

mkdir -p cmd/downstream/bin
CGO_ENABLED=0 go build -asmflags -trimpath -ldflags "$flags" -o cmd/downstream/bin/downstream cmd/downstream/*.go
go build -asmflags -trimpath -ldflags "$flags" -race -o cmd/downstream/bin/downstream-race cmd/downstream/*.go
CGO_ENABLED=0 go test -ldflags "$flags" -coverpkg="." -c -o cmd/downstream/bin/downstream-run-coverage cmd/downstream/*.go
CGO_ENABLED=0 go test -ldflags "$flags" -coverpkg="." -c -o bin/downstream-test-coverage cmd/downstream/*.go
go test -ldflags "$flags" -race -c -o cmd/downstream/bin/downstream-test cmd/downstream/*.go

if [ -z "$PATCH" ]; then
    if ! [ -z "${SEMVER}" ]; then
        if ! [ -z "${GITHUB_TOKEN}" ]; then
            github-release release --user karlmutch --repo platform-services --tag ${SEMVER} --pre-release || true
            github-release upload --user karlmutch --repo platform-services  --tag ${SEMVER} --name downstream --file cmd/downstream/bin/downstream
        fi
    fi
fi
