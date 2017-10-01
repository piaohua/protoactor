package robots

var robotmap map[string]RobotData

type RobotData struct {
	Id        uint32 `csv:"id"`
	Kind      uint32 `csv:"kind"`
	Changci   uint32 `csv:"changci"`
	Nickname  string `csv:"nickname"`
	Sex       uint32 `csv:"sex"`
	Level     uint32 `csv:"level"`
	Coin      uint32 `csv:"coin"`
	Diamond   uint32 `csv:"diamond"`
	Cointime  uint32 `csv:"cointime"`
	Coinup    uint32 `csv:"coinup"`
	Vip       uint32 `csv:"vip"`
	Headframe uint32 `csv:"headframe"`
	Phone     string `csv:"phone"`
}

func GetRobotList() []RobotData {
	return robotList
}

var robotList []RobotData

func GetRobot(phone string) *RobotData {
	if v, ok := robotmap[phone]; ok {
		return &v
	} else {
		return nil
	}

}

func init() {
	// robot.csv UTF-8 格式读取的v.Id为0,ANSI格式读取正常
	//f, err := os.Open("./csv/robot.csv")
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()

	//data, err := ioutil.ReadAll(f)
	//if err != nil {
	//	panic(err)
	//}
	//err = csv.Unmarshal(data, &robotList)
	//if err != nil {
	//	panic(err)
	//}
	//robotmap = make(map[string]RobotData)
	//for _, v := range robotList {
	//	robotmap[v.Phone] = v
	//	//glog.Infoln(v.Phone)
	//}
	//
	robotmap = make(map[string]RobotData)
	robotmap["13714592946"] = RobotData{Nickname: "test1"}
	robotmap["13714992949"] = RobotData{Nickname: "test2"}
	robotmap["13714592741"] = RobotData{Nickname: "test3"}
	robotmap["13714592480"] = RobotData{Nickname: "test4"}
	robotmap["10005000130"] = RobotData{Nickname: "孤酒浪人"}
	robotmap["10005000131"] = RobotData{Nickname: "我心已冷"}
	robotmap["10005000132"] = RobotData{Nickname: "百毒不侵"}
	robotmap["10005000133"] = RobotData{Nickname: "浪人与酒"}
	robotmap["10005000134"] = RobotData{Nickname: "世态炎凉"}
	robotmap["10005000135"] = RobotData{Nickname: "像个笑话"}
	robotmap["10005000136"] = RobotData{Nickname: "旧时猫巷"}
	robotmap["10005000137"] = RobotData{Nickname: "痴人说梦"}
	robotmap["10005000138"] = RobotData{Nickname: "孤独情话"}
	robotmap["10005000139"] = RobotData{Nickname: "败给现实"}
	robotmap["10005000140"] = RobotData{Nickname: "故人老街"}
	robotmap["10005000141"] = RobotData{Nickname: "神经领袖"}
	robotmap["10005000142"] = RobotData{Nickname: "顾与南歌"}
	robotmap["10005000143"] = RobotData{Nickname: "凉城凉忆"}
	robotmap["10005000144"] = RobotData{Nickname: "静梦蔷薇"}
	robotmap["10005000145"] = RobotData{Nickname: "盛夏物语"}
	robotmap["10005000146"] = RobotData{Nickname: "孤城少女"}
	robotmap["10005000147"] = RobotData{Nickname: "泪雨成殇"}
	robotmap["10005000148"] = RobotData{Nickname: "人心难测"}
	robotmap["10005000149"] = RobotData{Nickname: "孤身一人"}
	for _, v := range robotmap {
		robotList = append(robotList, v)
	}
}
