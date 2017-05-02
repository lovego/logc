#!/bin/bash

main() {
    org=$1
    if test -z "$org"; then
        echo "org is required"
        exit
    fi
    path=$(pwd)
    domain=$2
    logc_path="/usr/local/bin"
    logc_log_path="/var/log/logc"

    if [ ! -f $logc_path/logc ]; then
        download $logc_path
    fi
    make_log_dir $logc_log_path

    cd $path && nohup logc $org $domain > $logc_log_path/logc-$org.log 2>&1 &
}

download() {
    logc_url="http://github.com/lovego/logc/raw/master/release/logc"

    mkdir -p $1
    sudo wget -O $1/logc $logc_url
    cd $1
    sudo chmod +x logc
}

make_log_dir() {
    sudo mkdir -p $1
    user=$(id -un)
    group=$(id -gn)
    sudo chown $user:$group $1
}

main $1 $2
