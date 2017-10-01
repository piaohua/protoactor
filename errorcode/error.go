package errorcode

const (
	AccountLegal  = 13001 // 账号不合格
	LoginFail     = 13002 // 登录失败
	NicknameLegal = 13003 // 用户名不合格
	AccountExist  = 13004 // 账号已经存在
	PasswordLegal = 13007 // 密码不合格
	PasswordError = 13008 // 密码错误
	RegistFail    = 13009 // 注册失败

	NotInRoom        = 13010 // 你不在房间,针对房间的一切操作无效
	RoomNotExist     = 13011 // 房间不存在
	NotEnoughDiamond = 13012 // 钻石不足
	NotEnoughCoin    = 13013 // 金币不足
	CreateRoomFail   = 13014 // 创建房间失败
	NotYourTurn      = 13015 // 没轮到你打牌
	InTheGame        = 13016 // 游戏中不能离开
	InTheVote        = 13017 // 投票中不能准备
	NotInTheGame     = 13018 // 不在游戏中
)
