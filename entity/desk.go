package entity

import (
	"time"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//牌桌数据
type Desk struct {
	*DeskData          //房间类型基础数据
	Dealer    uint32   //庄家的座位
	Dice      uint32   //骰子
	Cards     []uint32 //没摸起的海底牌
	//-------------
	Lian  uint32 //私人房间连庄数
	Round uint32 //私人房间打牌局数
	//-------------
	Vote   uint32            //投票发起者座位号
	Votes  map[uint32]uint32 //投票同意解散de玩家座位号列表
	VoteT  *time.Timer       //投票定时器
	Record utils.Array       //打牌记录
	//-------------
	//操作提示优先级 : (胡)-(碰,杠)-(吃)
	//如果是同一玩家 : (胡,碰,杠), (胡,碰,杠,吃)
	//如果多个玩家胡 : (胡)-(碰,杠)-(吃)
	//map[位置]操作值
	Hu       map[uint32]uint32
	Pongkong map[uint32]uint32
	Chow     map[uint32]uint32
	Skip     map[uint32]bool //过圈
	//-------------
	Discard uint32 //出牌
	Draw    uint32 //模牌
	Seat    uint32 //当前模牌|出牌位置
	Timer   int    //计时
	Operate int    //操作状态
	Huing   int    //一炮多响胡牌
	Laipi   uint32 //赖皮
	Laizi   uint32 //赖子
	//-------------
	State   bool //房间状态
	Kong    bool //是否杠牌出牌
	Laidraw bool //是否出赖摸牌
	//-------------
	CloseCh chan bool //关闭通道
	//-------------
	Trusteeship map[uint32]bool //托管
	Ready       map[uint32]bool //是否准备
	//-------------
	OutCards  map[uint32][]uint32 //海底牌
	PongCards map[uint32][]uint32 //碰牌
	KongCards map[uint32][]uint32 //杠牌
	ChowCards map[uint32][]uint32 //吃牌(8bit-8-8)
	HandCards map[uint32][]uint32 //手牌
	OutLaizi  map[uint32]int      //
}

//牌桌数据
type DeskData struct {
	Rid     uint32               `xorm:"Int"`         //房间ID
	Rtype   uint32               `xorm:"Int"`         //房间类型
	Rname   string               `xorm:"Varchar(20)"` //房间名字
	Cid     string               `xorm:"Varchar(20)"` //房间创建人
	Expire  uint32               `xorm:"Int"`         //牌局设定的过期时间
	Code    string               `xorm:"Varchar(6)"`  //房间邀请码
	Rounds  uint32               `xorm:"Int"`         //总牌局数
	Round   uint32               `xorm:"Int"`         //已经打牌局数
	Ante    uint32               `xorm:"Int"`         //私人房底分
	Payment uint32               `xorm:"Int"`         //付费方式1=AA or 0=房主支付
	Cost    uint32               `xorm:"Int"`         //创建消耗
	Pao     uint32               `xorm:"Int"`         //1=大牌玩法,0=不能炮胡
	Count   uint32               `xorm:"Int"`         //房间限制数量2,4
	CTime   uint32               `xorm:"Int"`         //创建时间
	Score   map[string]int32     `xorm:"Text"`        //用户战绩积分key:userid,val:score
	Roles   map[string]uint32    `xorm:"Text"`        //用户位置key:userid,val:seat
	Seats   map[uint32]*DeskUser //用户位置key:seat,val:DeskUser
	Pid     *actor.PID           //房间进程ID
}

//牌桌玩家数据
type DeskUser struct {
	Userid   string     //玩家ID
	Nickname string     //玩家昵称
	Sex      uint32     //玩家性别
	Photo    string     //玩家头像地址
	Pid      *actor.PID //玩家进程ID
}
