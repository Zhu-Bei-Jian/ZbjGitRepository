
@echo off
set GOOS=%1
set BUILD_OUTPUT=%2
set BUILD_EXT=""

if "%GOOS%"=="windows" (
    set BUILD_EXT=".exe"
)

::测试用
go build -mod=vendor -o %BUILD_OUTPUT%/all-in-one%BUILD_EXT% ./cmd/all-in-one
echo build success: all-in-one

::集群内
go build -mod=vendor -o %BUILD_OUTPUT%/master%BUILD_EXT% ./cmd/master
echo build success: master
go build -mod=vendor -o %BUILD_OUTPUT%/gatesvr%BUILD_EXT% ./cmd/gatesvr
echo build success: gate
go build -mod=vendor -o  %BUILD_OUTPUT%/authsvr%BUILD_EXT% ./cmd/authsvr
echo build success: auth
go build -mod=vendor -o %BUILD_OUTPUT%/lobbysvr%BUILD_EXT% ./cmd/lobbysvr
echo build success: lobby
go build -mod=vendor -o %BUILD_OUTPUT%/entitysvr%BUILD_EXT% ./cmd/entitysvr
echo build success: entity
go build -mod=vendor -o %BUILD_OUTPUT%/gamesvr%BUILD_EXT% ./cmd/gamesvr
echo build success: game
go build -mod=vendor -o  %BUILD_OUTPUT%/adminsvr%BUILD_EXT% ./cmd/adminsvr
echo build success: admin
go build -mod=vendor -o  %BUILD_OUTPUT%/accountsvr%BUILD_EXT% ./cmd/accountsvr
echo build success: account

echo all done!




