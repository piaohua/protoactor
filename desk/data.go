package desk

import "protoactor/entity"

var DT int = 12 //出牌超时时间
var OT int = 10 //操作超时时间

func NewDeskData(rid, rounds, expire, rtype, ante, cost, payment, pao,
	count, ctime uint32, creator, rname, invitecode string) *entity.DeskData {
	return &entity.DeskData{
		Rid:     rid,
		Rtype:   rtype,
		Rname:   rname,
		Ante:    ante,
		Payment: payment,
		Cost:    cost,
		Cid:     creator,
		Expire:  expire,
		Rounds:  rounds,
		Code:    invitecode,
		Pao:     pao,
		Count:   count,
		CTime:   ctime,
		Score:   make(map[string]int32),
		Roles:   make(map[string]uint32),
		Seats:   make(map[uint32]*entity.DeskUser),
	}
}

//加入房间
func Add(u *entity.DeskUser, d *entity.DeskData) bool {
	if uint32(len(d.Roles)) >= d.Count {
		return false
	}
	if _, ok := d.Roles[u.Userid]; ok {
		return true //已经在房间
	}
	var seat uint32 = add_(u.Userid, d)
	if seat == 0 {
		return false
	}
	d.Seats[seat] = u
	return true
}

//加入房间,返回位置
func add_(userid string, d *entity.DeskData) uint32 {
	var i uint32
	for i = 1; i <= d.Count; i++ {
		if _, ok := d.Seats[i]; !ok {
			d.Roles[userid] = i
			return i
		}
	}
	return 0
}
