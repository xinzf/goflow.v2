package goflow

import (
	"github.com/xinzf/goflow.v2/functions"
	"github.com/xinzf/goflow.v2/spi"
)

func NewPropertSet(store spi.Store) *PropertSet {
	c := &PropertSet{
		store: store,
		funcs: make(map[string]func() spi.Function),
	}

	c.RegisterFunctions(
		func() spi.Function {
			return new(functions.CanJoin)
		},
	)

	return c
}

type PropertSet struct {
	store spi.Store
	funcs map[string]func() spi.Function
}

func (this *PropertSet) GetStore() spi.Store {
	return this.store
}

func (this *PropertSet) RegisterFunctions(funs ...func() spi.Function) {
	for _, fun := range funs {
		temp := fun()
		this.funcs[temp.GetName()] = fun
	}
}

func (this *PropertSet) GetFunction(name string) (spi.Function, bool) {
	if fn, found := this.funcs[name]; !found {
		return nil, false
	} else {
		return fn(), true
	}
}
