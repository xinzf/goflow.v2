package functions

import (
	"errors"
	"github.com/go-xorm/xorm"
	"github.com/xinzf/goflow.v2/spi"
	"github.com/xinzf/goflow.v2/tools"
	"gitlab.litudai.com/worker/server/models"
)

type CopyTo struct {
}

func (this *CopyTo) GetName() string {
	return "copyTo"
}

func (this *CopyTo) Eval(store spi.Store, transientVars *tools.TransientVars, args map[string]*tools.Value) (interface{}, error) {

	users, ok := args["users"]
	if !ok {
		return nil, errors.New("function:copyTo 缺少参数：users")
	}

	userIds := users.IntSlice(",")
	if len(userIds) == 0 {
		return nil, errors.New("function:copyTo 缺少参数：users")
	}

	entry := transientVars.Get(tools.Entry).GetData().(spi.Entry)
	copies := make([]*models.Copy, 0)
	err := store.GetConn().(*xorm.Session).Where("entry_id=?", entry.GetEntryId()).In("copy_to", userIds).Find(&copies)
	if err != nil {
		return nil, err
	}

	uids := make([]int, 0)
	for _, id := range userIds {
		exist := false
		for _, c := range copies {
			if c.CopyTo == id {
				exist = true
				break
			}
		}

		if !exist {
			uids = append(uids, id)
		}
	}

	if len(uids) > 0 {
		err = store.CreateCopy(entry.GetEntryId(), uids)
	}
	return nil, err
}
