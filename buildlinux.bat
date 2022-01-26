
set GOPATH=D:\\workspace_sgswx
set GOARCH=amd64

set GOOS=linux

go build -o ./bin/all_in_one ./cmd/all-in-one
go build -o ./bin/master ./cmd/master
go build -o ./bin/gatesvr ./cmd/gatesvr
go build -o ./bin/authsvr ./cmd/authsvr
go build -o ./bin/lobbysvr ./cmd/lobbysvr
go build -o ./bin/entitysvr ./cmd/entitysvr
go build -o ./bin/gamesvr ./cmd/gamesvr
go build -o ./bin/emailsvr ./cmd/emailsvr
go build -o ./bin/friendsvr ./cmd/friendsvr
go build -o ./bin/adminsvr ./cmd/adminsvr
go build -o ./bin/shopsvr ./cmd/shopsvr
go build -o ./bin/paysvr ./cmd/paysvr
go build -o ./bin/jobsvr ./cmd/jobsvr
go build -o ./bin/ranksvr ./cmd/ranksvr
go build -o ./bin/accountsvr ./cmd/accountsvr
go build -o ./bin/activitysvr ./cmd/activitysvr
go build -o ./bin/wechatsvr ./cmd/wechatsvr
go build -o ./bin/logsvr ./cmd/logsvr

go build -o ./bin/listsvr ./cmd/listsvr
