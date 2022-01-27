#!/bin/bash

GOPATH=/root/sgs_wx_release

# compile and build
#go build -o ./bins/gnatsd github.com/nats-io/gnatsd
go build -o ./bins/master ./cmd/master
go build -o ./bins/gatesvr ./cmd/gatesvr
go build -o ./bins/entitysvr ./cmd/entitysvr
go build -o ./bins/gamesvr ./cmd/gamesvr
go build -o ./bins/logserver ./cmd/logserver

go build -o ./bins/lobbysvr ./cmd/lobbysvr
go build -o ./bins/authsvr ./cmd/authsvr
go build -o ./bins/adminsvr ./cmd/adminsvr
go build -o ./bins/emailsvr ./cmd/emailsvr
go build -o ./bins/friendsvr ./cmd/friendsvr
go build -o ./bins/paysvr ./cmd/paysvr
go build -o ./bins/shopsvr ./cmd/shopsvr
go build -o ./bins/jobsvr ./cmd/jobsvr
go build -o ./bins/assistsvr ./cmd/assistsvr

go build -o ./bins/paycenter ./cmd/paycenter
go build -o ./bins/loadrecoder ./cmd/loadrecoder


chmod a+x ./bins/*svr


