package memory_adapter

import "context"

type RegionBizImpl struct {
	regionImpls map[string]*BizImpl
}

func (impl *RegionBizImpl) Get(ctxt context.Context, productId int64) *Product {
	return impl.regionImpls[regionOf(ctxt)].Get(ctxt, productId)
}

func (impl *RegionBizImpl) Set(ctxt context.Context, productId int64, pro *Product) {
	impl.regionImpls[regionOf(ctxt)].Set(ctxt, productId, pro)
}

func init() {
	InjectProductBiz(&RegionBizImpl{})
}
