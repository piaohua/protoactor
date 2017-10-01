package data

import (
	"log/syslog"
	"os"
	"protoactor/entity"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine
var Engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		panic(err)
	}

	engine.ShowSQL(true) //在控制台打印出生成的SQL语句
	//日志默认显示级别为INFO
	engine.Logger().SetLevel(core.LOG_DEBUG) //在控制台打印调试及以上的信息

	sql2log() //日志写入文件

	//engine.SetMaxIdleConns(5) //设置连接池的空闲数大小,default is 2
	//engine.SetMaxOpenConns(5) //设置最大打开连接数

	//设置时区
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

	//同步结构
	err = engine.Sync2(new(entity.Account), new(entity.User),
		new(entity.LogRegist), new(entity.LogLogin),
		new(entity.LogDiamond), new(entity.LogCoin))
	if err != nil {
		panic(err)
	}
	Engine = engine
}

//err := engine.Ping() //测试是否可以连接到数据库
//数据库有连接超时设置的,可以通过起一个定期Ping的Go程来保持连接鲜活
func Clone() *xorm.Engine {
	newEngine, err := engine.Clone()
	if err != nil {
		panic(err)
	}
	return newEngine
}

//将信息不仅打印到控制台，而是保存为文件
func sql2log() {
	f, err := os.Create("sql.log")
	if err != nil {
		//println(err.Error())
		//return
		panic(err)
	}
	engine.SetLogger(xorm.NewSimpleLogger(f))
}

//将日志记录到syslog中
func sql2syslog() {
	logWriter, err := syslog.New(syslog.LOG_DEBUG, "rest-xorm-example")
	if err != nil {
		//log.Fatalf("Fail to create xorm system logger: %v\n", err)
		panic(err)
	}
	logger := xorm.NewSimpleLogger(logWriter)
	logger.ShowSQL(true)
	engine.SetLogger(logger)
}

/*
名称映射规则
xorm内置了三种IMapper实现:
core.SnakeMapper
core.SameMapper
core.GonicMapper
*SnakeMapper 支持struct为驼峰式命名,表结构为下划线命名之间的转换,这个是默认的Maper;
*SameMapper  支持结构体名称和对应的表名称以及结构体field名称与对应的表字段名称相同的命名；
*GonicMapper 和SnakeMapper很类似,但是对于特定词支持更好,比如ID会翻译成id而不是i_d。
当前SnakeMapper为默认值，如果需要改变时，在engine创建完成后使用
engine.SetMapper(core.SameMapper{})
*/
