package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameutil"
)

const TableParkGameStatus_Start = 1
const TableParkGameStatus_Over = 2

type TableParkNotifier struct {
	config.TablePark
	gameId int32
}

func NewTableParkNotifier(gameId int32, config config.TablePark) *TableParkNotifier {
	return &TableParkNotifier{gameId: gameId, TablePark: config}
}

func (p *TableParkNotifier) NotifyGameStatusAsync(roomId int32, thirdAccountId string, status int32) {
	if !p.NoticeOpen {
		return
	}
	util.SafeGo(func() {
		err := p.notifyGameStatus(roomId, thirdAccountId, status)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"roomId":         roomId,
				"thirdAccountId": thirdAccountId,
			}).WithError(err).Error("tableparkNotifier notice gamestatus")
		}
	})
}

func (p *TableParkNotifier) NotifyGameDataAsync(thirdPartyId string, successNum, totalNum int32) {
	if !p.NoticeOpen {
		return
	}
	util.SafeGo(func() {
		err := p.notifyGameData(thirdPartyId, successNum, totalNum)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"thirdAccountId": thirdPartyId,
			}).WithError(err).Error("tableparkNotifier notice gamedata")
		}
	})
}

func (p *TableParkNotifier) notifyGameStatus(roomId int32, thirdPartyId string, status int32) error {
	type NotifyData struct {
		GameId     int32 `json:"gameId"`
		GameStatus int32 `json:"gameStatus"`
		RoomId     int32 `json:"roomId"`
		UserId     int64 `json:"userId"`
	}

	userId, err := gameutil.ToInt64(thirdPartyId)
	if err != nil {
		return err
	}

	d := NotifyData{
		GameId:     p.gameId,
		GameStatus: status,
		RoomId:     roomId,
		UserId:     userId,
	}

	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", p.UrlNoticeGameStatus, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	type Result struct {
		Status bool   `json:"status"`
		Code   int32  `json:"code"`
		Msg    string `json"msg"`
	}

	var ret Result
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return err
	}

	if !ret.Status {
		return fmt.Errorf("errCode %d msg:%s", ret.Code, ret.Msg)
	}

	return nil
}

func (p *TableParkNotifier) notifyGameData(thirdPartyId string, successNum, totalNum int32) error {
	type NotifyData struct {
		GameId     int32 `json:"gameId"`
		SuccessNum int32 `json:"successNum"`
		TotalNum   int32 `json:"totalNum"`
		UserId     int64 `json:"userId"`
	}

	userId, err := gameutil.ToInt64(thirdPartyId)
	if err != nil {
		return err
	}

	d := NotifyData{
		GameId:     p.gameId,
		SuccessNum: successNum,
		TotalNum:   totalNum,
		UserId:     userId,
	}

	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", p.UrlNoticeGameData, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	type Result struct {
		Status bool   `json:"status"`
		Code   int32  `json:"code"`
		Msg    string `json"msg"`
	}

	var ret Result
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return err
	}

	if !ret.Status {
		return fmt.Errorf("errCode %d msg:%s", ret.Code, ret.Msg)
	}

	return nil
}
