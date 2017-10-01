/**********************************************************
* Author        : piaohua
* Email         : 814004090@qq.com
* Last modified : 2016-12-27 17:00
* Filename      : res.go
* Description   : 玩牌协议消息响应
* *******************************************************/
package desk

import (
	"math"
	"protoactor/algo"
	"protoactor/entity"
	"protoactor/errorcode"
	"protoactor/messages"
)

//超时操作提示响应消息
func res_timerDiscard(userid string, card uint32) interface{} {
	return &messages.TimerDiscardReq{
		Userid: userid,
		Card:   card,
	}
}

//超时操作提示响应消息
func res_timerTurn(seat uint32) interface{} {
	return &messages.TimerTurnReq{Seat: seat}
}

//房间关闭,离开房间
func res_leaveDesk(seat uint32) interface{} {
	return &messages.LeaveDesk{Seat: seat}
}

//房间关闭
func res_deskClose(code string) interface{} {
	return &messages.DeskClose{Code: code}
}

//操作提示响应消息
func res_operate(seat, beseat, value, card uint32) interface{} {
	return &messages.SOperate{
		Error:  0,
		Seat:   seat,
		Card:   card,
		Value:  value,
		Beseat: beseat,
	}
}

//操作提示响应消息
func res_operate2(seat, beseat, value, qiang, card uint32) interface{} {
	return &messages.SOperate{
		Error:       0,
		Seat:        seat,
		Card:        card,
		Value:       value,
		Beseat:      beseat,
		Discontinue: qiang,
	}
}

//打牌响应消息
func res_discard(seat, card uint32) interface{} {
	return &messages.SDiscard{
		Error: 0,
		Card:  card,
		Seat:  seat,
	}
}

func res_discard2() interface{} {
	return &messages.SDiscard{
		Error: errorcode.NotYourTurn,
	}
}

func res_discard3(seat, v, card uint32) interface{} {
	return &messages.SDiscard{
		Error: 0,
		Card:  card,
		Seat:  seat,
		Value: v,
	}
}

//处理前面有玩家胡牌优先操作,如果该玩家跳过胡牌,此协议向有碰和明杠的玩家主动发送
func res_pengkong(seat, v, card uint32) interface{} {
	return &messages.SPengKong{
		Card:  card,
		Seat:  seat,
		Value: v,
	}
}

//进入房间响应消息
func res_camein(seat uint32, t *entity.Desk) interface{} {
	stoc := &messages.SCamein{}
	userinfo := &messages.RoomUser{}
	for k, u := range t.Seats {
		if k != seat {
			continue
		}
		userinfo = &messages.RoomUser{
			Ready:    t.Ready[k],
			Seat:     k,
			Score:    t.Score[u.Userid],
			Userid:   u.Userid,
			Nickname: u.Nickname,
			Sex:      u.Sex,
			Photo:    u.Photo,
			//Coin : u.Coin,
			//Diamond : u.Diamond,
			Coin:    0,
			Diamond: 0,
			Win:     0,
			Lost:    0,
			Ping:    0,
			Piao:    0,
		}
	}
	stoc.Userinfo = userinfo
	return stoc
}

//进入房间响应消息
func res_enter(seat uint32, t *entity.Desk) interface{} {
	round, expire := getRound(t)
	stoc := &messages.SEnterRoom{Error: 0}
	//房间信息
	roominfo := &messages.RoomData{
		Roomid:     t.Rid,
		Rtype:      t.Rtype,
		Rname:      t.Rname,
		Expire:     expire,
		Round:      round,
		Count:      t.Count,
		Invitecode: t.Code,
		Dealer:     t.Dealer,
		Userid:     t.Cid,
		State:      t.State,
		Cards:      uint32(len(t.Cards)),
		Dice:       t.Dice,
		Turn:       t.Seat,
		Laipi:      t.Laipi,
		Laizi:      t.Laizi,
	}
	stoc.Roominfo = roominfo
	//房间内人物数据
	userinfo := []*messages.RoomUser{}
	for k, u := range t.Seats {
		user := &messages.RoomUser{
			Ready:    t.Ready[k],
			Seat:     k,
			Score:    t.Score[u.Userid],
			Userid:   u.Userid,
			Nickname: u.Nickname,
			Sex:      u.Sex,
			Photo:    u.Photo,
			//Coin : u.Coin,
			//Diamond : u.Diamond,
			Coin:    0,
			Diamond: 0,
			Win:     0,
			Lost:    0,
			Ping:    0,
			Piao:    0,
		}
		userinfo = append(userinfo, user)
	}
	stoc.Userinfo = userinfo
	return stoc
}

//重新进入房间响应消息
func res_reEnter(seat uint32, t *entity.Desk) interface{} {
	round, expire := getRound(t)
	stoc := &messages.SEnterRoom{Error: 0}
	//房间信息
	roominfo := &messages.RoomData{
		Roomid:     t.Rid,
		Rtype:      t.Rtype,
		Rname:      t.Rname,
		Expire:     expire,
		Round:      round,
		Count:      t.Count,
		Invitecode: t.Code,
		Dealer:     t.Dealer,
		Userid:     t.Cid,
		State:      t.State,
		Cards:      uint32(len(t.Cards)),
		Dice:       t.Dice,
		Turn:       t.Seat,
		Laipi:      t.Laipi,
		Laizi:      t.Laizi,
	}
	stoc.Roominfo = roominfo
	//操作值
	var value uint32
	if v_h, ok := t.Hu[seat]; ok {
		value = v_h
	} else if v_p, ok := t.Pongkong[seat]; ok {
		value = v_p
	} else if v_c, ok := t.Chow[seat]; ok {
		value = v_c
	}
	//房间内人物数据
	userinfo := []*messages.RoomUser{}
	for k, u := range t.Seats {
		user := &messages.RoomUser{
			Ready:    t.Ready[k],
			Seat:     k,
			Score:    t.Score[u.Userid],
			Userid:   u.Userid,
			Nickname: u.Nickname,
			Sex:      u.Sex,
			Photo:    u.Photo,
			//Coin : u.Coin,
			//Diamond : u.Diamond,
			Coin:    0,
			Diamond: 0,
			Win:     0,
			Lost:    0,
			Ping:    0,
			Piao:    0,
		}
		if k == seat {
			user.Value = value //操作值
			user.Handcards = t.HandCards[k]
		}
		userinfo = append(userinfo, user)
	}
	stoc.Userinfo = userinfo
	//投票信息
	voteinfo := &messages.RoomVote{
		Seat: t.Vote,
	}
	if t.Vote > 0 { //投票中
		for k, v := range t.Votes {
			if v == 0 {
				voteinfo.Agree = append(voteinfo.Agree, k)
			} else {
				voteinfo.Disagree = append(voteinfo.Disagree, k)
			}
		}
	}
	stoc.Voteinfo = voteinfo
	//房间牌面数据
	cardinfo := []*messages.RoomCards{}
	for i, _ := range t.Seats {
		pcard := &messages.RoomCards{}
		pcard.Seat = i
		pongs := getPongCards(i, t) //碰牌数据
		for _, v := range pongs {
			j, card := algo.DecodePeng(v)
			var peng uint32 = j << 24
			peng |= (uint32(card) << 16)
			pcard.Peng = append(pcard.Peng, peng)
		}
		kongs := getKongCards(i, t) //杠牌数据
		for _, v := range kongs {
			j, card, classify := algo.DecodeKong(v)
			var kong uint32 = j << 24
			kong |= (uint32(card) << 16)
			kong |= (classify << 8)
			pcard.Kong = append(pcard.Kong, kong)
		}
		pcard.Outcards = getOutCards(i, t)
		//加上进行中的牌
		if t.Seat == i && t.Discard != 0 {
			pcard.Outcards = append(pcard.Outcards, t.Discard)
		}
	}
	stoc.Cardinfo = cardinfo
	return stoc
}

//游戏开始协议
func res_gamestart(dice, dealer, v, laipi, laizi uint32, cards []uint32) interface{} {
	return &messages.SGamestart{
		Dice:   dice,
		Dealer: dealer,
		Laipi:  laipi,
		Laizi:  laizi,
		Cards:  cards,
		Value:  v,
	}
}

//摸牌协议消息响应消息
func res_draw(seat, kong, value, card uint32) interface{} {
	return &messages.SDraw{
		Card:  card,
		Seat:  seat,
		Kong:  kong,
		Value: value,
	}
}

//玩家准备消息响应消息
func res_ready(seat uint32, ready bool) interface{} {
	return &messages.SReady{
		Error: 0,
		Ready: ready,
		Seat:  seat,
	}
}

func res_update(win, lost, ping, piao uint32) interface{} {
	return &messages.UpdateUser{
		Win:  win,
		Lost: lost,
		Ping: ping,
		Piao: piao,
	}
}

//结束牌局响应消息,huType:0:黄庄，1:自摸，2:炮胡
func res_over(id, huseat, seat, card, round, expire, count uint32,
	handCards map[uint32][]uint32, hu map[uint32]uint32,
	huFan, mingKong, beMingKong, anKong,
	buKong, total, coin, score map[uint32]int32) interface{} {
	stoc := &messages.SGameover{
		Roomid: id,
		Round:  round,
		Expire: expire,
		Seat:   huseat,
		Card:   card,
	}
	var huType uint32 = 0  //胡牌类型
	var paoSeat uint32 = 0 //放冲玩家
	var i uint32
	for i = 1; i <= count; i++ {
		hs := handCards[i] //手牌
		var val uint32 = 0 //胡牌掩码
		if v, ok := hu[i]; ok {
			if v&algo.PAOHU > 0 { //放冲
				huType = 2
				paoSeat = seat
				hs = append(hs, card) //加上打出的牌
			} else if v&algo.ZIMO > 0 { //自摸
				huType = 1
			}
			//val = v //胡牌掩码
			val = res_over_hu(v) //显示处理
		}
		over := &messages.RoomOver{
			Seat: i,
			//Userid: i, //TODO
			Hu: val,
			// 结算时要显示该玩家的手牌
			//Cards: handCards[i],
			Cards: hs,
			Total: total[i],
			Coin:  coin[i],
			Score: score[i],
			//算番, HuTypeFan前端只显示牌型分,所以给的一样
			HuTypeFan:  huFan[i],
			HuFan:      huFan[i],
			MingKong:   mingKong[i],
			BeMingKong: beMingKong[i],
			AnKong:     anKong[i],
			BuKong:     buKong[i],
		}
		stoc.Data = append(stoc.Data, over)
	}
	stoc.PaoSeat = paoSeat
	stoc.HuType = huType
	return stoc
}

//自摸平胡(只显示自摸),点炮(放冲)平胡(只显示点炮)
func res_over_hu(val uint32) uint32 {
	if val&algo.PAOHU > 0 && val&algo.HU_PING > 0 {
		return val ^ algo.HU_PING
	}
	if val&algo.ZIMO > 0 && val&algo.HU_PING > 0 {
		return val ^ algo.HU_PING
	}
	return val
}

//离开房间响应消息
func res_leave(seat uint32) interface{} {
	return &messages.SLeave{
		Error: 0,
		Seat:  seat,
	}
}

/*
//发起投票申请解散房间
func res_voteStart(seat uint32) interface{} {
	stoc := &messages.SLaunchVote{
		Seat: seat),
		Error: 0),
	}
	return stoc
}

//投票解散房间事件结果
func res_voteResult(vote uint32) interface{} {
	stoc := &messages.SVoteResult{
		Vote: vote),
		Error: 0),
	}
	return stoc
}

//投票
func res_vote(seat, vote uint32) interface{} {
	stoc := &messages.SVote{
		Vote: vote),
		Seat: seat),
		Error: 0),
	}
	return stoc
}

//托管消息响应
func res_trust(seat, kind uint32) interface{} {
	stoc := &messages.STrusteeship{
		Seat: seat),
		Kind: kind),
		Error: 0),
	}
	return stoc
}
*/

/*TODO:
抢杠胡	2	玩家补杠时，其他玩家胡牌	由放杠玩家支付*3
杠上花	2	玩家补杠或暗杠以后，摸牌构成胡牌	由三家支付番数
玩家直杠以后，摸牌构成胡牌	由放杠玩家支付*3
*/

/*
庄自摸	    庄炮胡	        闲自摸	            闲炮胡(闲放)	闲炮胡(庄放)
牌型*庄*3	牌型*放炮*庄	牌型+牌型+牌型*庄	牌型*放炮	    牌型*庄*放炮
6*牌型	    4*牌型	        4*牌型	            2*牌型	        4*牌型
庄家自摸时,其他玩家需每人多支付2倍,如:庄家平胡自摸所得番数为:1*2*3
庄家炮胡时,放炮者需支付给庄家2倍,其他玩家不需要给,如:庄家平胡炮胡所得番数为:1*2*2
其他玩家自摸时,庄家需要额外支付2倍,如:闲家平胡自摸所得番数为:1+1+1*2
其他玩家炮胡时,如果是闲家放炮,则庄家不需要给番数,所得番数为:1*2
其他玩家炮胡时,如果是庄家放炮,则所得番数为:1*2*2
*/
//结算,(明杠,放冲,庄家 - 收一家)
func gameOver(t *entity.Desk) (huFan, mingKong, beMingKong, anKong, buKong, total map[uint32]int32) {
	//huTypeFan  = make(map[uint32]int32)// 胡牌方式番数
	huFan = make(map[uint32]int32)      // 胡牌牌型番数
	mingKong = make(map[uint32]int32)   // 闷豆的番数
	beMingKong = make(map[uint32]int32) // 被点豆的负番数
	anKong = make(map[uint32]int32)     // 明豆的番数
	buKong = make(map[uint32]int32)     // 拐弯豆的番数
	total = make(map[uint32]int32)      // 总番数
	var k uint32
	for k = 1; k <= t.Count; k++ {
		//牌型分
		if v, ok := t.Hu[k]; ok {
			//f_t := algo.HuType2Tao(v) //胡牌牌型,多个牌型时相乘
			f_w := algo.HuWay2Tao(v) //胡牌方式,多个方式时相乘
			//牌型分
			//f_tw := (f_t + f_w + int32(t.OutLaizi[k])*2) * int32(t.Ante)
			f_tw := f_w * int32(math.Pow(2, float64(t.OutLaizi[k]))) * int32(t.Ante)
			//牌型分,t.seat=出牌(放冲)位置
			huFan = fanType(t.Dealer, t.Seat, k, t.Count, f_tw, huFan)
		}
		//杠牌分
		kongs := getKongCards(k, t) //杠牌数据
		for _, v := range kongs {   //杠牌分
			i, _, cy := algo.DecodeKong(v) //解码杠值
			f_k := algo.HuKong(cy)         //杠
			if cy == algo.MING_KONG {
				mingKong[k] += f_k * 1     //收一家
				beMingKong[i] += 0 - f_k*1 //被收一家
			} else if cy == algo.BU_KONG {
				buKong = over3(buKong, k, t.Count, f_k*1) //收三家
			} else if cy == algo.AN_KONG {
				anKong = over3(anKong, k, t.Count, f_k*2) //收三家
			}
		}
	}
	//总番数
	for k = 1; k <= t.Count; k++ {
		total[k] += huFan[k] + mingKong[k] + beMingKong[k] + anKong[k] + buKong[k]
	}
	return huFan, mingKong, beMingKong, anKong, buKong, total
}

//倍数(放冲,庄家 - 双倍) TODO:优化
//dealer=庄家位置,paoseat=放炮位置,huseat=胡牌位置
func fanType(dealer, paoseat, huseat, count uint32, f_tw int32,
	hf map[uint32]int32) map[uint32]int32 {
	if paoseat == huseat { //自摸,收三家
		if huseat == dealer { //庄家自摸
			//6 //其它3家*1 (庄家*1倍)
			hf = over3(hf, dealer, count, f_tw*1) //收三家
		} else { //闲家自摸
			//4 //庄家*1 + 其它2家*1
			var i uint32
			for i = 1; i <= count; i++ {
				if i == huseat {
					continue
				}
				if i == dealer {
					hf = over1(hf, huseat, i, f_tw*1) //收一家
				} else {
					hf = over1(hf, huseat, i, f_tw*1) //收一家
				}
			}
		}
	} else { //炮胡,收一家
		if huseat == dealer { //庄家胡(肯定闲家放炮)
			//4//收放炮的*1 (放1倍*庄1倍)
			hf = over1(hf, huseat, paoseat, f_tw*1) //收一家
		} else { //闲家胡
			if paoseat == dealer { //庄家放炮
				//4//收庄家的*1 (放炮1倍*庄家1倍)
				hf = over1(hf, huseat, dealer, f_tw*1) //收一家
			} else { //闲家放炮
				//2//收闲家的*1 (放炮1倍)
				hf = over1(hf, huseat, paoseat, f_tw*1) //收一家
			}
		}
	}
	return hf
}

//收三家,seat=收的位置,val=收的番数
func over3(m map[uint32]int32, seat, count uint32, val int32) map[uint32]int32 {
	var i uint32
	for i = 1; i <= count; i++ {
		var value int32
		if i != seat {
			value = 0 - val //为负数
		} else {
			value = 3 * val //收三家
		}
		m[i] += value
	}
	return m
}

//收一家,s1=收的位置,s2=出的位置,val=收的番数
func over1(m map[uint32]int32, s1, s2 uint32, val int32) map[uint32]int32 {
	m[s1] += val
	m[s2] -= val
	return m
}
