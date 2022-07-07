package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/json"
	"cuelang.org/go/encoding/jsonschema"
)

func main() {
	schema, err := readRemoteFile("https://raw.githubusercontent.com/compose-spec/compose-spec/master/schema/compose-spec.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonSchema, err := json.Extract("", []byte(schema))
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
		Strict:  true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(extractedSchema)
}

func readRemoteFile(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return " ", err
	}
	if response.StatusCode == http.StatusNotFound {
		return " ", fmt.Errorf("Not found")
	}

	defer response.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, response.Body)
	if err != nil {
		return " ", err
	}

	return buf.String(), nil
}
