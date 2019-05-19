package main

import (
	"github.com/xinzf/goflow.v2/tools"
	"log"
)

//func main() {
//	err := goflow.Start("/Users/xiangzhi/Work/Go/src/github.com/xinzf/goflow.v2/test/config")
//	if err != nil {
//		logrus.Errorln(err)
//	} else {
//		w, found := goflow.Get("bldz-20190516")
//		if !found {
//			logrus.Errorln("没有找到")
//		} else {
//			print(w.ToJson())
//		}
//	}
//}

type Tester interface {
	GetName() string
}
type Test struct {
	name string
}

func (this *Test) GetName() string {
	return this.name
}

func main() {
	trans := tools.NewTransientVars()
	trans.Put("test", &Test{name: "xiangzhi"})

	var t Tester
	t = trans.Get("test").GetData().(Tester)
	log.Println(t)
}
