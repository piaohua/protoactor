package algo

// 胡牌类型：平胡，混一色，碰碰胡，清一色，混一色碰碰胡，混幺九，纯幺九，大哥(清一色碰碰胡)，十三幺，字一色
// 胡牌方式：炮胡，自摸，抢杠胡，杠上花，天胡，地胡

//   平胡：普通3n+2牌型
// 混一色：同色数牌+字牌组成3n+2牌型
// 碰碰胡：1对将+n个刻了组成3n+2牌型(即大对子)
// 清一色：同色牌组成3n+2牌型
// 混一色碰碰胡：即满足混一色和碰碰胡组成3n+2牌型
// 混幺九：由数牌1,9+字牌组成3n+2牌型
// 纯幺九：由数牌1,9组成3n+2牌型
// 大哥(清一色碰碰胡)：即满足清一色和碰碰胡组成3n+2牌型
// 十三幺：由1、9序数牌和所有字牌组成3n+2牌型，需要一对作将牌(不许碰杠)
// 字一色：由字牌的刻子或杠组成的牌型组成3n+2牌型
// 注：除十三幺其它都可以碰杠，碰杠手牌都需满足3n+2牌型

// 注：大牌玩法只能自摸,只能平胡,没有天地胡

//赖皮 = 3, 赖子 = 4
//赖子在手牌中充当万能牌,但是打出去不能碰杠胡,也不能充当万能牌参与碰杠,
//打出去立刻摸一张(相当于换牌),手牌中只能有一张,有多张不能胡
//(只能有一张在手牌中充当万能牌胡牌,打出去继续摸牌)

const (
	// 牌局基础的常量
	TOTAL uint32 = 108 //一副贵州麻将的总数
	BING  uint32 = 2   //同子类型
	TIAO  uint32 = 1   //条子类型
	WAN   uint32 = 0   //万字类型
	FENG  uint32 = 4   //风牌类型
	ZI    uint32 = 5   //字牌类型

	HAND uint32 = 13 //手牌数量
	SEAT uint32 = 4  //最多可参与一桌打牌的玩家数量,不算旁观

	// 碰杠胡掩码,用32位每位代表不同的状态
	DRAW      uint32 = 0      // 摸牌
	DISCARD   uint32 = 1      // 打牌
	PENG      uint32 = 2 << 0 // 碰
	MING_KONG uint32 = 2 << 1 // 明杠
	AN_KONG   uint32 = 2 << 2 // 暗杠
	BU_KONG   uint32 = 2 << 3 // 补杠
	KONG      uint32 = 2 << 4 // 杠(代表广义的杠)
	CHOW      uint32 = 2 << 5 // 吃
	HU        uint32 = 2 << 6 // 胡(代表广义的胡)

	//胡牌方式
	ZIMO           uint32 = 2 << 8  // 自摸
	HEIMO          uint32 = 2 << 9  // 黑摸
	RUANMO         uint32 = 2 << 10 // 软摸
	PAOHU          uint32 = 2 << 11 // 炮胡,也叫放冲
	QIANG_GANG     uint32 = 2 << 12 // 抢杠,其他家胡你补杠那张牌
	REPAO          uint32 = 2 << 13 // 热炮,杠后出牌被胡
	HU_KONG_FLOWER uint32 = 2 << 14 // 杠上开花,杠完牌抓到的第一张牌自摸了
	TIAN_HU        uint32 = 2 << 15 // 天胡
	DI_HU          uint32 = 2 << 16 // 地胡

	//胡牌类型
	HU_PING          uint32 = 2 << 17 // 平胡
	HU_TWO_SUIT      uint32 = 2 << 18 // 混一色
	HU_PONGPONG      uint32 = 2 << 19 // 碰碰胡
	HU_ONE_SUIT      uint32 = 2 << 20 // 清一色
	HU_TWO_SUIT_PONG uint32 = 2 << 21 // 混一色碰碰胡
	HU_TWO_YAO       uint32 = 2 << 22 // 混幺九
	HU_ONE_YAO       uint32 = 2 << 23 // 纯幺九
	HU_ONE_SUIT_PONG uint32 = 2 << 24 // 大哥(清一色碰碰胡)
	HU_THIRTEEN_YAO  uint32 = 2 << 25 // 十三幺
	HU_ALL_ZI        uint32 = 2 << 26 // 字一色
)

const (
	//胡牌类型
	HU_SINGLE          uint32 = 0 // 十三烂
	HU_SINGLE_ZI       uint32 = 0 // 七星十三烂
	HU_SEVEN_PAIR_BIG  uint32 = 0 // 大七对
	HU_SEVEN_PAIR      uint32 = 0 // 小七对
	HU_SEVEN_PAIR_KONG uint32 = 0 // 豪华小七对

	HU_MENQING uint32 = 0 // 门清
	HU_DANDIAO uint32 = 0 // 单钓

	//组合牌型
	HU_ONE_SUIT_PAIR_BIG  uint32 = HU_ONE_SUIT | HU_SEVEN_PAIR_BIG
	HU_ONE_SUIT_PAIR      uint32 = HU_ONE_SUIT | HU_SEVEN_PAIR
	HU_ONE_SUIT_PAIR_KONG uint32 = HU_ONE_SUIT | HU_SEVEN_PAIR_KONG
	//---
	HU_ALL_ZI_PAIR_BIG  uint32 = HU_ALL_ZI | HU_SEVEN_PAIR_BIG
	HU_ALL_ZI_PAIR      uint32 = HU_ALL_ZI | HU_SEVEN_PAIR
	HU_ALL_ZI_PAIR_KONG uint32 = HU_ALL_ZI | HU_SEVEN_PAIR_KONG
	//---
	TIAN_HU_ONE_SUIT        uint32 = TIAN_HU | HU_ONE_SUIT
	DI_HU_ONE_SUIT          uint32 = DI_HU | HU_ONE_SUIT
	TIAN_HU_ALL_ZI          uint32 = TIAN_HU | HU_ALL_ZI
	DI_HU_ALL_ZI            uint32 = DI_HU | HU_ALL_ZI
	TIAN_HU_SEVEN_PAIR      uint32 = TIAN_HU | HU_SEVEN_PAIR
	DI_HU_SEVEN_PAIR        uint32 = DI_HU | HU_SEVEN_PAIR
	TIAN_HU_SEVEN_PAIR_KONG uint32 = TIAN_HU | HU_SEVEN_PAIR_KONG
	DI_HU_SEVEN_PAIR_KONG   uint32 = DI_HU | HU_SEVEN_PAIR_KONG
)
