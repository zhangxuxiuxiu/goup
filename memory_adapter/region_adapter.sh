#type ProductBiz interface {
#        RequireById(ctx context.Context, productId int64) *model.Product
#
#}
# generate following:
#type RegionProductBizAdapter struct{
#	regions map[string]*ProductBiz
#}
#
#func (impl *RegionProductBizAdapter) RequireById(ctx context.Context, productId int64) *model.Product{
#	return impl.regions[regionOf(ctx)].RequiredById(ctx, productId)
#}
#

sed  '/[[:blank:]]*type[[:blank:]]\{1,\}\([[:alnum:]]\{1,\}\)/ {
h
# type ProductBiz interface {
s/[[:blank:]]*type[[:blank:]]\{1,\}\([[:alnum:]]\{1,\}\).*/type Region\1Adapter struct { \n regions map[string]\*\1 \n}\n/
:x
n
/}[[:blank:]]*$/{
s/.*//
q
}
/^\s*$\|^\s*\/\//!{
G
# GetByIds(ctx context.Context, productIds []int64) []*model.Product
# type ProductBiz interface {
s/^[[:blank:]]*\(\([[:alnum:]_]\{1,\}\)(\([^)]\{1,\}\)).*\)\n[[:blank:]]*type[[:blank:]]\{1,\}\([[:alnum:]]\{1,\}\).*/func (impl *Region\4Adapter) \1{ \n\treturn impl.regions[regionOf(ctx)]\.\2(\3) \n}/M
}
bx
}' ../../insurance-product/src/product_factory/biz/product_biz.go  

awk '{ 
if( match($0, "\\s*type\\s+(\\w+)\\s+interface", structName) > 0 ) {
	printf("type Region%sAdapter struct{\n\tregions map[string]*%s\n}\n\n", structName[1], structName[1])
	while( !match($0,"}\\s*$")){
		if (getline <= 0) {
                	print("unexpected EOF or error:", ERRNO) > "/dev/stderr"
                	exit
            	}
		if ( match($0, "^\\s*(\\w+)\\((.*)\\)", fnArgs) > 0 ){
				fnName = fnArgs[1]
				# \w does not work in [], period(ie ".") requires no escape in []
				vals = gensub("(\\w+)\\s+[[:alnum:]_.]+\\s*(,?)", "\\1\\2","g",fnArgs[2])
				printf("func (impl *Region%sAdapter)%s(%s){\n\treturn impl.regions[regionOf(ctx)].%s(%s)\n}\n\n", structName[1],fnName,fnArgs[2],fnName,vals)
		}	
		delete fnArgs
	}
	exit
} else {
	print $0 "\n"
}
}' ../../insurance-product/src/product_factory/biz/product_biz.go  

