package account

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/redisclient"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"time"
)

const (
	UserSummaryExpireTime = time.Hour
	KeyUserSummary        = "user:summary:%d"
)

type UserCacheManager struct {
	client *redisclient.RedisClient
}

func newUserCacheManager(cfg *config.RedisConfig) (*UserCacheManager, error) {
	client := &redisclient.RedisClient{}

	err := client.Init(cfg)
	if err != nil {
		return nil, err
	}

	return &UserCacheManager{client: client}, nil
}

func (p *UserCacheManager) Close() {
	p.client.Close()
}

func (p *UserCacheManager) UpdateUserSummary(userId uint64, summary *gamedef.UserSummary) error {
	data, err := proto.Marshal(summary)
	if err != nil {
		return err
	}
	_, err = p.client.SetEx(userSummaryRedisKey(userId), data, UserSummaryExpireTime)
	return err
}

func (p *UserCacheManager) GetUserSummaries(userIds []uint64) (summaries map[uint64]*gamedef.UserSummary, unfindUserIds []uint64, err error) {
	summaries = make(map[uint64]*gamedef.UserSummary, len(userIds))

	if len(userIds) == 0 {
		return summaries, nil, nil
	}

	params := make([]string, 0, len(userIds))
	for _, userId := range userIds {
		params = append(params, userSummaryRedisKey(userId))
	}

	summaryDatas, err := p.client.MGetBytes(params...)
	if err != nil {
		return nil, nil, err
	}

	for i, data := range summaryDatas {
		userId := userIds[i]

		s, err := unmarshalUserSummary(userId, data)
		if err != nil {
			unfindUserIds = append(unfindUserIds, userId)
			continue
		}

		summaries[userId] = s
	}
	return
}

func userSummaryRedisKey(userId uint64) string {
	return fmt.Sprintf(KeyUserSummary, userId)
}

const NewestSummaryVersion = gamedef.UserSummary_VerInit

func unmarshalUserSummary(userId uint64, data []byte) (*gamedef.UserSummary, error) {
	if data == nil {
		return nil, errors.New("data == nil")
	}
	var summary gamedef.UserSummary
	err := proto.Unmarshal(data, &summary)
	if err != nil {
		return nil, err
	}

	if summary.UserBrief.UserID != userId || summary.Version != NewestSummaryVersion {
		return nil, errors.New("not the lastest value")
	}
	return &summary, nil
}
