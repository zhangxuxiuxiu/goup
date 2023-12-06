package main

import (
	"context"
	"fmt"
	"github.com/go_practice/config_insurance/dbops"
	"github.com/go_practice/config_insurance/domain"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"time"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) < 3 {
		showUsage()
		os.Exit(-1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "load":
		loadProduct()
	case "insert":
		fmt.Printf(dbops.Insert(unmarshalProduct(os.Args[2]), "product_tab"))
	case "delete":
		fmt.Printf(dbops.Delete(unmarshalProduct(os.Args[2]), "product_tab"))
	case "update":
		fmt.Printf("%s", dbops.Update(unmarshalProduct(os.Args[2]), unmarshalProduct(os.Args[3]), "product_tab"))
	default:
		fmt.Fprintf(os.Stderr, "invalid cmd:%s\n", cmd)
		showUsage()
	}
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage:\n"+
		"\t%s load product_id dev|test|uat id|vn\n"+
		"\t%s insert product.yml\n"+
		"\t%s delete product.yml\n"+
		"\t%s update v1.yml v2.yml\n",
		os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func loadProduct() {
	productId, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid product_id:%d", os.Args[1])
		os.Exit(-2)
	}

	region := "id"
	env := "test"
	if len(os.Args) > 3 {
		region = os.Args[3]
	}
	if len(os.Args) > 4 {
		env = os.Args[4]
	}
	marketEngine, err := xorm.NewEngineGroup("mysql",
		[]string{os.Getenv(fmt.Sprintf("ipc_mysql_%s_%s", region, env))})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in creating xorm engine with error:%s", err.Error())
		os.Exit(-3)
	}
	defer marketEngine.Close()

	marketEngine.ShowSQL(false)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := marketEngine.PingContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error in ping mysql with error:%s", err.Error())
		os.Exit(-4)
	}

	session := marketEngine.NewSession()
	defer session.Close()
	product := domain.BuildProduct(productId, session)
	if bytes, err := yaml.Marshal(product); err != nil {
		fmt.Fprintf(os.Stderr, "error in yaml marshal product with error:%s", err.Error())
		os.Exit(-5)
	} else {
		fmt.Printf("%s", bytes)
	}
}

func unmarshalProduct(file string) *domain.Product {
	bytes, err := os.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("file:%s not found\n", file))
	}
	var product domain.Product
	if err := yaml.Unmarshal(bytes, &product); err != nil {
		panic(fmt.Sprintf("Product Unmarshal failed with:%s\n", err.Error()))
	}
	return &product
}
