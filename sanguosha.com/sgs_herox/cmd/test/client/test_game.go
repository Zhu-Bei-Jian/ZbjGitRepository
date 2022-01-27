package client

import (
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"time"
)

func listenGameMsg(client *Client) {
	client.ListenMsg((*cmsg.SNoticeGameReady)(nil), func(msg proto.Message) {
		time.AfterFunc(1*time.Second, func() {
			client.SendMsg(&cmsg.CReqGameScene{})
		})
	})

	client.ListenMsg((*cmsg.SRespGameScene)(nil), func(msg proto.Message) {
		notice := msg.(*cmsg.SRespGameScene)
		client.seatid = notice.MySeatId
	})

	client.ListenMsg((*cmsg.SNoticeOp)(nil), func(msg proto.Message) {
		//notice := msg.(*cmsg.SNoticeOp)
		//switch notice.OpType {
		//case cmsg.SNoticeOp_DescribeInVoice:
		//	time.AfterFunc(3*time.Second, func() {
		//		client.SendMsg(&cmsg.CReqDescribe{})
		//	})
		//case cmsg.SNoticeOp_Describe:
		//	s := generateRandomString(5)
		//	time.AfterFunc(3*time.Second, func() {
		//		client.SendMsg(&cmsg.CReqDescribe{
		//			Desc: s,
		//		})
		//	})
		//case cmsg.SNoticeOp_Vote:
		//	seats := notice.GetTargetSeatIds()
		//	if len(seats) == 0 {
		//		logrus.Warn("Op seats is empty")
		//		return
		//	}
		//	seatsCanVote := getSeatsCanVote(seats, client.seatid)
		//	if len(seatsCanVote) == 0 {
		//		logrus.Warn("Op seats is empty")
		//		return
		//	}
		//	randomIndex := gameutil.RandomBetween(0, int32(len(seatsCanVote)-1))
		//	time.AfterFunc(3*time.Second, func() {
		//		client.SendMsg(&cmsg.CReqVote{
		//			SeatId: seats[randomIndex],
		//		})
		//	})
		//}
	})

	client.ListenMsg((*cmsg.SNoticeEnterPhase)(nil), func(msg proto.Message) {
		notice := msg.(*cmsg.SNoticeEnterPhase)
		if notice.Phase == gamedef.GamePhase_End {
			time.AfterFunc(1*time.Second, func() {
				client.SendMsg(&cmsg.CReqRoomLeave{})
				client.SendMsg(&cmsg.CReqRoomQuickJoin{
					GameMode: gameconf.GameModeTyp(generateRandomMode()),
				})
				client.SendMsg(&cmsg.CReqRoomReady{
					Ready: true,
				})
			})
		}
	})
}
func (c *Client) ClearGame() {
	if len(c.gameUuid) != 0 {
		c.gameMatch = false
		c.gameUuid = ""
		//logrus.Info("game over")
		pushCol("game_over", 1)
	}
}
func (c *Client) GoTestGame() {
	tickerReq := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-tickerReq.C:
				c.Post(func() {
					if len(c.gameUuid) == 0 {
						if !c.gameMatch {
							c.gameMatch = true
							time.AfterFunc(time.Duration(1+rand.Intn(5)), func() {
								//c.SendMsg(&cmsg.CReqMatch{
								//	Model:   9001,
								//	Section: 2002,
								//})
							})
						}
					}
					return
				})
			}
		}
	}()
}

func getSeatsCanVote(in []int32, myseat int32) []int32 {
	for i := 0; i < len(in); {
		if in[i] == myseat {
			in = append(in[:i], in[i+1:]...)
		} else {
			i++
		}
	}
	return in
}

func generateRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func generateRandomMode() int32 {
	return gameutil.RandomBetween31n(int32(gameconf.GameModeTyp_MGTSpyText), int32(gameconf.GameModeTyp_MGTSpyVoice))
}
