package main

import (
	"fmt"
	"io/ioutil"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/json"
	"cuelang.org/go/encoding/jsonschema"
)

func main() {
	schema, err := ioutil.ReadFile("compose-spec.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonSchema, err := json.Extract("", schema)
	if err != nil {
		fmt.Println(err)
		return
	}
	cueCtx := cuecontext.New()

	cueJsonSchemaExpr := cueCtx.BuildExpr(jsonSchema)
	if cueJsonSchemaExpr.Err() != nil {
		fmt.Println(cueJsonSchemaExpr.Err())
		return
	}

	extractedSchema, err := jsonschema.Extract(cueJsonSchemaExpr, &jsonschema.Config{
		PkgName: "myPkg",
		// Strict:  true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	val := cueCtx.BuildFile(extractedSchema)
	if val.Err() != nil {
		fmt.Println(val.Err())
		return
	}
	fmt.Println(val)
}
