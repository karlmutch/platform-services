#!/bin/bash -e
[ -z "$USER" ] && echo "env variable USER must be set" && exit 1;
nerdctl build -t platform-services:latest --build-arg USER=$USER --build-arg USER_ID=`id -u $USER` --build-arg USER_GROUP_ID=`id -g $USER` .
docker_name=`petname`
nerdctl run --name $docker_name -e GITHUB_TOKEN=$GITHUB_TOKEN -v $GOPATH:/project platform-services 
exit_code=`nerdctl inspect $docker_name --format='{{.State.ExitCode}}'`
if [ $exit_code -ne 0 ]; then
    echo "Error" $exit_code
    exit $exit_code
fi

echo "Build Done" ;

go install github.com/karlmutch/duat/cmd/semver@0.16.0
version=`$GOPATH/bin/semver`

for dir in cmd/*/ ; do
    base="${dir%%\/}"
    base="${base##*/}"
    if [ "$base" == "cli-experiment" ] ; then
        continue
    fi
    if [ "$base" == "cli-downstream" ] ; then
        continue
    fi
    cd $dir
    nerdctl build --namespace k8s.io -t $base:$version .
    cd -
done

exit 0

./push.sh
