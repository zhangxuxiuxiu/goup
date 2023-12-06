package memory_adapter

import (
	"context"
)

type BizImpl struct {
	//	products map[int64]*Product
	products map[string]map[int64]*Product
}

func (impl *BizImpl) Get(ctxt context.Context, productId int64) *Product {
	//	return impl.products[productId]
	return impl.products[regionOf(ctxt)][productId]
}

func (impl *BizImpl) Set(ctxt context.Context, productId int64, pro *Product) {
	//	impl.products[productId] = pro
	impl.products[regionOf(ctxt)][productId] = pro
}

func regionOf(ctxt context.Context) string {
	return ctxt.Value("region").(string)
}
