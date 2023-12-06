package memory_adapter

import "context"

type Product struct {
	ProductId uint64
	Name      string
}

type Biz interface {
	Get(ctxt context.Context, productId int64) *Product
	Set(ctxt context.Context, productId int64, pro *Product)
}

var productBiz Biz

func RefProductBiz() Biz {
	return productBiz
}

func InjectProductBiz(impl Biz) {
	productBiz = impl
}
