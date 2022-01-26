package admin

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/proto/smsg"
	"strconv"
	"sync"
	"time"
)

var (
	metrics_onlineVec, metrics_gameVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sgs_online_count",
		Help: "online player count in server",
	}, []string{"serverid"}), promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sgs_game_count",
		Help: "game count in server",
	}, []string{"serverid"})

	metrics_roomCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "sgs_room_count",
		Help:        "room count in server",
		ConstLabels: nil,
	})

	metrics_rpc_millseconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "",
		Subsystem:   "",
		Name:        "sgs_rpc_millseconds",
		Help:        "",
		ConstLabels: nil,
		Buckets:     []float64{10, 50, 200, 500},
	}, []string{"serverid"})
)

func startRecordMetrics() {
	serverTypes := []appframe.ServerType{
		sgs_herox.SvrTypeLobby,
		sgs_herox.SvrTypeGate,
		sgs_herox.SvrTypeGame,
	}

	var serverId2Metrics = make(map[uint32][]*smsg.AdAllRespMetrics_Metrics)

	util.SafeGo(func() {
		for {
			time.Sleep(time.Second * 10)
			wg := sync.WaitGroup{}

			availableServerIds := make(map[uint32][]*smsg.AdAllRespMetrics_Metrics)

			start := time.Now().UnixNano()
			for _, serverType := range serverTypes {
				serverIds := app.GetAvailableServerIDs(serverType)
				wg.Add(len(serverIds))

				for _, v := range serverIds {
					serverId := v
					app.GetServer(serverId).ReqSugar(&smsg.AdAllReqMetrics{
						ReqTime: start,
					}, func(resp *smsg.AdAllRespMetrics, err error) {
						defer wg.Done()

						if err != nil {
							logrus.WithError(err).Error("request metrics")
							return
						}
						elapse := (time.Now().UnixNano() - start) / 1e6
						metrics_rpc_millseconds.With(prometheus.Labels{"serverid": strconv.FormatInt(int64(serverId), 10)}).Observe(float64(elapse))

						availableServerIds[serverId] = resp.Metrics
					}, time.Second*10)
				}
			}

			wg.Wait()

			//不可用的服务器，指标参数清0
			for serverId, metrics := range serverId2Metrics {
				_, exist := availableServerIds[serverId]
				if exist {
					continue
				}
				for _, v := range metrics {
					v.Value = 0
				}
			}

			for serverId, v := range availableServerIds {
				serverId2Metrics[serverId] = v
			}

			for serverId, metrics := range serverId2Metrics {
				for _, met := range metrics {
					switch met.Key {
					case smsg.AdAllRespMetrics_OnlineCount:
						metrics_onlineVec.With(prometheus.Labels{"serverid": strconv.FormatInt(int64(serverId), 10)}).Set(float64(met.Value))
					case smsg.AdAllRespMetrics_GameCount:
						metrics_gameVec.With(prometheus.Labels{"serverid": strconv.FormatInt(int64(serverId), 10)}).Set(float64(met.Value))
					case smsg.AdAllRespMetrics_RoomCount:
						metrics_roomCount.Set(float64(met.Value))
					default:
						continue
					}
				}
			}
		}
	})
}
