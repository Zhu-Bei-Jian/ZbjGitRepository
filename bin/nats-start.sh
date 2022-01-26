#!/bin/bash

nohup ./nats-streaming-server -mb 10000MB -mm 0 -store file -dir datastore --file_sync=false -ft_group "ft" -cluster nats://localhost:6222 -routes nats://localhost:6223 -p 5222 -m 8222 >5222_info.log 2>5222_err.log &
nohup ./nats-streaming-server -mb 10000MB -mm 0 -store file -dir datastore --file_sync=false -ft_group "ft" -cluster nats://localhost:6223 -routes nats://localhost:6222 -p 5223 -m 8223 >5223_info.log 2>5223_err.log &
