package algo

import (
	"strconv"
	"testing"
)

// 测试
func TestExist(t *testing.T) {
	cards := []uint32{0x07, 0x08, 0x09, 0x18, 0x18, 0x41, 0x43, 0x44}
	value := existHu(cards)
	t.Log(cards, value, strconv.FormatInt(int64(value), 2))
	cs4 := []uint32{0x08, 0x08, 0x18, 0x18, 0x18, 0x41, 0x41, 0x41, 0x41}
	var ok bool = existFour2(cs4)
	t.Log("ok -> ", ok)
	//cs3 := []uint32{0x01, 0x05, 0x09, 0x11, 0x15, 0x19, 0x21, 0x25, 0x29,0x41,0x43,0x51,0x52,0x53}
	//cs3 := []uint32{0x01, 0x05, 0x09, 0x11, 0x19, 0x21, 0x29,0x41,0x42,0x43,0x44,0x51,0x52,0x53}
	cs3 := []uint32{0x01, 0x05, 0x09, 0x11, 0x15, 0x18, 0x21, 0x25, 0x29, 0x41, 0x42, 0x43, 0x52, 0x53}
	var thirteen uint32 = existThirteen(cs3)
	t.Log("thirteen -> ", thirteen)
	//cs7 := []uint32{0x01, 0x01, 0x02, 0x02, 0x11, 0x11, 0x19,0x19,0x21,0x21,0x44,0x44,0x52,0x52}
	cs7 := []uint32{0x01, 0x01, 0x01, 0x01, 0x11, 0x11, 0x19, 0x19, 0x21, 0x21, 0x44, 0x44, 0x52, 0x52}
	var seven uint32 = exist7pair(cs7)
	t.Log("seven -> ", seven)
	cs7kong := []uint32{0x01, 0x01, 0x01, 0x01, 0x11, 0x11, 0x21, 0x21, 0x21, 0x21, 0x44, 0x44, 0x44, 0x44}
	var haveKong int = haveKong(cs7kong)
	t.Log("haveKong -> ", haveKong)
}

// 大七对, 4刻子+2, 不可以吃/清一色/字一色
// 有序slices cs 包含吃，碰，杠
func TestHuTypeDetect(t *testing.T) {
	//大七对
	//cs := []uint32{0x01, 0x01, 0x01, 0x11, 0x11, 0x11, 0x21, 0x21, 0x21,0x44,0x44, 0x51, 0x51, 0x51}
	//清一色
	//cs := []uint32{0x01, 0x02, 0x03, 0x04, 0x04, 0x04,0x05,0x05,0x05,0x06,0x07,0x08,0x09,0x09}
	//字一色
	//cs := []uint32{0x41, 0x41, 0x42, 0x42, 0x43, 0x43, 0x44, 0x44, 0x51,0x51,0x51,0x52,0x52,0x52}
	//cs := []uint32{0x11,0x11,0x11,0x11,0x24,0x24,0x24,0x44,0x44,0x44,0x07,0x07,0x14,0x14,0x14}
	//cs := []uint32{0x01,0x02,0x03,0x04, 0x14,0x15,0x16, 0x27,0x27,0x27,0x41,0x41,0x41}
	//cs := []uint32{0x16,0x17,0x18,0x24,0x24,0x26,0x27,0x28,0x42,0x43,0x44,0x51,0x52,0x53}
	//cs := []uint32{0x01,0x01,0x04,0x04,0x04,0x14,0x15,0x16,0x25,0x25,0x25,0x17,0x17,0x17}
	//cs := []uint32{0x11,0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x07,0x08,0x09,0x41,0x43,0x44}
	//cs := []uint32{0x03,0x03,0x03,0x03,0x05,0x05,0x22,0x22,0x24,0x24,0x15,0x15,0x53,0x53}
	cs := []uint32{0x51, 0x52, 0x53, 0x51, 0x52, 0x53, 0x51, 0x52, 0x53, 0x51, 0x52, 0x53, 0x54, 0x54}
	var val uint32 = HuTypeDetect(false, false, true, cs)
	t.Log("val -> ", val)
}

// 测试
func TestNextSeat(t *testing.T) {
	t.Log("NextSeat -> ", NextSeat(1))
	t.Log("NextSeat -> ", NextSeat(2))
	t.Log("NextSeat -> ", NextSeat(3))
	t.Log("NextSeat -> ", NextSeat(4))
}

// 测试
func TestRemove(t *testing.T) {
	cs2chow := []uint32{0x01, 0x01, 0x01, 0x01, 0x11, 0x11, 0x21, 0x21, 0x21, 0x21, 0x44, 0x44, 0x44, 0x44}
	cs2chow = Remove(0x21, cs2chow)
	t.Log("cs2chow -> ", cs2chow, " len -> ", len(cs2chow))
	cs2chow = Remove(0x11, cs2chow)
	t.Log("cs2chow -> ", cs2chow, " len -> ", len(cs2chow))
	cs2chow = Remove(0x44, cs2chow)
	t.Log("cs2chow -> ", cs2chow, " len -> ", len(cs2chow))
}

// 测试
func TestCode(t *testing.T) {
	//c2 := uint32(5460561 >> 8 & 0xFF)
	//t.Log("c2 -> ", c2)
	//card := EncodeChow(0x23,0x24,0x25)
	//t.Logf("card -> %d ", card)
	//var card uint32 = 1381651
	//var card uint32 = 21073
	//var card uint32 = 5396
	var card uint32 = 1382680
	c1, c2, c3 := DecodeChow(card)
	t.Log("c1, c2, c3 -> ", c1, c2, c3)
	//c1, c2 := DecodeChow2(card)
	//t.Log("c1, c2 -> ", c1, c2)
	//var card2 uint32 = 1447444
	//var card2 uint32 = (22 << 8) | 20
	var card2 uint32 = 5652
	t.Logf("card2 -> %d ", card2)
	c21, c22 := DecodeChow2(card2)
	t.Log("c1, c2 -> ", c21, c22)
}

// 测试胡牌
func Test_hu(t *testing.T) {
	cards := []uint32{0x07, 0x08, 0x09, 0x18, 0x18, 0x41, 0x43, 0x44}
	value := existHu(cards)
	t.Log(cards, value, strconv.FormatInt(int64(value), 2))
	ok := Exist(0x18, cards, 3)
	t.Log("ok -> ", ok)
	cs, st := RemoveE(0x18, cards, 3)
	t.Log("cs -> ", cs, " st -> ", st)
}

// 测试胡牌
func TestMask(t *testing.T) {
	var v uint32 = 0x80082
	//胡牌方式
	if v&QIANG_GANG > 0 {
		t.Log("QIANG_GANG -> ", v)
	}
	if v&HU_KONG_FLOWER > 0 {
		t.Log("HU_KONG_FLOWER -> ", v)
	}
	if v&HU_MENQING > 0 {
		t.Log("HU_MENQING -> ", v)
	}
	if v&HU_DANDIAO > 0 {
		t.Log("HU_DANDIAO -> ", v)
	}
	if v&TIAN_HU > 0 {
		t.Log("TIAN_HU -> ", v)
	}
	if v&DI_HU > 0 {
		t.Log("DI_HU -> ", v)
	}
	//牌型
	if v&HU_PING > 0 {
		t.Log("HU_PING -> ", v)
	}
	if v&HU_SINGLE > 0 {
		t.Log("HU_SINGLE -> ", v)
	}
	if v&HU_SINGLE_ZI > 0 {
		t.Log("HU_SINGLE_ZI -> ", v)
	}
	if v&HU_SEVEN_PAIR_BIG > 0 {
		t.Log("HU_SEVEN_PAIR_BIG -> ", v)
	}
	if v&HU_SEVEN_PAIR > 0 {
		t.Log("HU_SEVEN_PAIR -> ", v)
	}
	if v&HU_SEVEN_PAIR_KONG > 0 {
		t.Log("HU_SEVEN_PAIR_KONG -> ", v)
	}
	if v&HU_ONE_SUIT > 0 {
		t.Log("HU_ONE_SUIT -> ", v)
	}
	if v&HU_ALL_ZI > 0 {
		t.Log("HU_ALL_ZI -> ", v)
	}
	//胡
	if v&PAOHU > 0 {
		t.Log("PAOHU -> ", v)
	}
	if v&ZIMO > 0 {
		t.Log("ZIMO -> ", v)
	}
	//操作
	if v&HU > 0 {
		t.Log("HU -> ", v)
	}
	if v&PENG > 0 {
		t.Log("PENG -> ", v)
	}
	if v&MING_KONG > 0 {
		t.Log("MING_KONG -> ", v)
	}
	if v&AN_KONG > 0 {
		t.Log("AN_KONG -> ", v)
	}
	if v&BU_KONG > 0 {
		t.Log("BU_KONG -> ", v)
	}
	if v&KONG > 0 {
		t.Log("KONG -> ", v)
	}
	if v&CHOW > 0 {
		t.Log("CHOW -> ", v)
	}
	t.Log("v -> ", v)
	//胡牌方式
}

// 测试
func TestMaCard(t *testing.T) {
	t.Log("MA_C -> ", MA_C)
	t.Log("MA_N -> ", MA_N)
	macards := []uint32{0x17, 0x23, 0x44, 0x52, 0x42}
	cs := MaCards(1, 3, macards)
	t.Logf("cs -> %+x", cs)
}

// 测试胡牌
func TestLaizi(t *testing.T) {
	cards := []uint32{0x07, 0x08, 0x09, 0x18, 0x18, 0x22, 0x23, 0x24}
	value := existHu3n2_normal(cards, len(cards))
	t.Log(cards, value, strconv.FormatInt(int64(value), 2))
	//cards = []uint32{0x07, 0x08, 0x08, 0x09, 0x09, 0x18, 0x18, 0x22, 0x23, 0x24}
	cards = []uint32{0x07, 0x07, 0x07, 0x07}
	value = existHu3n2_laizi(cards, len(cards))
	t.Log(cards, value, strconv.FormatInt(int64(value), 2))
	cards = []uint32{0x07, 0x07, 0x07, 0x08, 0x22}
	value, _ = existHuTao(cards, 0x08, len(cards))
	t.Log(cards, value, strconv.FormatInt(int64(value), 2))
	//
	cards = []uint32{0x07, 0x07, 0x07, 0x06, 0x07, 0x18, 0x22, 0x22}
	value, _ = existHuTao(cards, 0x18, len(cards))
	t.Log(cards, value, strconv.FormatInt(int64(value), 2))
}
