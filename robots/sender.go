/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-01-16 10:00
 * Filename      : sender.go
 * Description   : 机器人
 * *******************************************************/
package robots

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"protoactor/algo"
	"protoactor/messages"
	"utils"

	"github.com/golang/glog"
)

// 发送注册请求
func (c *Robot) SendRegist() {
	//随机初始化
	rl := GetRobotList()
	var r_l int32 = int32(len(rl))
	var num int32 = utils.RandInt32N(r_l)
	r_d := rl[num]
	c.data.Nickname = r_d.Nickname
	c.data.Sex = 0
	utils.SleepRand(3) //停留
	//开始注册
	h := md5.New()
	h.Write([]byte("piaohua")) // 需要加密的字符串为 123456
	pwd := hex.EncodeToString(h.Sum(nil))
	ctos := &messages.CRegist{
		Phone:    c.data.Phone,
		Nickname: c.data.Nickname,
		Password: pwd,
	}
	fmt.Printf("phone -> %s\n", c.data.Phone)
	fmt.Printf("nickname -> %s\n", c.data.Nickname)
	c.Sender(ctos)
}

// 发送登录请求
func (c *Robot) SendLogin() {
	h := md5.New()
	h.Write([]byte("piaohua")) // 需要加密的字符串为 123456
	pwd := hex.EncodeToString(h.Sum(nil))
	ctos := &messages.CLogin{
		Phone:    c.data.Phone,
		Password: pwd,
	}
	c.Sender(ctos)
}

// 获取玩家数据
func (c *Robot) SendUserData() {
	ctos := &messages.CUserData{
		Userid: c.data.Userid,
	}
	c.Sender(ctos)
}

// 玩家创建房间
func (c *Robot) SendCreate() {
	//var a1 []uint32 = []uint32{4, 8, 16}
	//var a2 []uint32 = []uint32{1, 5, 9}
	//var i int32 = utils.RandInt32N(3) //随机
	ctos := &messages.CCreateRoom{
		Round:   4,
		Rtype:   1,
		Ante:    1,
		Pao:     1,
		Payment: 0,
		Rname:   "ddd",
	}
	glog.Infof("create room phone -> %s", c.data.Phone)
	c.Sender(ctos)
}

// 玩家进入房间
func (c *Robot) SendEntry() {
	if c.code == "create" { //表示创建房间
		c.SendCreate() //创建一个房间
	} else { //表示进入房间
		ctos := &messages.CEnterRoom{
			Invitecode: c.code,
		}
		c.Sender(ctos)
	}
}

// 玩家准备
func (c *Robot) SendReady() {
	ctos := &messages.CReady{Ready: true}
	c.Sender(ctos)
}

// 离开
func (c *Robot) SendLeave() {
	ctos := &messages.CLeave{}
	c.Sender(ctos)
}

//// 解散
//func (c *Robot) SendVote() {
//	ctos := &messages.CVote{
//		Vote: 0,
//	}
//	c.Sender(ctos)
//}

// 庄家出牌
func (c *Robot) SendDiscard2() {
	utils.Sleep(6) //展示动画时间
	c.SendDiscard()
}

// 自己出牌
func (c *Robot) SendDiscard() {
	card := uint32(SearchCard(c.cards))
	ctos := &messages.CDiscard{Card: card}
	utils.SleepRand(3)
	c.Sender(ctos)
}

// 玩家碰杠操作
func (c *Robot) SendOperate(card, value uint32) {
	ctos := &messages.COperate{
		Card:  card,
		Value: value,
	}
	utils.SleepRand(3)
	c.Sender(ctos) // send data
}

// 胡牌
func (c *Robot) SendHu() {
	c.SendOperate(0, algo.HU)
}
