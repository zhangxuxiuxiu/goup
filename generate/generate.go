package gen

// run "go generate" in current dir
//
//go:generate go-enum -type=Gender
type Gender uint8

const (
	Unknown Gender = iota + 1
	Male
	Female
)

//go:generate schematyper -o schema_type.go --package gen ./schematyper.json

//https://github.com/globusdigital/deep-copy
//[--skip Selector1,Selector.Two --skip Selector2[i], Selector.Three[k]]
//go:generate deep-copy --type Schematyper --skip ID,MultipleOf --pointer-receiver -o ./copy_schema_typer.go .

//go:generate wire gen  -output_file_prefix=eg_
//go:generate sed -i "" -e "/go:generate/d" eg_wire_gen.go

//go:generate impl "s *Source" golang.org/x/oauth2.TokenSource
