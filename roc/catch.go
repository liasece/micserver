package roc

import (
	"github.com/liasece/micserver/util"
)

const (
	CATCH_POOL_GROUP_SUM = 8
)

type Catch struct {
	util.MapPool
}

func (this *Catch) Init() {
	this.MapPool.Init(CATCH_POOL_GROUP_SUM)
}
