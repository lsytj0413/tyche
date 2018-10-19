#!/bin/bash

tyche_start() {
    echo "Start tyche server..."
    docker run -d -p :443:443 --restart always --name tyche -v /root/keys:/keys tyche:0.0.1 tyche -wx-appid ${WXAPPID} -wx-appsecret ${WXAPPSECRET} -wx-token ${WXTOKEN} -wx-aeskey ${WXAESKEY} -listen-client-url "https://0.0.0.0:443" -client-key-file /keys/www.soren.vip.key -client-cert-file /keys/www.soren.vip.pem

    if [ $? -ne 0 ]; then
        echo "Start tyche Failed"
        exit 1
    fi
}

tyche_stop() {
    echo "Stop tyche server..."
    docker stop $(docker ps -aq -f "name=tyche")
    docker rm $(docker ps -aq -f "name=tyche")
}

usage() {
    echo "use start|stop|restart arguments"
}

case $1 in
    start)    tyche_start
              exit
              ;;
    stop)     tyche_stop
              exit
              ;;
    restart)  tyche_stop
              tyche_start
              exit
              ;;
    *)        usage
              exit 1
              ;;
esac