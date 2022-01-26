package accountservice

import (
	"errors"
	"fmt"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"time"
)

var RequestTimeout = 15 * time.Second
var ErrNotFound = errors.New("user not exist")

type Service struct {
	app *appframe.Application
}

func New(app *appframe.Application) *Service {
	app.RegisterService(sgs_herox.SvrTypeAccount, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeAccount))
	app.RegisterResponse((*smsg.PuAcRespUserSummary)(nil))
	app.RegisterResponse((*smsg.PuAcRespQueryUserID)(nil))

	return &Service{app: app}
}

func (p *Service) GetUserSummarySync(userIds []uint64) (map[uint64]*gamedef.UserSummary, error) {
	var resp smsg.PuAcRespUserSummary
	err := p.app.GetService(sgs_herox.SvrTypeAccount).CallSugar(&smsg.PuAcReqUserSummary{Userids: userIds}, &resp, RequestTimeout)
	if err != nil {
		return nil, err
	}

	if resp.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("PuAcReqUsersHeadInfo err:%d", resp.ErrCode))
	}
	return resp.Summaries, nil
}

func (p *Service) GetUserSummaryAsync(userIds []uint64, cbk func(map[uint64]*gamedef.UserSummary, error)) {
	p.app.GetService(sgs_herox.SvrTypeAccount).ReqSugar(&smsg.PuAcReqUserSummary{Userids: userIds}, func(resp *smsg.PuAcRespUserSummary, err error) {
		if err != nil {
			cbk(nil, err)
			return
		}

		if resp.ErrCode != 0 {
			cbk(nil, errors.New(fmt.Sprintf("PuAcReqUsersHeadInfo err:%d", resp.ErrCode)))
			return
		}

		cbk(resp.Summaries, nil)
	}, RequestTimeout)
}

func (p *Service) GetUserIdByUnionIdSync(unionId string, accountType gameconf.AccountLoginTyp) (uint64, error) {
	userIds, err := p.QueryUserIDSync(0, unionId, accountType, smsg.PuAcReqQueryUserID_ByUnionId)
	if err != nil {
		return 0, err
	}
	if len(userIds) == 0 {
		return 0, ErrNotFound
	}

	return userIds[0], nil
}

func (p *Service) GetUserIDByAccountSync(account string) (uint64, error) {
	userIds, err := p.QueryUserIDSync(0, account, 0, smsg.PuAcReqQueryUserID_ByAccount)
	if err != nil {
		return 0, err
	}
	if len(userIds) == 0 {
		return 0, ErrNotFound
	}

	return userIds[0], nil
}

func (p *Service) GetUserIdsByNickNameSync(nickname string, like bool) ([]uint64, error) {
	queryType := smsg.PuAcReqQueryUserID_ByNickName
	if like {
		queryType = smsg.PuAcReqQueryUserID_LikeNickName
	}
	userIds, err := p.QueryUserIDSync(0, nickname, 0, queryType)
	if err != nil {
		return nil, err
	}
	return userIds, nil
}

func (p *Service) GetUserIDByAccountAsync(account string, cbk func(uint64, error)) {
	p.QueryUserIDAsync(0, account, smsg.PuAcReqQueryUserID_ByAccount, func(userIds []uint64, err error) {
		if err != nil {
			cbk(0, err)
			return
		}
		if len(userIds) == 0 {
			cbk(0, ErrNotFound)
			return
		}

		cbk(userIds[0], nil)
	})
}

func (p *Service) GetUserIdsByNickNameAsync(nickname string, like bool, cbk func([]uint64, error)) {
	queryType := smsg.PuAcReqQueryUserID_ByNickName
	if like {
		queryType = smsg.PuAcReqQueryUserID_LikeNickName
	}

	p.QueryUserIDAsync(0, nickname, queryType, func(userIds []uint64, err error) {
		if err != nil {
			cbk(nil, err)
			return
		}

		cbk(userIds, nil)
	})
}

func (p *Service) QueryUserIDAsync(paramInt uint64, paramStr string, queryType smsg.PuAcReqQueryUserID_QueryType, cbk func([]uint64, error)) {
	p.app.GetService(sgs_herox.SvrTypeAccount).ReqSugar(&smsg.PuAcReqQueryUserID{
		QueryType: queryType,
		ParamInt:  paramInt,
		ParamStr:  paramStr,
	}, func(resp *smsg.PuAcRespQueryUserID, err error) {
		if err != nil {
			cbk(nil, err)
			return
		}

		if resp.ErrCode != 0 {
			cbk(nil, errors.New(fmt.Sprintf("PuAsRespUserIDByShowID err:%d", resp.ErrCode)))
			return
		}

		cbk(resp.UserIds, nil)
	}, RequestTimeout)
}

func (p *Service) QueryUserIDSync(paramInt uint64, paramStr string, accountType gameconf.AccountLoginTyp, queryType smsg.PuAcReqQueryUserID_QueryType) ([]uint64, error) {
	var resp smsg.PuAcRespQueryUserID
	err := p.app.GetService(sgs_herox.SvrTypeAccount).CallSugar(&smsg.PuAcReqQueryUserID{
		QueryType:   queryType,
		ParamInt:    paramInt,
		ParamStr:    paramStr,
		AccountType: accountType,
	}, &resp, RequestTimeout)

	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("PuAsRespUserIDByShowID err:%d", resp.ErrCode))
	}
	return resp.UserIds, nil
}
