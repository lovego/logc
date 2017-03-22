#!/bin/bash

main() {
    token=$1
    if test -z "$token"; then
        echo "token is required"
        exit
    fi
    domain=$2
    logc_path="/usr/local/bin"
    logc_log_path="/var/log/logc"

    if [ ! -f $logc_path/logc ]; then
        download $logc_path
    fi
    make_log_dir $logc_log_path

    nohup logc $token $domain > $logc_log_path/logc-$token.log 2>&1 &
}

download() {
    logc_url="http://github.com/logc-monitor/logc/raw/master/logc"

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
