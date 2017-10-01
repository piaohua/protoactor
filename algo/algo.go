/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2016-12-23 10:00
 * Filename      : algo.go
 * Description   : 玩牌算法
 * *******************************************************/
package algo

//// external function

// 打牌检测,胡牌, 放炮胡检测
func DiscardHu(card uint32, cs []uint32) uint32 {
	cards := make([]uint32, len(cs)+1)
	copy(cards, cs)
	cards[len(cards)-1] = card
	//胡
	var status uint32 = existHu(cards)
	if status > 0 {
		status |= PAOHU
	}
	return status
}

// 打牌检测,胡牌,放炮胡检测
func DiscardHuTao(laizi, card uint32, cs []uint32) uint32 {
	cards := make([]uint32, len(cs)+1)
	copy(cards, cs)
	cards[len(cards)-1] = card
	//胡
	status, _ := existHuTao(cards, laizi, len(cards))
	if status > 0 {
		status |= PAOHU
	}
	return status
}

//平胡不能炮胡
func PaoHuPing(status uint32) uint32 {
	if status == (HU | HU_PING | PAOHU) {
		return 0
	}
	return status
}

//十三幺可以抢任何杠, 其它只能抢补杠
func QiangKong(status, mask uint32) bool {
	if status == 0 {
		return false
	}
	if status&HU_THIRTEEN_YAO > 0 {
		return true
	}
	if mask == BU_KONG {
		return true
	}
	return false
}

//计算马牌(庄,胡牌玩家,马牌) TODO:优化
func MaCards(dealer, seat uint32, macards []uint32) []uint32 {
	var cs []uint32 = []uint32{}
	if len(macards) == 0 {
		return cs
	}
	//庄到胡牌玩家距离为key(逆时针)
	var key uint32 = between(dealer, seat)
	var ms map[int]int = make(map[int]int)
	if ms2, ok := MA_N[int(key)]; !ok {
		panic(key)
	} else {
		ms = ms2
	}
	for _, v := range macards {
		if !type_zi(v) { //字牌
			var v1 int
			if v2, ok := MA_C[v]; !ok {
				panic(v)
			} else {
				v1 = v2
			}
			if _, ok := ms[v1]; ok {
				cs = append(cs, v)
			}
		} else {
			if _, ok := ms[int(v&0xf)]; ok {
				cs = append(cs, v)
			}
		}
	}
	return cs
}

//两个数相差值,返回绝对值
func between(n, m uint32) uint32 {
	if n > m {
		return n - m
	}
	return m - n
}

//去掉平胡(有大牌型时去掉平胡番值,在牌局时才去掉)
func CancelHuPing(status uint32) uint32 {
	if status&HU_PING == 0 {
		return status
	}
	switch {
	case status&HU_SINGLE > 0:
		status ^= HU_PING //去掉平胡
	case status&HU_SINGLE_ZI > 0:
		status ^= HU_PING //去掉平胡
	case status&HU_SEVEN_PAIR_BIG > 0:
		status ^= HU_PING //去掉平胡
	case status&HU_SEVEN_PAIR > 0:
		status ^= HU_PING //去掉平胡
	case status&HU_SEVEN_PAIR_KONG > 0:
		status ^= HU_PING //去掉平胡
	case status&HU_ONE_SUIT > 0:
		status ^= HU_PING //去掉平胡
	case status&HU_ALL_ZI > 0:
		status ^= HU_PING //去掉平胡
	case status&TIAN_HU > 0:
		status ^= HU_PING //去掉平胡
	case status&DI_HU > 0:
		status ^= HU_PING //去掉平胡
	}
	return status
}

// 打牌检测,明杠/碰牌
func DiscardPong(card uint32, cs []uint32) uint32 {
	var status uint32
	//碰杠
	if existMingKong(card, cs) {
		status |= MING_KONG
		status |= KONG
	}
	if existPeng(card, cs) {
		status |= PENG
	}
	return status
}

// 打牌检测,吃
func DiscardChow(s1, s2 uint32, card uint32, cs []uint32) uint32 {
	var status uint32
	//吃
	if NextSeat(s1) == s2 {
		if existChow(card, cs) {
			status |= CHOW
		}
	}
	return status
}

// 摸牌检测,胡牌／暗杠／补杠
func DrawDetectTao(laizi, card uint32, cs []uint32, ps []uint32) uint32 {
	le := len(cs)
	cards := make([]uint32, le)
	copy(cards, cs)
	//自摸胡检测
	status, mo := existHuTao(cards, laizi, le)
	if status > 0 {
		if mo {
			status |= ZIMO | RUANMO
		} else {
			status |= ZIMO | HEIMO
		}
	}
	if len(existAnKongTao(laizi, cards)) > 0 {
		status |= AN_KONG | KONG
	} else if existBuKong(card, ps) {
		status |= BU_KONG | KONG
	}
	return status
}

// 摸牌检测,胡牌／暗杠／补杠
func DrawDetect2(card uint32, cs []uint32, ps []uint32) uint32 {
	le := len(cs)
	cards := make([]uint32, le)
	copy(cards, cs)
	var status uint32
	//自摸胡检测
	status = existHu(cards)
	if status > 0 {
		status |= ZIMO
	}
	if len(existAnKong(cards)) > 0 {
		status |= AN_KONG
		status |= KONG
	} else if existBuKong(card, ps) {
		status |= BU_KONG
		status |= KONG
	}
	return status
}

// 摸牌检测,胡牌／暗杠／补杠
func DrawDetect(card uint32, cs []uint32, ch, ps, ks []uint32) uint32 {
	le := len(cs)
	cards := make([]uint32, le)
	copy(cards, cs)
	var status uint32
	//自摸胡检测
	status = existHu(cards)
	if status > 0 {
		status |= ZIMO
	}
	if status > 0 {
		if le == 14 || menQing(ch, ps, ks) {
			status |= HU_MENQING
			if status&HU_PING > 0 {
				status ^= HU_PING //去掉平胡
			}
		}
	}
	if len(existAnKong(cards)) > 0 {
		status |= AN_KONG
		status |= KONG
	} else if existBuKong(card, ps) {
		status |= BU_KONG
		status |= KONG
	}
	return status
}

// 门清,(检测是否全暗杠或全没有返回true),TODO:优化,可以放在结算时判断
func menQing(cs, ps, ks []uint32) bool {
	if len(cs) > 0 {
		return false
	}
	if len(ps) > 0 {
		return false
	}
	if len(ks) == 0 {
		return true
	}
	for _, v2 := range ks { //杠
		_, _, v := DecodeKong(v2) //解码
		if v != AN_KONG {         //不是暗杠
			return false
		}
	}
	return true
}

// 有暗杠时不算点炮(放冲,单钓)
func DanDiao(v uint32, ks []uint32) uint32 {
	if v&HU_DANDIAO == 0 {
		return v
	}
	if len(ks) == 0 {
		return v
	}
	for _, v2 := range ks { //杠
		_, _, v1 := DecodeKong(v2) //解码
		if v1 == AN_KONG {         //有暗杠
			return v ^ HU_DANDIAO //去掉暗杠
		}
	}
	return v
}

// 胡牌牌型检测
func HuTypeDetect(hu, chow, kong bool, cs []uint32) uint32 {
	Sort(cs, 0, len(cs)-1)
	return existHuType(hu, chow, kong, cs)
}

// 胡牌牌型检测 (广东)
func HuTypeDetect2(cs []uint32) uint32 {
	return existHuTypeYao(cs)
}

// 检测手牌是否有杠
func KongDetect(cs []uint32) bool {
	var n uint32 = kongs(cs)
	return n > 0
}

// 验证吃, c1,c2,c3有序
func VerifyChow(c1, c2, c3 uint32) bool {
	if c1 != c2 && c2 != c3 {
		if c1+1 == c2 && c2+1 == c3 {
			return true
		}
		if c1 >= 0x41 && c3 <= 0x44 {
			return true
		}
	}
	return false
}

// 碰杠吃数据
func EncodePeng(seat uint32, card uint32) uint32 {
	seat = seat << 8
	seat |= uint32(card)
	return seat
}

func DecodePeng(value uint32) (seat uint32, card uint32) {
	seat = value >> 8
	card = uint32(value & 0xFF)
	return
}

func EncodeKong(seat uint32, card uint32, value uint32) uint32 {
	value = value << 16
	value |= (seat << 8)
	value |= uint32(card)
	return value
}

func DecodeKong(value uint32) (seat uint32, card uint32, v uint32) {
	v = value >> 16
	seat = (value >> 8) & 0xFF
	card = uint32(value & 0xFF)
	return
}

func EncodeChow(c1, c2, c3 uint32) (value uint32) {
	value = uint32(c1) << 16
	value |= uint32(c2) << 8
	value |= uint32(c3)
	return
}

func DecodeChow(value uint32) (c1, c2, c3 uint32) {
	c1 = uint32(value >> 16)
	c2 = uint32(value >> 8 & 0xFF)
	c3 = uint32(value & 0xFF)
	return
}

func DecodeChow2(value uint32) (c1, c2 uint32) {
	c1 = uint32(value >> 8)
	c2 = uint32(value & 0xFF)
	return
}

// 正常流程走牌令牌移到下一家
func NextSeat(seat uint32) uint32 {
	if seat == 4 {
		return 1
	}
	return seat + 1
}

// 是否存在n个牌
func Exist(c uint32, cs []uint32, n int) bool {
	for _, v := range cs {
		if n == 0 {
			return true
		}
		if c == v {
			n--
		}
	}
	return n == 0
}

// 移除一个牌
func Remove(c uint32, cs []uint32) []uint32 {
	for i, v := range cs {
		if c == v {
			cs = append(cs[:i], cs[i+1:]...)
			break
		}
	}
	return cs
}

// 移除n个牌,返回是否存在n个牌
func RemoveE(c uint32, cs []uint32, n int) ([]uint32, bool) {
	var m int = n
	for n > 0 {
		for i, v := range cs {
			if c == v {
				cs = append(cs[:i], cs[i+1:]...)
				m--
				break
			}
		}
		n--
	}
	return cs, n == m
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

// 对牌值从小到大排序，采用快速排序算法
func Sort(arr []uint32, start, end int) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for arr[i] < key {
				i++
			}
			for arr[j] > key {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		if start < j {
			Sort(arr, start, j)
		}
		if end > i {
			Sort(arr, i, end)
		}
	}
}

//组合牌型番值(取最大值),TODO:优化
func huType2(v uint32) int32 {
	var f int32 = 0
	if v&HU_ONE_SUIT > 0 {
		f = max(f, huTypeOneSuit(v^HU_ONE_SUIT))
	}
	if v&HU_ALL_ZI > 0 {
		f = max(f, huTypeAllZi(v^HU_ALL_ZI))
	}
	if v&TIAN_HU > 0 {
		f = max(f, huTypeTian(v^TIAN_HU))
	}
	if v&DI_HU > 0 {
		f = max(f, huTypeDi(v^DI_HU))
	}
	return f
}

func huTypeOneSuit(v uint32) int32 {
	var f int32 = 0
	if v&HU_ONE_SUIT_PAIR_BIG > 0 {
		f = max(f, FAN[HU_ONE_SUIT_PAIR_BIG])
	}
	if v&HU_ONE_SUIT_PAIR > 0 {
		f = max(f, FAN[HU_ONE_SUIT_PAIR])
	}
	if v&HU_ONE_SUIT_PAIR_KONG > 0 {
		f = max(f, FAN[HU_ONE_SUIT_PAIR_KONG])
	}
	return f
}

func huTypeAllZi(v uint32) int32 {
	var f int32 = 0
	if v&HU_ALL_ZI_PAIR_BIG > 0 {
		f = max(f, FAN[HU_ALL_ZI_PAIR_BIG])
	}
	if v&HU_ALL_ZI_PAIR > 0 {
		f = max(f, FAN[HU_ALL_ZI_PAIR])
	}
	if v&HU_ALL_ZI_PAIR_KONG > 0 {
		f = max(f, FAN[HU_ALL_ZI_PAIR_KONG])
	}
	return f
}

func huTypeTian(v uint32) int32 {
	var f int32 = 0
	if v&TIAN_HU_ONE_SUIT > 0 {
		f = max(f, FAN[TIAN_HU_ONE_SUIT])
	}
	if v&TIAN_HU_ALL_ZI > 0 {
		f = max(f, FAN[TIAN_HU_ALL_ZI])
	}
	if v&TIAN_HU_SEVEN_PAIR > 0 {
		f = max(f, FAN[TIAN_HU_SEVEN_PAIR])
	}
	if v&TIAN_HU_SEVEN_PAIR_KONG > 0 {
		f = max(f, FAN[TIAN_HU_SEVEN_PAIR_KONG])
	}
	return f
}

func huTypeDi(v uint32) int32 {
	var f int32 = 0
	if v&DI_HU_ONE_SUIT > 0 {
		f = max(f, FAN[DI_HU_ONE_SUIT])
	}
	if v&DI_HU_ALL_ZI > 0 {
		f = max(f, FAN[DI_HU_ALL_ZI])
	}
	if v&DI_HU_SEVEN_PAIR > 0 {
		f = max(f, FAN[DI_HU_SEVEN_PAIR])
	}
	if v&DI_HU_SEVEN_PAIR_KONG > 0 {
		f = max(f, FAN[DI_HU_SEVEN_PAIR_KONG])
	}
	return f
}

func max(n, m int32) int32 {
	if n > m {
		return n
	}
	return m
}

//算番(牌型) TODO:优化
func HuType(v uint32, cs []uint32) int32 {
	var f int32 = 0
	f = huType2(v)
	if f > 0 {
		return f //组合牌型时只取组合牌型
	}
	f = 1 //相乘时为1
	//var f int32 = 1 //相加时应该为0
	if v&HU_PING > 0 {
		f *= FAN[HU_PING]
	}
	if v&HU_SINGLE > 0 {
		f *= FAN[HU_SINGLE]
	}
	if v&HU_SINGLE_ZI > 0 {
		f *= FAN[HU_SINGLE_ZI]
	}
	if v&HU_SEVEN_PAIR_BIG > 0 {
		f *= FAN[HU_SEVEN_PAIR_BIG]
	}
	if v&HU_SEVEN_PAIR > 0 {
		f *= FAN[HU_SEVEN_PAIR]
	}
	if v&HU_SEVEN_PAIR_KONG > 0 {
		var n uint32 = kongs(cs)
		f *= max(16, FAN[HU_SEVEN_PAIR_KONG]*(1<<n)) //n=手牌中暗杠数:4*2^n,最大16番
	}
	if v&HU_ONE_SUIT > 0 {
		f *= FAN[HU_ONE_SUIT]
	}
	if v&HU_ALL_ZI > 0 { //没牌型4番,有8番
		if v == (HU|HU_ALL_ZI|ZIMO) ||
			v == (HU|HU_ALL_ZI|PAOHU) {
			f *= FAN[HU_ALL_ZI] / 2
		} else {
			f *= FAN[HU_ALL_ZI]
		}
	}
	//天地胡时取最大值
	if v&TIAN_HU > 0 {
		f = max(f, FAN[TIAN_HU])
	}
	if v&DI_HU > 0 {
		f = max(f, FAN[DI_HU])
	}
	return f
}

//算番(胡牌方式) TODO:优化
func HuWay(v uint32) int32 {
	var f int32 = 1
	if v&QIANG_GANG > 0 {
		f *= FAN[QIANG_GANG]
	}
	if v&HU_KONG_FLOWER > 0 {
		f *= FAN[HU_KONG_FLOWER]
	}
	if v&HU_MENQING > 0 {
		if menQingFan(v) {
			f *= FAN[HU_MENQING]
		}
	}
	if v&HU_DANDIAO > 0 {
		f *= FAN[HU_DANDIAO]
	}
	//if v&TIAN_HU > 0 {
	//	f *= FAN[TIAN_HU]
	//}
	//if v&DI_HU > 0 {
	//	f *= FAN[DI_HU]
	//}
	return f
}

//算番(杠)
func HuKong(v uint32) int32 {
	return FAN[v]
}

//天地胡去掉平胡
func HuType2TianDi(v uint32, f int32) int32 {
	if v&HU_PING > 0 {
		f -= FAN[HU_PING]
	}
	return f
}

//算番(牌型) TODO:优化
func HuType2Yao(v uint32) int32 {
	switch {
	case v&HU_TWO_SUIT > 0:
		return FAN[HU_TWO_SUIT]
	case v&HU_PONGPONG > 0:
		return FAN[HU_PONGPONG]
	case v&HU_ONE_SUIT > 0:
		return FAN[HU_ONE_SUIT]
	case v&HU_TWO_SUIT_PONG > 0:
		return FAN[HU_TWO_SUIT_PONG]
	case v&HU_TWO_YAO > 0:
		return FAN[HU_TWO_YAO]
	case v&HU_ONE_YAO > 0:
		return FAN[HU_ONE_YAO]
	case v&HU_ONE_SUIT_PONG > 0:
		return FAN[HU_ONE_SUIT_PONG]
	case v&HU_THIRTEEN_YAO > 0:
		return FAN[HU_THIRTEEN_YAO]
	case v&HU_ALL_ZI > 0:
		return FAN[HU_ALL_ZI]
	case v&HU_PING > 0:
		return FAN[HU_PING]
	default:
		return 1
	}
}

//算番(胡牌方式) TODO:优化
func HuWay2Yao(v uint32) (int32, bool) {
	switch {
	case v&QIANG_GANG > 0:
		return FAN[QIANG_GANG], false
	case v&HU_KONG_FLOWER > 0:
		return FAN[HU_KONG_FLOWER], false
	case v&TIAN_HU > 0:
		return FAN[TIAN_HU], true
	case v&DI_HU > 0:
		return FAN[DI_HU], true
	case v&ZIMO > 0:
		return FAN[ZIMO], false
	case v&PAOHU > 0:
		return FAN[PAOHU], false
	default:
		return 1, false
	}
}

//(小七对,十三烂, 天地胡)不算门清番
func menQingFan(v uint32) bool {
	if v&HU_SEVEN_PAIR > 0 {
		return false
	}
	if v&HU_SINGLE > 0 {
		return false
	}
	if v&HU_SINGLE_ZI > 0 {
		return false
	}
	if v&TIAN_HU > 0 {
		return false
	}
	if v&DI_HU > 0 {
		return false
	}
	return true
}

//// internal function

// 手牌有多少个杠牌
func kongs(cards []uint32) uint32 {
	var i uint32 = 0
	var m = make(map[uint32]int)
	for _, v := range cards {
		if n, ok := m[v]; ok {
			m[v] = n + 1
		} else {
			m[v] = 1
		}
	}
	for _, v := range m {
		if v == 4 {
			i += 1
		}
	}
	return i
}

//是否存在暗杠
func existAnKong(cards []uint32) (kong []uint32) {
	le := len(cards)
	for j := 0; j < le-3; j++ {
		count := 0
		for i := j + 1; i < le; i++ {
			if cards[j] == cards[i] {
				count = count + 1
				if count == 3 {
					kong = append(kong, cards[i])
					break
				}
			}
		}
	}
	return
}

//是否存在暗杠
func existAnKongTao(laizi uint32, cards []uint32) (kong []uint32) {
	le := len(cards)
	for j := 0; j < le-3; j++ {
		count := 0
		for i := j + 1; i < le; i++ {
			if cards[j] == cards[i] && cards[j] != laizi {
				count = count + 1
				if count == 3 {
					kong = append(kong, cards[i])
					break
				}
			}
		}
	}
	return
}

//是否存在碰
func existPeng(card uint32, cards []uint32) bool {
	le := len(cards)
	count := 0
	for i := 0; i < le; i++ {
		if card == cards[i] {
			count = count + 1
			if count == 2 {
				return true
			}
		}
	}
	return false
}

//是否存在补杠
func existBuKong(card uint32, pongs []uint32) bool {
	le := len(pongs)
	for i := 0; i < le; i++ {
		_, c := DecodePeng(pongs[i])
		if card == c {
			return true
		}
	}
	return false
}

//是否存在明杠
func existMingKong(card uint32, cards []uint32) bool {
	le := len(cards)
	count := 0
	for i := 0; i < le; i++ {
		if card == cards[i] {
			count = count + 1
			if count == 3 {
				return true
			}
		}
	}
	return false
}

//检测,是否存在吃 TODO:返回吃牌组合 [[1,2,3],[2,3,4],[3,4,5]]
func existChow(card uint32, cs []uint32) bool {
	le := len(cs)
	if le == 1 {
		return false
	}
	cards := make([]uint32, le)
	copy(cards, cs)
	var t uint32 = card >> 4
	if uint32(t) >= FENG {
		return existChow1(t, card, cards)
	}
	return existChow2(card, cards)
}

//风牌,字牌两张不同即可
func existChow1(t, card uint32, cs []uint32) bool {
	count := 0
	var c uint32
	for _, v := range cs {
		if v != card && v != c && v>>4 == t {
			count += 1
			c = v
		}
	}
	if count >= 2 {
		return true
	}
	return false
}

//数牌三张
func existChow2(card uint32, cs []uint32) bool {
	var cards []uint32 = make([]uint32, 5)
	cards[2] = card
	for _, v := range cs {
		var s int = int(card) - int(v) + 2
		if s >= 0 && s <= 4 {
			cards[s] = v
		}
	}
	count := 0
	for _, v := range cards {
		if count >= 3 {
			break
		}
		if v > 0 {
			count += 1
		} else {
			count = 0
		}
	}
	if count >= 3 {
		return true
	}
	return false
}

//有序slices,判断是否有4个相同的牌
func existFour2(cards []uint32) bool {
	var c uint32
	var i int
	for _, v := range cards {
		if c == v {
			i += 1
		} else {
			c = v
			i = 1
		}
		if i == 4 {
			return true
		}
	}
	return false
}

//有序slices,包含多少个杠
func haveKong(cards []uint32) int {
	var c uint32
	var i int
	var n int
	for _, v := range cards {
		if c == v {
			i += 1
		} else {
			c = v
			i = 1
		}
		if i == 4 {
			n += 1
		}
	}
	return n
}

// 胡牌后检测
// 清一色(同一花色,必须有牌型,没有字牌)/字一色(全部是字牌)
// 大七对, 1对子(将)+(刻子,杠) + 不能吃 + 杠(非手牌中杠)
// 有序slices cs 包含吃，碰，杠, chow=false没有吃,=true有吃
// kong=false没有杠,=true有杠, hu=false没有牌型,=true有牌型
func existHuType(hu, chow, kong bool, cs []uint32) uint32 {
	var all_zi bool = true
	var one_suit bool = true
	var seven_pair_big bool = true
	var b bool
	var c uint32
	var huType uint32
	var m = make(map[uint32]int)
	for _, v := range cs {
		if n, ok := m[v]; ok {
			m[v] = n + 1
		} else {
			m[v] = 1
		}
		if !one_suit && !all_zi {
			continue
		}
		if c == 0 {
			c = v
		} else if c>>4 != v>>4 {
			one_suit = false
		} else if uint32(v>>4) >= FENG {
			one_suit = false
		}
		if uint32(v>>4) < FENG {
			all_zi = false
		}
	}
	for _, v := range m {
		if v == 2 && !b {
			b = true
			continue
		}
		if v < 3 || chow || kong {
			seven_pair_big = false
			break
		}
	}
	if seven_pair_big {
		huType |= HU
		huType |= HU_SEVEN_PAIR_BIG
	}
	// (huType > 0 || hu) //有牌型
	if one_suit && (huType > 0 || hu) {
		huType |= HU
		huType |= HU_ONE_SUIT
	}
	if all_zi {
		huType |= HU
		huType |= HU_ALL_ZI
	}
	return huType
}

// 13烂(数牌相隔2个或以上,字牌不重复)/七星13烂(13烂基础上包含7个不同字牌)
// 有序slices cs && len(cs) == 14
func existThirteen(cs []uint32) uint32 {
	var thirteen_single bool = true
	var thirteen_single_zi bool = true
	var c uint32
	var n int
	var huType uint32
	for _, v := range cs {
		if uint32(v>>4) < FENG {
			if c == 0 || c>>4 != v>>4 {
				c = v
				continue
			}
			if v-c < 0x03 {
				thirteen_single = false
				break
			}
			c = v
		} else {
			n++
			if v == c {
				thirteen_single_zi = false
				break
			}
			c = v
		}
	}
	if thirteen_single_zi && thirteen_single_zi && n == 7 {
		huType |= HU_SINGLE_ZI
	} else if thirteen_single && thirteen_single_zi {
		huType |= HU_SINGLE
	}
	return huType
}

// 判断是否小七对(7个对子),豪华小七对(7个对子,其中有杠)
// 有序slices cs && len(cs) == 14
func exist7pair(cs []uint32) uint32 {
	var seven_pair bool = true
	var c uint32
	var i int
	var huType uint32
	var le int = len(cs)
	for n := 0; n < le-1; n += 2 {
		if cs[n] != cs[n+1] {
			seven_pair = false
			break
		}
		if i != 4 && cs[n] == c {
			i += 2
		} else if i != 4 {
			c = cs[n]
			i = 2
		}
	}
	if seven_pair && i == 4 {
		huType |= HU_SEVEN_PAIR_KONG
	} else if seven_pair {
		huType |= HU_SEVEN_PAIR
	}
	return huType
}

// -------广东-------

// 胡牌后检测
// 除十三幺，无序slices cs 包含碰，杠
func existHuTypeYao(cs []uint32) uint32 {
	var c_n uint32
	var suit_n bool = true
	var n = make(map[uint32]int)
	var m = make(map[uint32]int)
	for _, v := range cs {
		if type_zi(v) {
			n[v] += 1
			c_n, suit_n = same_suit(c_n, v, suit_n)
		} else {
			m[v] += 1
		}
	}
	var l_n int = len(n)
	var l_m int = len(m)
	//字一色(不再计算碰碰胡)
	if l_n == 0 {
		return HU_ALL_ZI
	}
	var pp bool = true //碰碰胡
	var yj bool = true //幺九
	for c, v := range n {
		if !one_nine(c) {
			yj = false
		}
		if v == 1 {
			pp = false
		}
	}
	//大哥(清一色碰碰胡)
	if l_m == 0 && suit_n && pp {
		return HU_ONE_SUIT_PONG
	}
	//纯幺九
	if l_m == 0 && yj {
		return HU_ONE_YAO
	}
	//混幺九
	if l_n != 0 && l_m != 0 && yj {
		return HU_TWO_YAO
	}
	//混一色碰碰胡
	if l_n != 0 && suit_n && pp {
		return HU_TWO_SUIT_PONG
	}
	//清一色
	if l_m == 0 && suit_n {
		return HU_ONE_SUIT
	}
	//碰碰胡
	if pp {
		return HU_PONGPONG
	}
	//混一色
	if l_n != 0 && suit_n {
		return HU_TWO_SUIT
	}
	return 0
}

// 同色
func same_suit(c1, c2 uint32, ok bool) (uint32, bool) {
	if !ok {
		return c2, ok
	}
	if c1 == 0 {
		return c2, ok
	}
	if c1>>4 != c2>>4 {
		return c2, false
	}
	return c2, ok
}

// 幺九
func one_nine(c uint32) bool {
	var v uint32 = c & 0xf
	if v == 1 || v == 9 {
		return true
	}
	return false
}

// 是否字牌
func type_zi(c uint32) bool {
	return (c >> 4) < 4
}

// 十三幺，有序slices cs && len(cs) == 14
func existThirteenYao(cs []uint32) uint32 {
	var n = make(map[uint32]int)
	var m = make(map[uint32]int)
	for _, v := range cs {
		if type_zi(v) {
			if !one_nine(v) {
				return 0
			}
			n[v] += 1
		} else {
			m[v] += 1
		}
	}
	if len(m) != 7 || len(n) != 6 {
		return 0
	}
	return HU_THIRTEEN_YAO
}

// 判断是否胡牌,0表示不胡牌,非0用32位表示不同的胡牌牌型
func existHu(cs []uint32) uint32 {
	le := len(cs)
	var huType uint32
	//单钓胡牌
	if le == 2 && cs[0] == cs[1] {
		huType |= HU
		huType |= HU_PING
		return huType
	}
	//排序slices
	Sort(cs, 0, le-1)
	// 14张牌型胡牌
	if le == 14 {
		// 十三幺牌型胡牌
		t2 := existThirteenYao(cs)
		if t2 > 0 {
			huType |= HU
			huType |= t2
			return huType
		}
	}
	// 3n +2 牌型胡牌
	//if (le - 2) % 3 != 0 {
	//	return huType
	//}
	return existHu3n2(cs, le)
}

// 3n +2 牌型胡牌
// 有序slices cs
func existHu3n2(cs []uint32, le int) (huType uint32) {
	list := make([]uint32, le)
	for n := 0; n < le-1; n++ {
		if cs[n] == cs[n+1] { //
			copy(list, cs)
			list[n] = 0x00
			list[n+1] = 0x00
			for i := 0; i < le-2; i++ {
				if list[i] > 0 {
					for j := i + 1; j < le-1; j++ {
						if list[j] > 0 && list[i] > 0 {
							for k := j + 1; k < le; k++ {
								if list[k] > 0 && list[i] > 0 && list[j] > 0 {
									//刻子
									if list[i] == list[j] && list[j] == list[k] {
										list[i], list[j], list[k] = 0x00, 0x00, 0x00
										break
									}
									if list[i] >= 0x41 { //广东字牌不成顺子
										break
									}
									//顺子
									if list[i]+1 == list[j] && list[j]+1 == list[k] {
										list[i], list[j], list[k] = 0x00, 0x00, 0x00
										break
									}
									//宜春东南西北任3张成顺 TODO fix bug
									//if list[i] >= 0x41 && list[k] <= 0x44 &&
									//list[i] != list[j] && list[j] != list[k] {
									//	list[i], list[j], list[k] = 0, 0, 0
									//	break
									//}
								}
							}
						}
					}
				}
			}
			num := false
			for i := 0; i < le; i++ {
				if list[i] > 0 {
					num = true
					break
				}
			}
			if !num {
				huType |= HU
				huType |= HU_PING
				return huType
			}
		}
	}
	return huType
}

//-------------赖子玩法---------

//算番(牌型) TODO:优化
func HuType2Tao(v uint32) int32 {
	switch {
	case v&HU_TWO_SUIT > 0:
		return FAN[HU_TWO_SUIT]
	case v&HU_PONGPONG > 0:
		return FAN[HU_PONGPONG]
	case v&HU_ONE_SUIT > 0:
		return FAN[HU_ONE_SUIT]
	case v&HU_TWO_SUIT_PONG > 0:
		return FAN[HU_TWO_SUIT_PONG]
	case v&HU_TWO_YAO > 0:
		return FAN[HU_TWO_YAO]
	case v&HU_ONE_YAO > 0:
		return FAN[HU_ONE_YAO]
	case v&HU_ONE_SUIT_PONG > 0:
		return FAN[HU_ONE_SUIT_PONG]
	case v&HU_THIRTEEN_YAO > 0:
		return FAN[HU_THIRTEEN_YAO]
	case v&HU_ALL_ZI > 0:
		return FAN[HU_ALL_ZI]
	case v&HU_PING > 0:
		return FAN[HU_PING]
	default:
		return 1
	}
}

//算番(胡牌方式) TODO:优化
func HuWay2Tao(v uint32) int32 {
	switch {
	case v&TIAN_HU > 0:
		return FAN[TIAN_HU]
	case v&DI_HU > 0:
		return FAN[DI_HU]
	case v&HU_KONG_FLOWER > 0:
		return FAN[HU_KONG_FLOWER]
	case v&REPAO > 0:
		return FAN[REPAO]
	case v&QIANG_GANG > 0:
		return FAN[QIANG_GANG]
	case v&HEIMO > 0:
		return FAN[HEIMO]
	case v&RUANMO > 0:
		return FAN[RUANMO]
	case v&PAOHU > 0:
		return FAN[PAOHU]
	case v&ZIMO > 0:
		return FAN[ZIMO]
	default:
		return 1
	}
}

//算番(杠)
func HuKongTao(v uint32) int32 {
	return FAN[v]
}

//玩法1：能接炮法如8为赖子
//123 44 789 11 12 可接胡
//123 44 567 11 12 可接胡
//183 44 567 11 12 不可接，只能摸
//
//第一个是赖子刚好是一句牌 叫黑摸
//第二个是没赖子 叫黑摸
//第三个是有赖子但跟第一种不同，所以只能摸 叫软摸

//玩法2: 不能接炮 8为赖子
//除了不能接炮与玩法1一样的

//黑摸倍数*2 扔赖子倍数*2 暗扛倍数*2 明扛倍数*1
//点扛=底分*1 算在放扛的人头上

// 判断是否胡牌,是否软胡,0表示不胡牌,非0用32位表示不同的胡牌牌型
func existHuTao(cards []uint32, laizi uint32, le int) (uint32, bool) {
	//排序slices
	Sort(cards, 0, le-1)
	// TODO 优化
	var n int
	cs1 := make([]uint32, 0)
	for i := 0; i < le; i++ {
		if cards[i] == laizi {
			n++
		} else {
			cs1 = append(cs1, cards[i])
		}
		if n > 1 { //多个不能胡
			return 0, false
		}
	}
	//TODO 其它牌型
	val := existHu3n2_normal(cards, le)
	if val > 0 { //黑摸
		return val, false
	}
	if n == 0 {
		return 0, false
	}
	val = existHu3n2_laizi(cs1, len(cs1))
	if val > 0 { //软摸
		return val, true
	}
	return 0, false
}

// 有序slices cards, 3n +2 牌型胡牌
func existHu3n2_normal(cards []uint32, le int) uint32 {
	cs := make([]uint32, le)
	for n := 0; n < le-1; n++ {
		if cards[n] == cards[n+1] { //
			copy(cs, cards)
			cs[n], cs[n+1] = 0, 0
			for i := 0; i < le-2; i++ {
				if cs[i] > 0 {
					for j := i + 1; j < le-1; j++ {
						if cs[j] > 0 && cs[i] > 0 {
							for k := j + 1; k < le; k++ {
								if cs[k] > 0 && cs[i] > 0 && cs[j] > 0 {
									//刻子
									if cs[i] == cs[j] && cs[j] == cs[k] {
										cs[i], cs[j], cs[k] = 0, 0, 0
										break
									}
									//顺子
									if cs[i]+1 == cs[j] && cs[j]+1 == cs[k] {
										cs[i], cs[j], cs[k] = 0, 0, 0
										break
									}
								}
							}
						}
					}
				}
			}
			if existHu_normal(cs, le) {
				return HU | HU_PING
			}
		}
	}
	return 0
}

//slice cs 有序, 是否可以胡
func existHu_normal(cs []uint32, le int) bool {
	for i := 0; i < le; i++ {
		if cs[i] > 0 {
			return false
		}
	}
	return true
}

// 有序slices cs, 3n +2 牌型胡牌
func existHu3n2_laizi(list []uint32, le int) uint32 {
	cs := make([]uint32, le)
	for n := 0; n < le; n++ {
		copy(cs, list)
		for i := n; i < le-2; i++ {
			if cs[i] > 0 {
				for j := i + 1; j < le-1; j++ {
					if cs[j] > 0 && cs[i] > 0 {
						for k := j + 1; k < le; k++ {
							if cs[k] > 0 && cs[i] > 0 && cs[j] > 0 {
								//刻子
								if cs[i] == cs[j] && cs[j] == cs[k] {
									cs[i], cs[j], cs[k] = 0, 0, 0
									break
								}
								//顺子
								if cs[i]+1 == cs[j] && cs[j]+1 == cs[k] {
									cs[i], cs[j], cs[k] = 0, 0, 0
									break
								}
							}
						}
					}
				}
			}
		}
		//是否可以胡
		if existHu_laizi(cs, le) {
			return HU | HU_PING
		}
	}
	return 0
}

//slice cs 有序, 是否可以胡
func existHu_laizi(cs []uint32, le int) bool {
	ls := make([]uint32, 0)
	for i := 0; i < le; i++ {
		if cs[i] > 0 {
			ls = append(ls, cs[i])
		}
	}
	ls_l := len(ls)
	if ls_l == 1 { //做将,手牌中有杠的情况也可以胡
		return true
	}
	if ls_l == 4 { //有一对将,另一对做顺或刻
		return hu_3n2(ls)
	}
	return false
}

//slice ls 有序,两个数同色且相差值是是否小于3(可以组成顺子或刻子)
func hu_3n2(ls []uint32) bool {
	if ls[0] == ls[1] {
		return hu_3n(ls[2], ls[3])
	}
	if ls[2] == ls[3] {
		return hu_3n(ls[0], ls[1])
	}
	if ls[1] == ls[2] {
		return hu_3n(ls[0], ls[3])
	}
	return false
}

//两个数同色且相差值是是否小于3(可以组成顺子或刻子)
func hu_3n(n, m uint32) bool {
	if n>>4 != m>>4 {
		return false
	}
	if n > m {
		return (n - m) < 3
	}
	return (m - n) < 3
}
