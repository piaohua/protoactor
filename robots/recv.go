/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-04-01 21:12:30
 * Filename      : recv.go
 * Description   : 机器人
 * *******************************************************/
package robots

import (
	"fmt"
	"protoactor/algo"
	"protoactor/errorcode"
	"protoactor/messages"

	"github.com/golang/glog"
)

func (r *Robot) Recv(id uint32, msg interface{}) {
	switch id {
	case 1000:
		r.recvLogin(msg)
	case 4000:
		r.recvComein(msg)
	case 4010:
		r.recvCreate(msg)
	case 4029:
		r.recvDiscard(msg)
	case 4027:
		r.recvDraw(msg)
	case 4036:
		r.recvGameover(msg)
	case 4002:
		r.recvLeave(msg)
	case 1036:
		r.recvdata(msg)
	case 4025:
		r.recvGamestart(msg)
	case 4030:
		r.recvOperate(msg)
	case 4032:
		r.recvPengKong(msg)
	case 1002:
		r.recvRegist(msg)
	default:
	}
}

// 接收到服务器登录返回
func (c *Robot) recvRegist(msg interface{}) {
	stoc := msg.(*messages.SRegist)
	var errcode uint32 = stoc.GetError()
	fmt.Printf("regist -> %d\n", errcode)
	switch {
	case errcode == errorcode.AccountExist:
		c.SendLogin() //已经注册,登录
	case errcode == 0:
		var uid string = stoc.GetUserid()
		c.data.Userid = uid
		fmt.Printf("regist uid -> %s\n", uid)
		c.SendLogin() //注册成功,登录
	default:
		glog.Infof("regist err -> %d", errcode)
	}
}

// 接收到服务器登录返回
func (c *Robot) recvLogin(msg interface{}) {
	stoc := msg.(*messages.SLogin)
	var errcode uint32 = stoc.GetError()
	fmt.Printf("login -> %d\n", errcode)
	switch {
	case errcode == 0:
		c.data.Userid = stoc.GetUserid()
		glog.Infof("logined uid -> %s\n", c.data.Userid)
		glog.Infof("logined phone -> %s\n", c.data.Phone)
		c.SendUserData() // 获取玩家数据
	default:
		glog.Infof("login err -> %d", errcode)
	}
}

// 接收到玩家数据
func (c *Robot) recvdata(msg interface{}) {
	stoc := msg.(*messages.SUserData)
	var errcode uint32 = stoc.GetError()
	if errcode != 0 {
		glog.Infof("get data err -> %d", errcode)
	}
	data := stoc.GetData()
	// 设置数据
	c.data.Userid = data.GetUserid()     // 用户id
	c.data.Nickname = data.GetNickname() // 用户昵称
	c.data.Sex = data.GetSex()           // 用户性别,男1 女0
	c.data.Coin = data.GetCoin()         // 金币
	c.data.Diamond = data.GetDiamond()   // 钻石
	//查找房间-进入房间
	c.SendEntry()
}

// 离开房间
func (c *Robot) recvLeave(msg interface{}) {
	stoc := msg.(*messages.SLeave)
	var seat uint32 = stoc.GetSeat()
	if seat == c.seat {
		c.Close() //下线
	}
	if seat >= 1 && seat <= 4 && seat != c.seat {
		c.SendLeave() //离开
	}
}

// 创建房间
func (c *Robot) recvCreate(msg interface{}) {
	stoc := msg.(*messages.SCreateRoom)
	var errcode uint32 = stoc.GetError()
	switch {
	case errcode == 0:
		var code string = stoc.GetRdata().GetInvitecode()
		if code != "" {
			glog.Infof("create room code -> %s", code)
			c.code = code //设置邀请码
			c.SendEntry() //进入房间
		} else {
			glog.Errorf("create room code empty -> %s", code)
		}
	default:
		glog.Infof("create room err -> %d", errcode)
		c.Close() //进入出错,关闭
	}
}

// 进入房间
func (c *Robot) recvComein(msg interface{}) {
	stoc := msg.(*messages.SEnterRoom)
	var errcode uint32 = stoc.GetError()
	switch {
	case errcode == 0:
		userinfo := stoc.GetUserinfo()
		for _, v := range userinfo {
			if v.GetUserid() == c.data.Userid {
				c.seat = v.GetSeat()
				c.SendReady() //准备
				break
			}
		}
	default:
		glog.Infof("comein err -> %d", errcode)
		c.Close() //进入出错,关闭
	}
}

//// 解散
//func (c *Robot) recvVote(msg interface{}) {
//  stoc := msg.(*messages.SLaunchVote)
//	var seat uint32 = stoc.GetSeat()
//	glog.Infof("vote seat -> %d", seat)
//	c.SendVote()
//}

//开始游戏
func (c *Robot) recvGamestart(msg interface{}) {
	stoc := msg.(*messages.SGamestart)
	var dealer uint32 = stoc.GetDealer()
	var laipi uint32 = stoc.GetLaipi()
	var laizi uint32 = stoc.GetLaizi()
	var cards []uint32 = stoc.GetCards()
	var value uint32 = stoc.GetValue()
	glog.Infof("laipi -> %d, laizi -> %d", laipi, laizi)
	c.cards = cards
	if c.seat == dealer {
		if value == 0 {
			c.SendDiscard2()
		} else {
			c.operate(value, 0) //胡或暗杠
		}
	}
}

// 结束
func (c *Robot) recvGameover(msg interface{}) {
	stoc := msg.(*messages.SGameover)
	var round uint32 = stoc.GetRound()
	if round == 0 {
		c.Close() //结束下线
	} else {
		c.SendReady() //准备
	}
}

// 此协议向有碰和明杠的玩家主动发送
func (c *Robot) recvPengKong(msg interface{}) {
	stoc := msg.(*messages.SPengKong)
	card := stoc.GetCard()
	value := stoc.GetValue()
	c.operate(value, card) //胡碰杠吃
}

// 操作结果返回,同步手牌
func (c *Robot) recvOperate(msg interface{}) {
	stoc := msg.(*messages.SOperate)
	seat := stoc.GetSeat()               // 操作位置
	card := stoc.GetCard()               // 被碰或杠牌的牌值
	value := stoc.GetValue()             // 碰或值杠,统一掩码标示
	discontinue := stoc.GetDiscontinue() // 抢杠
	if seat != c.seat {                  // 不是自己操作
		return
	}
	// 根据不同操作同步cards
	switch {
	case value&algo.CHOW > 0:
		c1, c2, _ := algo.DecodeChow(card) //解码
		c.cards = RemoveN(uint32(c1), c.cards, 1)
		c.cards = RemoveN(uint32(c2), c.cards, 1)
		c.SendDiscard() //出牌
	case value&algo.PENG > 0:
		c.cards = RemoveN(uint32(card), c.cards, 2)
		c.SendDiscard() //出牌
	case value&algo.MING_KONG > 0:
		c.cards = RemoveN(uint32(card), c.cards, 3)
	case value&algo.AN_KONG > 0:
		c.cards = RemoveN(uint32(card), c.cards, 4)
	case value&algo.BU_KONG > 0:
		c.cards = RemoveN(uint32(card), c.cards, 1)
	case discontinue == algo.QIANG_GANG:
		c.SendHu()
	case discontinue > 0: //不同地址掩码值不一样
		c.SendHu()
	default:
		glog.Infof("SOperate err -> %d", value)
	}
}

// 玩家出牌广播
func (c *Robot) recvDiscard(msg interface{}) {
	stoc := msg.(*messages.SDiscard)
	var errcode uint32 = stoc.GetError()
	var seat uint32 = stoc.GetSeat()
	var card uint32 = stoc.GetCard()
	var value uint32 = stoc.GetValue()
	if errcode != 0 {
		glog.Infof("Discard err -> %d", errcode)
		return
	}
	if seat == c.seat { //自己出牌
		c.cards = RemoveN(uint32(card), c.cards, 1) //移除
		return
	}
	if value == 0 { //没有操作
		return
	}
	c.operate(value, card) //胡碰杠吃
}

// 抓牌
func (c *Robot) recvDraw(msg interface{}) {
	stoc := msg.(*messages.SDraw)
	var value uint32 = stoc.GetValue()
	var seat uint32 = stoc.GetSeat()
	var card uint32 = stoc.GetCard()
	//c.cards = stoc.GetCards()
	//PrintCards(c.cards)
	if seat == c.seat { //自己摸牌
		//glog.Infof("Discard seat -> %d", seat)
		c.cards = append(c.cards, uint32(card))
		if value == 0 {
			c.SendDiscard()
		} else {
			c.operate(value, card) //胡或暗杠
		}
	}
}

//操作一定要正确,服务器有些操作没做验证
func (c *Robot) operate(value, card uint32) {
	switch {
	case value&algo.HU > 0:
		c.SendHu()
	case value&algo.MING_KONG > 0:
		c.SendOperate(card, algo.MING_KONG)
	case value&algo.BU_KONG > 0:
		c.SendOperate(card, algo.BU_KONG)
	case value&algo.AN_KONG > 0:
		var kongcard uint32 = findKong(c.cards) //找一个暗扛
		c.SendOperate(uint32(kongcard), algo.AN_KONG)
	case value&algo.PENG > 0:
		c.SendOperate(card, algo.PENG)
	case value&algo.CHOW > 0:
		var chowcard uint32 = findChow(uint32(card), c.cards)
		c.SendOperate(chowcard, algo.CHOW)
	}
}

//判断是否有4个相同的牌
func findKong(cards []uint32) uint32 {
	m := make(map[uint32]int, len(cards))
	for _, v := range cards {
		if i, ok := m[v]; ok {
			if i == 3 {
				return v
			}
			m[v] = i + 1
		} else {
			m[v] = 1
		}
	}
	return 0
}

//选择一个吃牌,TODO:选最优吃牌组合
func findChow(card uint32, cards []uint32) uint32 {
	var t uint32 = card >> 4
	var chow []uint32
	if uint32(t) >= algo.FENG {
		chow = findChow1(t, card, cards)
	} else {
		chow = findChow2(card, cards)
	}
	chow = RemoveN(card, chow, 1)
	chow = RemoveN(0, chow, len(chow))
	if len(chow) < 2 {
		panic(fmt.Sprintf("chow err -> %+x, %x, %+x", chow, card, cards))
	}
	return (uint32(chow[0]) << 8) | uint32(chow[1])
}

//风牌,字牌两张不同即可
func findChow1(t, card uint32, cs []uint32) []uint32 {
	var chow []uint32 = []uint32{}
	var c uint32
	for _, v := range cs {
		if v != card && v != c && v>>4 == t {
			chow = append(chow, v)
			c = v
		}
	}
	return chow
}

//数牌三张
func findChow2(card uint32, cs []uint32) []uint32 {
	var cards []uint32 = make([]uint32, 5)
	cards[2] = card
	for _, v := range cs {
		var s int = int(card) - int(v) + 2
		if s >= 0 && s <= 4 {
			cards[s] = v
		}
	}
	var chow []uint32 = []uint32{}
	count := 0
	for _, v := range cards {
		if count >= 3 {
			break
		}
		if v > 0 {
			chow = append(chow, v)
			count += 1
		} else {
			chow = []uint32{}
			count = 0
		}
	}
	return chow
}

// 移除n个牌
func RemoveN(c uint32, cs []uint32, n int) []uint32 {
	for n > 0 {
		for i, v := range cs {
			if c == v {
				cs = append(cs[:i], cs[i+1:]...)
				break
			}
		}
		n--
	}
	return cs
}

func PrintCards(cards []uint32) (str string) {
	for _, card := range cards {
		str += fmt.Sprintf("%x,", card)
	}
	glog.Infof("cards -> %s", str)
	return
}

// 选择最优出牌,TODO:优化
func SearchCard(cards []uint32) uint32 {
	m := make(map[uint32][]int)
	for _, card := range cards {
		var k uint32 = uint32(card >> 4)
		var v uint32 = uint32(card & 0x0f)
		var n []int
		if n1, ok := m[k]; ok {
			n = n1
		} else {
			switch {
			case k == algo.FENG:
				n = []int{3: 0}
			case k == algo.ZI:
				n = []int{2: 0}
			default:
				n = []int{8: 0}
			}
		}
		n[v-1] += 1
		m[k] = n
	}
	//cards := []uint32{0x07, 0x08, 0x09, 0x18, 0x18, 0x41,0x43,0x44}
	//map[0:[0 0 0 0 0 0 1 1 1] 1:[0 0 0 0 0 0 0 2 0] 4:[1 0 1 1]]
	for k, v := range m { //优先了顺子
		if k == algo.FENG {
			m[k] = shun_feng(v)
		} else {
			m[k] = shun(v)
		}
	}
	//TODO:优化
	for k, v := range m {
		var v_l int = len(v)
		for i := 0; i < v_l; i++ {
			var n int = 0
			if i == 0 {
				n = v[i] + v[i+1] + v[i+2]
			} else if i == v_l-1 {
				n = v[i-2] + v[i-1] + v[i]
			} else {
				n = v[i-1] + v[i] + v[i+1]
			}
			if v[i] == 1 && n == 1 {
				return uint32(int(k<<4) | (i + 1))
			}
		}
	}
	//
	for k, v := range m {
		for i, j := range v {
			if j == 1 {
				return uint32(int(k<<4) | (i + 1))
			}
		}
	}
	return cards[0]
}

func shun_feng(v []int) []int {
	if v[0] >= 1 && v[1] >= 1 && v[2] >= 1 {
		v[0] -= 1
		v[1] -= 1
		v[2] -= 1
	} else if v[0] >= 1 && v[2] >= 1 && v[3] >= 1 {
		v[0] -= 1
		v[2] -= 1
		v[3] -= 1
	} else if v[1] >= 1 && v[2] >= 1 && v[3] >= 1 {
		v[1] -= 1
		v[2] -= 1
		v[3] -= 1
	}
	return v
}

func shun(v []int) []int {
	var v_l int = len(v)
	var i, j int = 0, 0
	for i = 0; i < v_l; i++ {
		if v[i] >= 1 {
			j++
		} else {
			j = 0
		}
		if j == 3 {
			v[i] -= 1
			v[i-1] -= 1
			v[i-2] -= 1
			j = 0
			return shun(v)
		}
	}
	return v
}

// vim: set fdm=marker:
