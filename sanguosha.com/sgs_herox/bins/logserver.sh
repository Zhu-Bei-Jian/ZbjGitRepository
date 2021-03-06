#!/bin/bash

exefile="logserver"
service="logserver"

basePath=$(cd `dirname $0` ; pwd)
startCommand=" $basePath/$exefile --name $service"
pidPath="$basePath/var/run"
logPath="$basePath/var/log"
pidFile="$pidPath/$service.pid"
logFile="$logPath/$service.log"

start(){
    num=`ps -ef | grep "${startCommand}" | grep -v grep | wc -l`
    echo ${num}
    if [ ${num} -eq 0 ]
    then
        echo "starting..."
        echo ${startCommand}
        if [ ! -d "$logPath" ]; then
            mkdir -p "$logPath"
        fi
        #nohup ${startCommand} > /dev/null 2>&1 &
        nohup ${startCommand} > /dev/null 2>>${logFile} &
        if [ $? -ne 0 ]
        then
            echo "start failed, please check the log"
            exit $?
        else
			if [ ! -d "$pidPath" ]; then
				mkdir -p "$pidPath"
			fi
            echo $! > ${pidFile}
            echo "start success"
        fi
    else
        echo "$service is already running"
    fi
}
stop(){
    if [ ! -f "$pidFile" ]; then
        echo "$service is not running"
    else
        echo "stopping..."
        echo ${startCommand}
        PROCESS=`ps -ef | grep "${startCommand}" | grep -v grep | awk '{print $2}'`
        for i in ${PROCESS}
        do
            kill -SIGINT ${i}
        done
        #kill -SIGINT `cat ${pidFile}`
        if [ $? -ne 0 ]
        then
            echo "stop failed, may be $service isn't running"
            exit $?
        else
            rm -rf ${pidFile}
            echo "stop success"
        fi
    fi
}

status(){
    num=`ps -ef | grep "${startCommand}" | grep -v grep | wc -l`
    if [ ${num} -eq 0 ]
    then
        echo "$service isn't running"
    else
        echo "$service is running"
    fi
}

case $1 in
    start)      start ;;
    stop)       stop ;;
    status)     status ;;
    *)          echo "Usage: $0 {start|stop|status}" ;;
esac

exit 0
