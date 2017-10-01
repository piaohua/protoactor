/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-01-23 10:34
 * Filename      : error.go
 * Description   : 前端请求异常的错误号和对应的错误信息
 *********************************************************/
package errorcode

const (
	DataFormatError       = 13001 // 请求的数据格式错误
	UserLoginAleady       = 13002 // 用户已经登录
	UsernameOrPwdError    = 13009 // 用户名或者密码错误
	UsernameEmpty         = 13010 // 用户名为空
	RegistError           = 13012 // 注册失败
	PwdEmpty              = 13011 // 密码为空
	MiguLoginError        = 13013 // 咪咕登录失败
	NicknameExist         = 13014 // 昵称已经存在
	EmailRegisted         = 13015 // 邮箱已经注册
	PhoneRegisted         = 13016 // 手机已经注册
	PhoneNumberError      = 13017 // 不是手机号码
	PwdFormatError        = 13018 // // 验证只能由数字字母下划线组成的5-17位密码字符串
	PhoneNumberEnpty      = 13019 // // 电话号码为空
	UserDataNotExist      = 14001 // 用户数据不存在
	OtherLoginThisAccount = 14002 // 您的帐号在其它地方登录
	WechatLoingFailReAuth = 14003 // 微信登录失败，请重新授权
	GetWechatUserInfoFail = 14004 // 获取微信用户数据失败

	NotInRoomCannotLeave   = 20002 // 不在房间，没有离开房间这一说
	GameStartedCannotLeave = 20001 // 正在牌局中不能离开
	NotYourTurn            = 20003 // 没轮到你打牌
	CardNotExistInYourHand = 20006 // 你没有这张牌
	CardValueZero          = 20007 // 牌值为0
	SearchRoomError        = 20008 // 房间匹配错误
	IndexOutofRange        = 20010 // 插牌位置超出
	InTheRoomAlready       = 20012 // 已经在房间了，请您请求断线重连接口
	YouNotExistPengAndKong = 20014 // 你不存在碰杠
	YouNotExistHu          = 20016 // 你没听牌不能胡
	NotInRoom              = 20018 // 你不在房间,针对房间的一切操作无效
	NotEnoughCoin          = 20019 // 金币不足
	EnoughCoin             = 20020 // 你没达到领取破产救济金的条件
	NotBankruptCoin        = 20021 // 你今天的破产补助已经全部领取完了
	CantTianTing           = 20022 // 已经开始打牌，不能天听
	StartedNotKick         = 20023 // 已经开始游戏不能踢人

	CreateRoomFail          = 30012 //创建房间失败
	RoomTypeNotExist        = 30013 //房间类型不存在
	RoomNotExist            = 30016 //房间不存在
	RoomFull                = 30018 //房间已满
	NotInPrivateRoom        = 30019 //玩家不在私人房间
	InTherPrivateRoom       = 30020 //玩家在私人房间
	NotonDesk               = 30021 //你不在牌桌座位上
	CannotCreatePrivateRoom = 30023 //不能创建非私人房间
	ReadyYet                = 30024 //你已经按过准备游戏了
	AudienceCannotOperate   = 30025 //你是观众不能对牌局操作
	NoRecordInSocial        = 30026 //  你没有在该私人房的打牌纪录
	NoRecordInRoom          = 30027 //该房间还没有打牌纪录
	GameHasBegun            = 30028 //牌局已经开始不能进入牌桌
	TingError               = 30029 // 玩家停牌出错
	AlreadyHasPrivateRoom   = 30030 //已经有私人房间
	MaCountIEllegal         = 30031 //马的数量不合法，合法范围(0-8) ,且是双数
	RunningNotVote          = 30032 //牌局已经开始不能投票
	VotingCantLaunchVote    = 30033 //房间里已经有玩家发起投票了
	NotVoteTime             = 30034 //先有人发起才能投票

	GlobalConfigError = 40001 //全局数据请求错误
	NameTooLong       = 40002 //取名太长了
	NameTooShort      = 40003 //取名太短了
	SexValueRangeout  = 40004 //性别取值错误
	FeedfackError     = 40005 //反馈失败
	NotEnoughDiamond  = 40007 //钻石不足
	NoticeListEnpty   = 40008 //没有公告
	DataOutOfRange    = 40009 //数据超出范围

	NotExistRankingCoin     = 50002 //目前没有财富总排行
	NotExistRankingCoinGain = 50003 //目前没有昨日收入排行
	NotExistRankingWin      = 50004 //目前没有胜局排行
	NotExistRankingExp      = 50005 //目前没有经验排行
	NotInRanking            = 50006 //为上榜没有奖励
	AlreadyRanking          = 50007 //已领取奖励

	SignFail    = 60001 //  签到失败，你可能今天已经签过到 了
	SignDayFail = 60002 //  你未签过到，你的签到列表为空
	SignReFail  = 60003 //  补签失败，全部签到or次数不够

	TaskIdError        = 61001 //  任务ID错误
	TaskRewardFail     = 61002 //  领取奖励失败or已经领取
	ActivityIdError    = 61003 //  活动ID错误
	ActivityRewardFail = 61004 //  活动领取奖励失败or已经领取
	TradeIdError       = 61005 //  兑换ID错误or其他
	TradeNotEnough     = 61006 //  兑换失败/全部兑换完
	TradeNotEnoughEx   = 61007 //  兑换失败/兑换券不足
	TradeNotNumber     = 61008 //  兑换失败/兑换次数用完
	TradeListEmpty     = 61009 //  兑换列表为空
	TradeRecordEmpty   = 61010 //  兑换记录列表为空

	IpayOrderFail  = 62001 // 爱贝支付下单失败
	AppleOrderFail = 62002 // 苹果支付下单失败

	PostboxEmpty     = 69001 // 你的邮箱没有邮件
	PostNotExist     = 69002 // 邮件不存在
	AppendixNotExist = 69003 // 邮件没有附件

	NotExistInMatch    = 68001 // 你不在该比赛场
	CantEntryMatch     = 68002 //不能进入该比赛场/时间限制
	CostNotEnough      = 68003 //不能进入消耗不足
	PrivateRecordEmpty = 68004 // 没有私人局记录
)
