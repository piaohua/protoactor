package messages

import (
	"reflect"
	"testing"

	proto "github.com/gogo/protobuf/proto"
)

// 启动一个服务
func TestProto(t *testing.T) {
	packet := &Request{}
	t.Logf("packet -> %+v", packet.String())
	s := proto.MessageName(&Request{})
	t.Logf("s -> %s", s)
	//---
	//p := proto.MessageType(s)
	p := reflect.TypeOf(Request{})
	t.Logf("p -> %v", p) //reflect.Type
	v := reflect.New(p)
	t.Logf("v -> %v", v)
	t.Logf("v -> %v", v.Interface())
	//---
	reg := &Request{
		UserName: "xxxx",
	}
	b, err := reg.Marshal()
	t.Log("err -> ", err)
	t.Logf("b -> %#v", b)
	say := &Request{}
	t.Logf("say -> %#v", say)
	err = say.Unmarshal(b)
	t.Log("err -> ", err)
	t.Logf("say -> %#v", say)
	//---
	ss := proto.MessageName(&Request{})
	mt := proto.MessageType(ss)
	mv := reflect.New(mt.Elem())
	t.Logf("mv -> %#v", mv)
	mb, err := proto.Marshal(mv.Interface().(proto.Message))
	t.Log("err -> ", err)
	t.Logf("mb -> %#v", mb)

	msg := &SRegist{
		Userid: "xxx",
	}
	b2, err := msg.Marshal()
	t.Log(b2, err)
}
