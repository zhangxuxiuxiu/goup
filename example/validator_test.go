package util

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"testing"
)

//https://pkg.go.dev/github.com/go-playground/validator#hdr-Minimum
//https://github.com/go-playground/validator/blob/v9/_examples/struct-level/main.go

type Request struct {
	Name     string `validate:"required,min=3,oneof=ni na"`
	Children []int  `validate:"unique,min=3"`
}

func TestValidate(t *testing.T) {
	valid := validator.New()
	if err := valid.Struct(Request{Name: "ni", Children: []int{2, 3, 3}}); err != nil {
		fmt.Printf("err:%s\n", err.Error())
	} else {
		fmt.Printf("valid struct\n")
	}
	if err := valid.Struct(Request{Name: "kax", Children: []int{2, 3}}); err != nil {
		fmt.Printf("err:%s\n", err.Error())
	} else {
		fmt.Printf("valid struct\n")
	}

}
