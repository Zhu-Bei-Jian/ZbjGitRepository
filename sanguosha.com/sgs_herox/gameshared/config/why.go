package config

import "github.com/sirupsen/logrus"

//添加物品原因  ADD_what_BY_reason ADD_BY_reason
const (
	//------------商城相关---------------------
	ShopBuy int32 = 301 //商城购买

	//------------邮件----------------------------
	EmailAttach int32 = 400 //邮件符件 GetID:邮件ID

	//------------充值相关-------------------------
	RMBCharge int32 = 500 //充值 GetID:OrderID

	//------------------游戏相关------------------------------------
	PVPLastSeasonReward int32 = 600 //排位赛上赛季奖励
	PVPLevelReward      int32 = 601 //排位赛赛季段位奖励

	PVP1V1GuideameOver  int32 = 781 // 引导模式1v1游戏结束
	PVPQualifyGameOver  int32 = 780 //排位赛游戏结束
	PVPRPGGameOver      int32 = 792 //身份场模式结束
	PVPLandlordGameOver int32 = 793 // 斗地主模式结算

	RPGChangeHero int32 = 800 //身份场换将

	ChangeGeneral int32 = 801 // 换将卡换将
	ChangeCards   int32 = 802 // 手气卡换牌

	Cailanzi int32 = 803 // 使用菜篮子物品

	//------------------------活动任务相关--------------------------------
	TaskRewardReceive = 1000 //任务奖励领取 GetID:任务ID

	SignIn             int32 = 1060 //签到 GetID:第几个格子
	Login7DayReward    int32 = 1061 //七日登陆 GetID:对应天数
	Happy7DayBuy       int32 = 1062 //七天乐抢购
	InviteGetAward     int32 = 1063 //领取邀请奖励
	LuckyBag           int32 = 1064 //福袋
	LuckyBagGoldReward int32 = 1065 //福袋发起人元宝奖励
	FirstCharge        int32 = 1066 //首充奖励
	LuckySign          int32 = 1067 //幸运签
	ReceiveRelief      int32 = 1068 //领取救济金
	HeroGift           int32 = 1069 //新手福利-送武将
	WeChatAPPReward    int32 = 1070 // 添加到"我的小程序"奖励

	OpenBox = 1071 //开盒子

	// ----------------------分享-------------------------------------
	Share        int32 = 1100
	ShareGame    int32 = 1101 // 分享游戏结果
	ShareGeneral int32 = 1102 // 分享获得武将
	// ----------------------好友赠送-------------------------------------
	WechatFriendGive int32 = 1201 // 微信好友赠送

	//-----其int32它-------
	AddBox            int32 = 1601 //加宝箱（产出有变化）
	SyncFromMobileSgs int32 = 1602 // 手杀同步武将
	AddHeroDestory    int32 = 1603 //加英雄分解
	AddHeroTryDestory int32 = 1604 //加试用英雄分解
	AddSkinDestory    int32 = 1605 //加皮肤分解
	AddSkinTryDestory int32 = 1606 //int32加试用皮肤分解

	//-----------------其它----------------
	GMAdd = 1700 //GM加的
	GMDel = 1701 //GM删的

)

func GetWhyName(whyID int32) string {
	switch whyID {
	case ShopBuy:
		return "商城购买"
	case EmailAttach:
		return "邮件符件"
	case RMBCharge:
		return "充值"
	case PVPLastSeasonReward:
		return "排位赛上赛季奖励"
	case PVPLevelReward:
		return "排位赛赛季段位奖励"
	case PVP1V1GuideameOver:
		return "引导模式1v1游戏结束"
	case PVPQualifyGameOver:
		return "2v2游戏结束"
	case PVPRPGGameOver:
		return "身份场模式结束"
	case PVPLandlordGameOver:
		return "斗地主模式结束"
	case RPGChangeHero:
		return "身份场换将"
	case ChangeGeneral:
		return "换将卡换将"
	case ChangeCards:
		return "手气卡换牌"
	case TaskRewardReceive:
		return "任务奖励领取"
	case SignIn:
		return "签到"
	case Login7DayReward:
		return "七日登陆"
	case GMAdd:
		return "GM加的"
	case GMDel:
		return "GM删的"
	case Happy7DayBuy:
		return "七天乐抢购"
	case AddBox:
		return "加宝箱"
	case LuckyBag:
		return "福袋"
	case LuckyBagGoldReward:
		return "福袋发起人奖励"
	case FirstCharge:
		return "首充奖励"
	case LuckySign:
		return "幸运签"
	case ReceiveRelief:
		return "领取救济金"
	case Share:
		return "分享有奖"
	case ShareGame:
		return "分享游戏结果"
	case HeroGift:
		return "新手福利-送武将"
	case WeChatAPPReward:
		return "添加到我的小程序奖励"
	case OpenBox:
		return "开盒子"
	case SyncFromMobileSgs:
		return "手杀武将同步"
	case AddHeroDestory:
		return "加英雄分解"
	case AddHeroTryDestory:
		return "加试用英雄分解"
	case AddSkinDestory:
		return "加皮肤分解"
	case AddSkinTryDestory:
		return "加试用皮肤分解"
	default:
		logrus.WithField("whyID", whyID).Debug("whyName null")
	}
	return ""
}
