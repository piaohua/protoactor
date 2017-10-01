package algo

// 全部牌值，
// 高四位表示色值(0:万,1:条,2:饼)，
// 低四位表示1-9的牌值
var CARDS = []uint32{
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
}

//番值
var FAN map[uint32]int32 = make(map[uint32]int32)

//初始化
func init() {
	//牌型
	FAN[HU_SINGLE] = 2          // 十三烂
	FAN[HU_SINGLE_ZI] = 3       // 七星十三烂
	FAN[HU_SEVEN_PAIR_BIG] = 3  // 大七对
	FAN[HU_SEVEN_PAIR] = 4      // 小七对
	FAN[HU_SEVEN_PAIR_KONG] = 4 // 豪华小七对
	//胡牌方式
	FAN[HU_MENQING] = 0 // 门清
	FAN[HU_DANDIAO] = 0 // 单钓

	//牌型
	FAN[HU_PING] = 1          // 平胡
	FAN[HU_TWO_SUIT] = 2      // 混一色
	FAN[HU_PONGPONG] = 2      // 碰碰胡
	FAN[HU_ONE_SUIT] = 5      // 清一色
	FAN[HU_TWO_SUIT_PONG] = 6 // 混一色碰碰胡
	FAN[HU_TWO_YAO] = 7       // 混幺九
	FAN[HU_ONE_YAO] = 9       // 纯幺九
	FAN[HU_ONE_SUIT_PONG] = 9 // 大哥(清一色碰碰胡)
	FAN[HU_THIRTEEN_YAO] = 13 // 十三幺
	FAN[HU_ALL_ZI] = 13       // 字一色
	//胡牌方式
	FAN[ZIMO] = 2           // 自摸
	FAN[RUANMO] = 1         // 软摸
	FAN[HEIMO] = 2          // 黑摸
	FAN[PAOHU] = 1          // 炮胡
	FAN[QIANG_GANG] = 2     // 抢杠,其他家胡你补杠那张牌
	FAN[REPAO] = 2          // 炮胡
	FAN[HU_KONG_FLOWER] = 2 // 杠上开花,杠完牌抓到的第一张牌自摸了
	FAN[TIAN_HU] = 3        // 天胡
	FAN[DI_HU] = 3          // 地胡
	//杠牌
	FAN[MING_KONG] = 1 // 明杠
	FAN[AN_KONG] = 2   // 暗杠
	FAN[BU_KONG] = 1   // 补杠
}

//奖马
var MA_C map[uint32]int = make(map[uint32]int)
var MA_N map[int]map[int]int = make(map[int]map[int]int)

//初始化
func init() {
	MA_C = map[uint32]int{
		0x41: 10,
		0x42: 11,
		0x43: 12,
		0x44: 13,
		0x51: 14,
		0x52: 15,
		0x53: 16,
	}
	MA_N[0] = map[int]int{
		1:  1,
		5:  1,
		9:  1,
		13: 1,
	}
	MA_N[1] = map[int]int{
		2:  1,
		6:  1,
		10: 1,
		14: 1,
	}
	MA_N[2] = map[int]int{
		3:  1,
		7:  1,
		11: 1,
		15: 1,
	}
	MA_N[3] = map[int]int{
		4:  1,
		8:  1,
		12: 1,
		16: 1,
	}
}
