package memory_adapter

import (
	"context"
	"encoding/json"
	"fmt"
)

type Cache interface {
	Set(ctx context.Context, key, value []byte) error
	Get(ctx context.Context, key []byte) ([]byte, error)
}

type CacheCenter struct {
	BizCache Cache
}

func (biz *CacheCenter) SetAny(ctx context.Context, key string, val any) error {
	body, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return biz.BizCache.Set(ctx, []byte(key), body)
}

func (biz *CacheCenter) GetAny(ctx context.Context, key string, val any) error {
	body, err := biz.BizCache.Get(ctx, []byte(key))
	if err != nil {
		return err
	}
	return json.Unmarshal(body, val)
}

var cacheCenter CacheCenter

func GetCacheCenter() *CacheCenter {
	return &cacheCenter
}

type ProductBiz struct{}

func (biz *ProductBiz) ReloadProducts(ctx context.Context) error {
	var products []Product // get from DB or redis
	return GetCacheCenter().SetAny(ctx, "product", products)
}

func (biz *ProductBiz) GetProducts(ctx context.Context) ([]Product, error) {
	var products []Product
	if err := GetCacheCenter().GetAny(ctx, "product", products); err != nil {
		return nil, err
	}
	return products, nil
}

func (biz *ProductBiz) GetProductById(ctx context.Context, productId uint64) (*Product, error) {
	products, err := biz.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	for _, pro := range products {
		if pro.ProductId == productId {
			return &pro, nil
		}
	}
	return nil, fmt.Errorf("product with id:%d not found", productId)
}
