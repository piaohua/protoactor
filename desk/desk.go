/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-03-28 21:05:38
 * Filename      : desk.go
 * Description   : 玩牌逻辑
 * *******************************************************/
package desk

import (
	"fmt"
	"math/rand"
	"protoactor/algo"
	"protoactor/entity"
	"time"
	"utils"

	"github.com/golang/glog"
)

//// external function

//新建一张牌桌
func NewDesk(data *entity.DeskData) *entity.Desk {
	desk := &entity.Desk{
		DeskData: data,
		//------
		Votes:       make(map[uint32]uint32),
		Record:      utils.NewArray(true),
		Trusteeship: make(map[uint32]bool),
		Ready:       make(map[uint32]bool),
	}
	return desk
}

//玩家准备
func Ready(userid string, ready bool, t *entity.Desk) int {
	if t.Vote != 0 { //投票中不能准备
		return 1
	}
	seat := getRoles(userid, t)
	if seat == 0 {
		return 2
	}
	t.Ready[seat] = ready //设置状态
	msg := res_ready(seat, ready)
	Broadcast(msg, t) //广播消息
	diceing(t)        //主动打骰
	return 0
}

//房间消息广播
func Broadcast(msg interface{}, t *entity.Desk) {
	for _, u := range t.Seats {
		u.Pid.Tell(msg)
	}
}

//玩家打骰子切牌,发牌
func diceing(t *entity.Desk) bool {
	if isDiceing(t) { //是否打骰
		gameStart(t) //开始牌局
		return true
	} else {
		return false
	}
}

//是否可以打骰
func isDiceing(t *entity.Desk) bool {
	if t.State { //已经开始
		return false
	}
	if t.Dealer != 0 { //已经打庄
		return false
	}
	if uint32(len(t.Seats)) != t.Count { //人数不够
		return false
	}
	if !allReady(t) { //没有准备
		return false
	}
	return true
}

//是否全部准备状态
func allReady(t *entity.Desk) bool {
	if uint32(len(t.Ready)) != t.Count {
		return false
	}
	for _, ok := range t.Ready {
		if !ok {
			return false
		}
	}
	return true
}

//开始游戏
func gameStart(t *entity.Desk) {
	gameStartInit(t) //初始化
	glog.Infof("gameStart -> %d, seat -> %d", t.Rid, t.Seat)
	//打庄消息通知+抽水
	//var coin int32 = -100 //TODO:抽水
	//var cost int32 = -1 * int32(t.data.Cost)
	//for _, u := range t.Seats {
	//	//u.Pid.Tell()
	//	//resource.ChangeRes(p, resource.COIN, coin, 1) //抽水
	//	//开始游戏扣除创建房间钻石
	//	if t.Round == 0 {
	//		//TODO
	//	}
	//}
	//打骰(两个骰子)
	dice1 := uint32(utils.RandInt32N(5) + 11)
	dice2 := uint32(utils.RandInt32N(5) + 1)
	t.Dice = dice1 + dice2 //1-6的骰子数
	shuffle(t)             //洗牌
	dealer_(t)             //打庄
	deal(t)                //发牌
	//等待玩家操作
}

//初始化
func gameStartInit(t *entity.Desk) {
	t.State = true //设置房间状态
	t.Timer = -6   //重置计时器,6秒给前端展示动画
	if t.CloseCh == nil {
		t.CloseCh = make(chan bool, 1)
		go ticker(t) //计时器goroutine
	}
	//------
	operateInit(t)
	//------
	t.OutCards = make(map[uint32][]uint32)  //海底牌
	t.PongCards = make(map[uint32][]uint32) //碰牌
	t.KongCards = make(map[uint32][]uint32) //杠牌
	t.ChowCards = make(map[uint32][]uint32) //吃牌(8bit-8-8)
	t.HandCards = make(map[uint32][]uint32) //手牌
	t.OutLaizi = make(map[uint32]int)       //打出赖子数量
}

//初始化操作值
func operateInit(t *entity.Desk) {
	t.Hu = make(map[uint32]uint32)
	t.Pongkong = make(map[uint32]uint32)
	t.Chow = make(map[uint32]uint32)
}

//计时器
func ticker(t *entity.Desk) {
	tick := time.Tick(time.Second)
	glog.Infof("ticker -> %d", t.Rid)
	for {
		select {
		case <-tick:
			//超时判断
			if t.State {
				if t.Timer == OT || t.Timer == DT {
					timerL(t) //逻辑处理
				} else {
					t.Timer++
				}
			}
		case <-t.CloseCh:
			glog.Infof("close desk -> %d", t.Rid)
			return
		}
	}
}

//超时处理
func timerL(t *entity.Desk) {
	//操作(胡,碰杠,吃)超时处理
	if t.Timer == OT && t.Discard != 0 {
		t.Timer = 0
		t.Pid.Tell(res_timerTurn(t.Seat))
		//TurnL(t) //操作(自动胡牌,操作超时...)
	} else if t.Timer == DT && t.Draw != 0 {
		//出牌超时处理
		t.Timer = 0
		p := getPlayer(t.Seat, t)
		t.Pid.Tell(res_timerDiscard(p.Userid, t.Draw))
		//DiscardL(t.Seat, t.Draw, true, t) //超时打出摸到的牌
	} else {
		t.Timer++
	}
}

//洗牌
func shuffle(t *entity.Desk) {
	rand.Seed(time.Now().UnixNano())
	d := make([]uint32, algo.TOTAL, algo.TOTAL)
	copy(d, algo.CARDS)
	//测试暂时去掉洗牌
	for i := range d {
		j := rand.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
	t.Cards = d
}

//打庄处理
func dealer_(t *entity.Desk) {
	//选择庄家
	if t.Lian == 0 {
		t.Dealer = uint32(utils.RandInt32N(int32(t.Count)) + 1)
	} else {
		t.Dealer = t.Lian
	}
}

//发牌
func deal(t *entity.Desk) {
	for s, _ := range t.Seats {
		var hand int = int(algo.HAND)
		if s == t.Dealer { //判断庄家发14张牌
			hand += 1
		}
		cards := make([]uint32, hand, hand)
		tmp := t.Cards[:hand]
		copy(cards, tmp)
		t.HandCards[s] = cards
		t.Cards = t.Cards[hand:]
	}
	//翻赖子
	lai(t)
	//第一个操作为庄家
	t.Seat = t.Dealer
	//发牌协议消息
	for s, u := range t.Seats {
		var cards []uint32 = getHandCards(s, t)
		if t.Dealer == s {
			//庄家提示处理
			var v uint32 = algo.DrawDetectTao(t.Laizi, uint32(0), cards, []uint32{})
			v |= heType2(v, t.Seat, 0, cards, t)
			if v > 0 {
				t.Hu[s] = v //设置操作状态值
			}
			t.Draw = cards[len(cards)-1] //庄家最后一张默认为摸牌
			//庄家消息
			msg := res_gamestart(t.Dice, t.Dealer, v, t.Laipi, t.Laizi, cards)
			u.Pid.Tell(msg)
		} else {
			//闲家消息
			msg := res_gamestart(t.Dice, t.Dealer, 0, t.Laipi, t.Laizi, cards)
			u.Pid.Tell(msg)
		}
	}
}

//赖子
func lai(t *entity.Desk) {
	t.Laipi = t.Cards[0]
	if (t.Laipi & 0x0f) == 9 {
		t.Laizi = (t.Laipi & 0xf0) | 0x01
	} else {
		t.Laizi = t.Laipi + 1
	}
	t.Cards = t.Cards[1:]
}

//进入房间
func EnterMsg(userid string, t *entity.Desk) (bool, interface{}) {
	seat := getRoles(userid, t)
	if seat == 0 {
		return false, nil
	}
	msg := res_enter(seat, t)
	msg1 := res_camein(seat, t)
	Broadcast_(seat, msg1, t)
	return true, msg
}

//重复进入或重连
func Enter(userid string, t *entity.Desk) (bool, interface{}) {
	seat := getRoles(userid, t)
	if seat == 0 {
		return false, nil
	}
	msg := res_reEnter(seat, t)
	//glog.Infof("Enter userid: %s, t: %#v", userid, msg)
	return true, msg
}

//玩家离开
func Leave(userid string, t *entity.Desk) int {
	if t.State { //游戏中不能离开
		return 1
	}
	seat := getRoles(userid, t)
	if seat == 0 {
		return 2
	}
	//广播消息
	msg := res_leave(seat)
	Broadcast(msg, t)
	//清除数据
	delete(t.Ready, seat)
	delete(t.Seats, seat)
	delete(t.Roles, userid)
	return 0
}

/*
//踢除玩家
func (t *Desk) Kick(cid string, seat uint32) bool {
	if t.State { //游戏中不能离开
		return false
	}
	if cid != t.Cid {
		return false
	}
	if p, ok := t.Seats[seat]; ok {
		//p.Pid.Tell(leave room)
		//p.ClearRoom() //清除玩家房间数据
	}
	//广播消息
	msg := res_leave(seat)
	Broadcast(msg, t)
	//清除数据
	delete(t.Ready, seat)
	delete(t.Seats, seat)
	delete(t.Roles, userid)
	return true
}

//托管处理,Kind=1:托管;Kind=0:取消托管
func trust_(seat, kind uint32, t *entity.Desk) {
	var tru bool = false
	if kind == 1 { //托管
		tru = true
	}
	if getTrust(seat, t) == tru { //相同设置不处理
		return
	}
	t.Trusteeship[seat] = tru //设置状态
	msg := res_trust(seat, kind)
	Broadcast(msg, t) //广播消息
}
*/

//// internal function

//模牌,kong==false普通摸牌,kong==true扛后摸牌
func drawcard(t *entity.Desk) {
	if len(t.Cards) == 0 { //TODO 出牌放入outCards再结束
		he(0, t) //结束牌局
		return
	}
	var value uint32 = 0
	if t.Laidraw { //出赖后摸牌,不做位置切换,和杠后摸牌一样
		var outs []uint32 = getOutCards(t.Seat, t)
		outs = append(outs, t.Discard)
		t.OutCards[t.Seat] = outs
		t.OutLaizi[t.Seat] += 1
		//不做位置切换,继续摸牌
	} else if t.Kong { //杠后摸牌
		value = 1
	} else { //普通摸牌, t.Discrad != 0
		var outs []uint32 = getOutCards(t.Seat, t)
		outs = append(outs, t.Discard)
		t.OutCards[t.Seat] = outs
		//位置切换
		t.Seat = algo.NextSeat(t.Seat)
	}
	var card uint32 = t.Cards[0]
	t.Draw = card //设置摸牌状态
	t.Discard = 0 //清除打牌状态
	t.Operate = 0 //清除操作状态
	t.Huing = 0   //清除胡牌状态
	t.Cards = t.Cards[1:]
	var cards []uint32 = in(t.Seat, card, t)
	var ps []uint32 = getPongCards(t.Seat, t)
	var v uint32 = algo.DrawDetectTao(t.Laizi, card, cards, ps)
	v |= heType2(v, t.Seat, 0, cards, t)
	//杠上开花
	if t.Kong && v&algo.HU > 0 {
		v = v | algo.HU_KONG_FLOWER
	}
	if v > 0 { //摸牌全部记录为胡(只自己操作)
		t.Hu[t.Seat] = v
	}
	//记录 TODO:暂时没添加
	//t.Record(algo.DRAW, card, t.seat)
	//摸牌协议消息通知
	for s, o := range t.Seats {
		if s == t.Seat {
			//摸牌玩家消息
			msg := res_draw(t.Seat, value, v, card)
			//o.Send(msg)
			o.Pid.Tell(msg)
		} else {
			//其他玩家消息
			msg := res_draw(t.Seat, value, 0, 0)
			//o.Send(msg)
			o.Pid.Tell(msg)
		}
	}
	//托管处理
	if getTrust(t.Seat, t) {
		t.Timer = DT - 2 //出牌时间
	} else {
		t.Timer = 0 //重置计时
	}
}

//出牌,加锁,托管自摸自动胡牌
func DiscardL(userid string, card uint32, ok bool, t *entity.Desk) int {
	if !t.State {
		glog.Infof("DiscardL err -> %s", userid)
		return 2
	}
	seat := getRoles(userid, t)
	if seat == 0 {
		return 1
	}
	if seat < 1 || seat > t.Count {
		glog.Infof("Discard seat err -> %d", seat)
		return 1
	}
	if seat != t.Seat || t.Draw == 0 {
		p := getPlayer(seat, t)
		msg := res_discard2()
		//p.Send(msg)
		p.Pid.Tell(msg)
		glog.Infof("Discard seat err -> %d", seat)
		return 3
	}
	if ok { //超时操作
		//trust_(seat, 1, t) //设置玩家超时托管
		if _, ok := t.Hu[seat]; ok {
			he(seat, t) //托管自摸自动胡牌
			return 0
		}
	} else {
		//trust_(seat, 0, t) //设置玩家超时托管
	}
	discard_(card, t)
	return 0
}

//出牌,没加锁
func discard_lai(card uint32, t *entity.Desk) bool {
	if card != t.Laizi {
		if t.Laidraw { //出赖标识
			t.Laidraw = false //出赖标识清除
		}
		return false
	}
	t.Laidraw = true
	return true
}

//出牌,没加锁
func discard_(card uint32, t *entity.Desk) {
	//记录 TODO:暂时没添加
	//t.Record(algo.DISCARD, card, t.seat)
	operateInit(t)   //清除操作记录
	t.Discard = card //设置打牌状态
	t.Draw = 0       //清除摸牌状态
	t.Operate = 0    //清除操作状态
	t.Huing = 0      //清除胡牌状态
	var out_lai bool = discard_lai(card, t)
	//检测(胡,碰杠,吃)
	for s, _ := range t.Seats {
		if s == t.Seat { //出牌人跳过
			out(t.Seat, card, t) //移除牌
			continue
		}
		if out_lai { //出赖不做任何操作
			continue
		}
		var cards []uint32 = getHandCards(s, t)
		if dapai(t) { //是否大牌玩法
			//胡,杠碰,吃检测
			v_h := algo.DiscardHuTao(t.Laizi, card, cards) //胡
			v_h |= heType2(v_h, s, card, cards, t)
			v_h = algo.PaoHuPing(v_h) //平胡不能炮胡
			if v_h > 0 {
				if t.Kong { //杠后出牌热炮
					v_h |= algo.REPAO
				}
				t.Hu[s] = v_h
			}
		}
		v_p := algo.DiscardPong(card, cards) //碰杠
		if v_p > 0 {
			t.Pongkong[s] = v_p
		}
	}
	if t.Kong { //杠操作出牌标识
		t.Kong = false //杠后出牌清除
	}
	order(t) //按规则优先级设置
}

//按规则优先级设置
func order(t *entity.Desk) {
	//TODO:优化order + turn
	var l_h, l_p, l_c int = getLength(t)
	if l_h == 0 && l_p == 0 && l_c == 0 { //无操作
		//出牌协议消息通知
		msg := res_discard(t.Seat, t.Discard)
		Broadcast(msg, t) //消息广播
		drawcard(t)       //摸牌
		return
	}
	//多人胡,碰杠一人,吃一人
	if l_p == 1 {
		for k, v := range t.Pongkong {
			if vc, ok := t.Chow[k]; ok { //碰杠吃为同一人
				t.Pongkong[k] = v | vc
				t.Chow = make(map[uint32]uint32) //清除
				break
			}
		}
	} else if l_p == 0 && l_c != 0 { //吃,胡一起提示,TODO:优化
		t.Pongkong = t.Chow
		t.Chow = make(map[uint32]uint32) //清除
	}
	if l_h == 1 { //一人胡
		for k, v := range t.Hu {
			if vc, ok := t.Pongkong[k]; ok { //胡碰杠为同一人
				t.Hu[k] = v | vc
				t.Pongkong = make(map[uint32]uint32) //清除
				break
			}
		}
	}
	//多人胡按顺序操作
	turn(false, t)
}

//超时操作,加锁
func TurnL(t *entity.Desk) {
	turn(true, t)
}

//超时操作,没加锁 ok表示是否超时处理
func turn(ok bool, t *entity.Desk) {
	var l_h, l_p, l_c int = getLength(t)
	if l_h == 0 && l_p == 0 && l_c == 0 { //无操作
		drawcard(t) //摸牌
		return
	}
	if ok { //超时
		if l_h != 0 {
			for k, _ := range t.Hu {
				he(k, t) //超时自动胡牌
				return
			}
		} else if l_p != 0 {
			t.Pongkong = make(map[uint32]uint32)
		} else if l_c != 0 {
			t.Chow = make(map[uint32]uint32)
		}
		//操作超时不托管,只有打牌超时才托管
		//消除提示操作消息,下个操作消息同时进行
		turn(false, t) //发送下一个操作消息
		return
	}
	//操作按胡-碰杠-吃顺序,相同操作(胡),TODO:优化
	//t.hu -> t.pongkong -> t.chow
	var m map[uint32]uint32 = make(map[uint32]uint32)
	if l_h != 0 {
		m = t.Hu
	} else if l_p != 0 {
		m = t.Pongkong
	} else if l_c != 0 {
		m = t.Chow
	} else {
		panic(fmt.Sprintf("turn error:%d", t.Seat))
	}
	var card uint32 = operateCard(t) //所操作的牌
	//广播操作消息
	for k, v := range m {
		p := getPlayer(k, t)
		if t.Operate == 0 { //第一次提示操作
			t.Operate += 1 //操作状态变化
			msg1 := res_discard3(t.Seat, v, card)
			p.Pid.Tell(msg1)
			//p.Send(msg1)
			msg2 := res_discard(t.Seat, card)
			Broadcast_(k, msg2, t) //消息广播
		} else { //二次提示操作
			t.Operate += 1 //操作状态变化
			msg3 := res_pengkong(t.Seat, v, card)
			//p.Send(msg3)
			p.Pid.Tell(msg3)
		}
		//托管处理
		if getTrust(k, t) {
			t.Timer = OT - 2 //出牌时间
		} else {
			t.Timer = 0 //重置计时
		}
	}
}

//胡碰杠吃操作 TODO:优化
func OperateL(card, value uint32, userid string, t *entity.Desk) int {
	if !t.State {
		glog.Infof("operate err -> %s", userid)
		return 2
	}
	seat := getRoles(userid, t)
	if seat == 0 {
		return 1
	}
	if seat < 1 || seat > t.Count {
		glog.Infof("operate seat err -> %d", seat)
		return 1
	}
	var l_h, l_p, l_c int = getLength(t)
	if l_h == 0 && l_p == 0 && l_c == 0 { //无操作
		glog.Infof("operate err -> %d", seat)
		return 3
	}
	//trust_(seat, 0, t) //设置玩家超时托管
	if l_h > 1 { //多人胡,有胡
		operateH(card, value, seat, t) //TODO:区别处理??
	} else if l_h == 1 { //一人胡,有胡碰杠吃
		operateH(card, value, seat, t)
	} else if l_p == 1 { //没有胡,有碰杠吃
		operateP(card, value, seat, t)
	} else if l_c == 1 { //没有碰杠,有吃
		operateC(card, value, seat, t)
	} else { //错误
		glog.Infof("unexpected operate err -> %d", seat)
	}
	return 0
}

//胡碰杠吃操作,TODO:优化
func operateH(card, value, seat uint32, t *entity.Desk) {
	if v, ok := t.Hu[seat]; ok {
		if value == 0 {
			delete(t.Hu, seat)
			cancelOperate(seat, t.Seat, value, card, t)
			if len(t.Hu) == t.Huing && t.Huing != 0 {
				for k, _ := range t.Hu {
					he(k, t) //选择位置,TODO:超时时huing处理
					return
				}
			} else if seat != t.Seat { //如果暗杠时取消不应该摸牌,应该打牌
				turn(false, t) //取消时进入下一个操作
			}
		} else if value&algo.HU > 0 {
			t.Huing++                 //选择胡牌人数
			if len(t.Hu) == t.Huing { //多人胡
				he(seat, t)
			}
		} else {
			operateA(card, value, seat, v, t)
		}
	} else {
		//不存在操作
		glog.Errorf("not exist operate err -> %d", seat)
	}
}

//胡碰杠吃操作
func operateP(card, value, seat uint32, t *entity.Desk) {
	if v, ok := t.Pongkong[seat]; ok {
		if value == 0 {
			delete(t.Pongkong, seat)
			cancelOperate(seat, t.Seat, value, card, t)
			turn(false, t) //取消时进入下一个操作
		} else {
			operateA(card, value, seat, v, t)
		}
	} else {
		//不存在操作
		glog.Errorf("not exist operate err -> %d", seat)
	}
}

//胡碰杠吃操作
func operateC(card, value, seat uint32, t *entity.Desk) {
	if v, ok := t.Chow[seat]; ok {
		if value == 0 {
			delete(t.Chow, seat) //先清除
			cancelOperate(seat, t.Seat, value, card, t)
			turn(false, t) //取消时进入下一个操作
		} else {
			operateA(card, value, seat, v, t)
		}
	} else {
		//不存在操作
		glog.Errorf("not exist operate err -> %d", seat)
	}
}

//胡碰杠吃操作,TODO:碰杠验证处理
func operateA(card, value, seat, v uint32, t *entity.Desk) {
	if value&algo.CHOW > 0 && v&algo.CHOW > 0 {
		chow_(card, value, seat, t)
	} else if value&algo.PENG > 0 && v&algo.PENG > 0 {
		pong_(card, value, seat, t)
	} else if value&algo.MING_KONG > 0 && v&algo.MING_KONG > 0 {
		kong_(card, value, seat, v, algo.MING_KONG, t)
	} else if value&algo.BU_KONG > 0 && v&algo.BU_KONG > 0 {
		kong_(card, value, seat, v, algo.BU_KONG, t)
	} else if value&algo.AN_KONG > 0 && v&algo.AN_KONG > 0 {
		kong_(card, value, seat, v, algo.AN_KONG, t)
	} else {
		//不存在操作
		glog.Errorf("not exist operate err -> %d", seat)
	}
}

//胡牌,(多人胡牌时同时胡),t.seat=出牌(放冲)位置
func he(seat uint32, t *entity.Desk) {
	for k, v := range t.Seats {
		glog.Infof("he seat %d -> uid %s", k, v.Userid)
	}
	glog.Infof("id -> %d, he -> %+x", t.Rid, t.Hu)
	glog.Infof("seat -> %d, t.seat -> %d", seat, t.Seat)
	glog.Infof("discard -> %x, draw -> %x", t.Discard, t.Draw)
	glog.Infof("handCards -> %+x", t.HandCards)
	glog.Infof("pongCards -> %+x", t.PongCards)
	glog.Infof("kongCards -> %+x", t.KongCards)
	glog.Infof("chowCards -> %+x", t.ChowCards)
	glog.Infof("outCards -> %+x", t.OutCards)
	glog.Infof("outLaizi -> %+x", t.OutLaizi)
	var card uint32 = operateCard(t) //所胡的牌
	tianHe(t)                        //天地胡
	glog.Infof("id -> %d, he -> %+x", t.Rid, t.Hu)
	glog.Infof("seat -> %d, t.seat -> %d", seat, t.Seat)
	glog.Infof("seat -> %d, t.dealer -> %d", seat, t.Dealer)
	qiangKongHe(seat, card, t)   //抢杠胡处理
	var l_h bool = len(t.Hu) > 0 //是否胡牌
	//算番
	huFan, mingKong, beMingKong, anKong, buKong, total := gameOver(t)
	coin := make(map[uint32]int32)
	score := make(map[uint32]int32)
	//TODO:金币,积分结算
	for i, v := range total {
		p := getPlayer(i, t)
		//v *= int32(t.Ante)           //番*底分
		t.Score[p.Userid] += v       //总分
		coin[i] = v                  //当局分
		score[i] = t.Score[p.Userid] //总分
	}
	updateUser(l_h, t) //更新玩家数据
	glog.Infof("score -> %+v", t.Score)
	glog.Infof("huFan:%+v, mingKong:%+v, beMingKong:%+v, anKong:%+v, buKong:%+v, total:%+v",
		huFan, mingKong, beMingKong, anKong, buKong, total)
	t.Round++ //局数
	round, expire := getRound(t)
	//结算消息广播
	msg := res_over(t.Rid, seat, t.Seat, card, round, expire, t.Count, t.HandCards,
		t.Hu, huFan, mingKong, beMingKong, anKong, buKong, total, coin, score)
	Broadcast(msg, t)
	//t.privaterecord(coin)         //日志记录 TODO:goroutine
	lianDealer(l_h, seat, t)           //连庄
	overSet(t)                         //重置状态
	closeDesk(round, expire, false, t) //结束牌局
}

//更新玩家数据
func updateUser(l_h bool, t *entity.Desk) {
	for i, u := range t.Seats {
		if !l_h {
			msg := res_update(0, 0, 1, 0)
			u.Pid.Tell(msg) //玩家数据更新
			continue
		}
		if _, ok := t.Hu[i]; ok {
			msg := res_update(1, 0, 0, 0)
			u.Pid.Tell(msg) //玩家数据更新
		} else {
			msg := res_update(0, 1, 0, 0)
			u.Pid.Tell(msg) //玩家数据更新
		}
	}
}

//1.胡的牌(自摸,放炮),2.抢杠胡时操作的牌(补杠人的摸牌),TODO:优化
func operateCard(t *entity.Desk) uint32 {
	if t.Discard != 0 {
		return t.Discard
	} else {
		return t.Draw
	}
}

//抢杠胡处理,TODO:多人抢杠胡
//TODO:十三幺可以抢所有杠,杠番值暂时一样,可以暂时不处理
func qiangKongHe(seat, card uint32, t *entity.Desk) {
	if v, ok := t.Hu[seat]; ok {
		if v&algo.QIANG_GANG > 0 {
			msg := res_operate2(seat, t.Seat, algo.QIANG_GANG,
				algo.QIANG_GANG, card)
			Broadcast(msg, t)
			//去掉被抢杠玩家(t.seat)的补杠,TODO:要不要还原之前的碰?
			kongs := getKongCards(t.Seat, t)
			var cs []uint32
			for i, v2 := range kongs { //杠
				_, c, mask := algo.DecodeKong(v2) //解码
				if c == card && mask == algo.BU_KONG {
					cs = append(kongs[:i], kongs[i+1:]...)
					t.KongCards[t.Seat] = cs
					break
				}
			}
		}
	}
}

//是否大牌玩法
func dapai(t *entity.Desk) bool {
	return t.Pao == 1 //是否大牌玩法
}

//天地胡
func tianHe(t *entity.Desk) {
	if dapai(t) { //是否大牌玩法
		return
	}
	var l_s uint32 = uint32(len(t.Cards))
	var l_h uint32 = algo.TOTAL - (algo.HAND*t.Count + 1)
	if l_s != l_h {
		return
	}
	//TIAN_HU,DI_HU
	var l_o int = len(t.OutCards)
	var l_p int = len(t.PongCards)
	var l_k int = len(t.KongCards)
	var l_c int = len(t.ChowCards)
	if l_o == 0 && l_p == 0 &&
		l_k == 0 && l_c == 0 {
		if t.Seat == t.Dealer && t.Discard == 0 {
			for k, v := range t.Hu {
				t.Hu[k] = v | algo.TIAN_HU
			}
		} else {
			for k, v := range t.Hu {
				t.Hu[k] = v | algo.DI_HU
			}
		}
	}
}

//胡牌牌型检测
func heType2(val, seat, card uint32, hands []uint32, t *entity.Desk) uint32 {
	if !dapai(t) { //是否大牌玩法
		return 0
	}
	if val == 0 {
		return 0
	}
	if val&algo.HU == 0 {
		return 0
	}
	var cs []uint32 = []uint32{}
	if card != 0 {
		cs = append(cs, card)
	}
	cs = append(cs, hands...)
	kongs := getKongCards(seat, t)
	for _, v2 := range kongs { //杠
		_, c, _ := algo.DecodeKong(v2) //解码
		cs = append(cs, c, c, c, c)
	}
	pongs := getPongCards(seat, t)
	for _, v3 := range pongs { //碰
		_, c := algo.DecodePeng(v3) //解码
		cs = append(cs, c, c, c)
	}
	var v uint32 = algo.HuTypeDetect2(cs)
	if v == 0 {
		return 0
	}
	if card == 0 { //自摸
		return v | algo.ZIMO
	} else { //放炮
		return v | algo.PAOHU
	}
}

//连庄
func lianDealer(l_h bool, seat uint32, t *entity.Desk) {
	if l_h { //胡牌
		if len(t.Hu) > 1 {
			t.Lian = t.Seat //一炮多响,放炮者当庄
		} else {
			t.Lian = seat //胡牌玩家连庄
		}
	} else { //黄庄
		t.Lian = t.Seat //最后摸牌玩家
	}
}

//结束牌局重置状态数据
func overSet(t *entity.Desk) {
	t.State = false //牌局状态
	t.Dealer = 0    //庄家重置
	t.Timer = 0     //重置计时
	t.Discard = 0   //重置打牌
	t.Draw = 0      //重置摸牌
	t.Dice = 0      //重置骰子
	t.Seat = 0      //清除位置
	t.Huing = 0     //清除胡牌
	t.Kong = false  //清除杠牌
	t.Trusteeship = make(map[uint32]bool)
	t.Ready = make(map[uint32]bool)
}

//结束牌局,ok=true投票解散
func closeDesk(round, expire uint32, ok bool, t *entity.Desk) {
	var n uint32 = uint32(utils.Timestamp())
	if (round > 0 && expire > n) && !ok {
		return
	}
	if t.CloseCh != nil {
		close(t.CloseCh) //关闭计时器
		t.CloseCh = nil  //消除计时器
	}
	for k, u := range t.Seats {
		u.Pid.Tell(res_leaveDesk(k)) //清除玩家房间数据
		msg := res_leave(k)
		Broadcast(msg, t)
	}
	t.Pid.Tell(res_deskClose(t.Code))
	//RoomsPID.Tell(desk close) //从房间列表中清除
}

//碰操作,已经验证通过
func pong_(card, value, seat uint32, t *entity.Desk) {
	var cards []uint32 = ponging(seat, card, t)
	//碰操作协议消息通知
	msg := res_operate(seat, t.Seat, value, card)
	Broadcast(msg, t)
	//状态设置
	t.Seat = seat                //位置切换
	t.Timer = 0                  //重置计时
	t.Draw = cards[len(cards)-1] //设置摸牌,超时出牌时打出
	t.Discard = 0                //重置出牌
	operateInit(t)               //清除操作记录,操作成功后消除,防止重复提示
	//等待出牌
}

//杠操作,已经验证通过
func kong_(card1, value, seat, v, mask uint32, t *entity.Desk) {
	var card uint32 = uint32(card1)
	switch mask {
	case algo.BU_KONG:
		buKong(seat, card, t)
	case algo.MING_KONG:
		mingKong(seat, card, t)
	case algo.AN_KONG:
		anKong(seat, card, t)
	}
	//杠操作协议消息通知
	msg := res_operate(seat, t.Seat, value, card1)
	Broadcast(msg, t)
	//状态设置
	t.Kong = true //杠操作出牌标识
	t.Seat = seat //位置切换
	t.Timer = 0   //重置计时
	//抢杠处理
	var ok bool = qiangKong(card, mask, t) //检测是否抢杠
	if ok {                                //TODO:优化
		turn(false, t) //抢杠操作
	} else {
		drawcard(t) //摸牌
	}
}

//吃操作,已经验证通过
func chow_(card, value, seat uint32, t *entity.Desk) {
	c1, c2 := algo.DecodeChow2(card)
	var ok bool = chowing(seat, c1, c2, t)
	if !ok {
		glog.Errorf("chow card error -> %d", card)
	}
	//吃操作协议消息通知
	var card2 uint32 = algo.EncodeChow(c1, c2, t.Discard)
	msg := res_operate(seat, t.Seat, value, card2)
	Broadcast(msg, t)
	//状态设置
	t.Seat = seat //位置切换
	t.Timer = 0   //重置计时
	var cards []uint32 = t.HandCards[seat]
	t.Draw = cards[len(cards)-1] //设置摸牌,超时出牌时打出
	t.Discard = 0                //重置出牌
	operateInit(t)               //清除操作记录
	//等待出牌
}

//补杠被抢杠,抢杠处理,抢杠胡牌
func qiangKong(card, mask uint32, t *entity.Desk) bool {
	operateInit(t) //清除操作记录
	//if mask != algo.BU_KONG {
	//	return false
	//}
	//检测(抢杠胡)
	for s, _ := range t.Seats {
		if s == t.Seat { //出牌人跳过
			continue
		}
		var cards []uint32 = getHandCards(s, t)
		v_h := algo.DiscardHuTao(t.Laizi, card, cards) //胡
		v_h |= heType2(v_h, s, card, cards, t)
		//十三幺可以抢任何杠, 其它只能抢补杠
		if algo.QiangKong(v_h, mask) {
			t.Hu[s] = v_h | algo.QIANG_GANG
		}
	}
	if len(t.Hu) > 0 {
		return true
	}
	return false
}

//取消操作时消息通知
func cancelOperate(seat, beseat, value, card uint32, t *entity.Desk) {
	msg := res_operate(seat, beseat, value, card)
	u := getPlayer(seat, t)
	u.Pid.Tell(msg)
}

//--------获取方法
//获取玩家
func getPlayer(seat uint32, t *entity.Desk) *entity.DeskUser {
	if v, ok := t.Seats[seat]; ok && v != nil {
		return v
	}
	panic(fmt.Sprintf("getPlayer error:%d", seat))
}

//获取玩家位置
func getRoles(userid string, t *entity.Desk) uint32 {
	if v, ok := t.Roles[userid]; ok {
		return v
	}
	return 0
}

//获取手牌
func getHandCards(seat uint32, t *entity.Desk) []uint32 {
	if v, ok := t.HandCards[seat]; ok && v != nil {
		return v
	}
	panic(fmt.Sprintf("getHandCards error:%d", seat))
}

//获取海底牌
func getOutCards(seat uint32, t *entity.Desk) []uint32 {
	if v, ok := t.OutCards[seat]; ok && v != nil {
		return v
	}
	return []uint32{}
}

//获取碰牌
func getPongCards(seat uint32, t *entity.Desk) []uint32 {
	if v, ok := t.PongCards[seat]; ok && v != nil {
		return v
	}
	return []uint32{}
}

//获取杠牌
func getKongCards(seat uint32, t *entity.Desk) []uint32 {
	if v, ok := t.KongCards[seat]; ok && v != nil {
		return v
	}
	return []uint32{}
}

//获取吃牌
func getChowCards(seat uint32, t *entity.Desk) []uint32 {
	if v, ok := t.ChowCards[seat]; ok && v != nil {
		return v
	}
	return []uint32{}
}

//获取玩家托管状态
func getTrust(seat uint32, t *entity.Desk) bool {
	if v, ok := t.Trusteeship[seat]; ok {
		return v
	}
	return false
}

//获取玩家准备状态
func getReady(seat uint32, t *entity.Desk) bool {
	if v, ok := t.Ready[seat]; ok {
		return v
	}
	return false
}

//获取剩余局数,结束时间
func getRound(t *entity.Desk) (uint32, uint32) {
	var expire uint32 = 0
	var now int64 = utils.Timestamp()
	if int64(t.Expire) > now {
		expire = t.Expire
	}
	var round uint32 = t.Rounds - t.Round
	if round < 0 {
		round = 0
	}
	return round, expire
}

//获取操作值长度
func getLength(t *entity.Desk) (int, int, int) {
	return len(t.Hu), len(t.Pongkong), len(t.Chow)
}

//房间消息广播(除seat外)
func Broadcast_(seat uint32, msg interface{}, t *entity.Desk) {
	for i, u := range t.Seats {
		if i != seat {
			u.Pid.Tell(msg)
		}
	}
}

//--------操作
//摸牌
func in(seat uint32, card uint32, t *entity.Desk) []uint32 {
	var cards []uint32 = getHandCards(seat, t)
	cards = append(cards, card)
	t.HandCards[seat] = cards
	return cards
}

//玩家出牌
func out(seat uint32, card uint32, t *entity.Desk) {
	var cards []uint32 = getHandCards(seat, t)
	cards = algo.Remove(card, cards)
	t.HandCards[seat] = cards
}

//吃牌操作
func chowing(seat, c1, c2 uint32, t *entity.Desk) bool {
	var c []uint32 = []uint32{t.Discard, c1, c2}
	algo.Sort(c, 0, 2)
	//验证吃, c1,c2,c3有序
	if !algo.VerifyChow(c[0], c[1], c[2]) {
		return false
	}
	var cards []uint32 = getHandCards(seat, t)
	var isExist1 bool = algo.Exist(c1, cards, 1)
	if !isExist1 {
		glog.Errorf("chowing card error -> %d", c1)
		return false
	}
	var isExist2 bool = algo.Exist(c2, cards, 1)
	if !isExist2 {
		glog.Errorf("chowing card error -> %d", c2)
		return false
	}
	cards = algo.Remove(c1, cards)
	cards = algo.Remove(c2, cards)
	var cs []uint32 = getChowCards(seat, t)
	cs = append(cs, algo.EncodeChow(c1, c2, t.Discard))
	t.HandCards[seat] = cards
	t.ChowCards[seat] = cs
	return true
}

//碰牌操作
func ponging(seat, card uint32, t *entity.Desk) []uint32 {
	var cards []uint32 = getHandCards(seat, t)
	var isExist bool = algo.Exist(card, cards, 2)
	if !isExist {
		glog.Errorf("ponging card error -> %d", card)
		return cards
	}
	cards = algo.RemoveN(card, cards, 2)
	var cs []uint32 = getPongCards(seat, t)
	cs = append(cs, algo.EncodePeng(seat, card))
	t.HandCards[seat] = cards
	t.PongCards[seat] = cs
	return cards
}

//暗扛操作
func anKong(seat, card uint32, t *entity.Desk) {
	var cards []uint32 = getHandCards(seat, t)
	var isExist bool = algo.Exist(card, cards, 4)
	if !isExist {
		glog.Errorf("anKong card error -> %d", card)
		return
	}
	cards = algo.RemoveN(card, cards, 4)
	var cs []uint32 = getKongCards(seat, t)
	cs = append(cs, algo.EncodeKong(0, card, algo.AN_KONG))
	t.HandCards[seat] = cards
	t.KongCards[seat] = cs
}

//明杠操作
func mingKong(seat, card uint32, t *entity.Desk) {
	var cards []uint32 = getHandCards(seat, t)
	var isExist bool = algo.Exist(card, cards, 3)
	if !isExist {
		glog.Errorf("mingKong card error -> %d", card)
		return
	}
	cards = algo.RemoveN(card, cards, 3)
	var cs []uint32 = getKongCards(seat, t)
	cs = append(cs, algo.EncodeKong(t.Seat, card, algo.MING_KONG))
	t.HandCards[seat] = cards
	t.KongCards[seat] = cs
}

//补杠操作
func buKong(seat, card uint32, t *entity.Desk) {
	var cards []uint32 = getHandCards(seat, t)
	var isExist bool = algo.Exist(card, cards, 1)
	if !isExist {
		glog.Errorf("buKong card error -> %d", card)
		return
	}
	var pongs []uint32 = getPongCards(seat, t)
	for i, v := range pongs {
		_, c := algo.DecodePeng(v)
		if c == card {
			pongs = append(pongs[:i], pongs[i+1:]...)
			break
		}
	}
	cards = algo.Remove(card, cards)
	var cs []uint32 = getKongCards(seat, t)
	cs = append(cs, algo.EncodeKong(0, card, algo.BU_KONG))
	t.HandCards[seat] = cards
	t.KongCards[seat] = cs
	t.PongCards[seat] = pongs
}
