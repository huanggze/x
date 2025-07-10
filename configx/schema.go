package configx

import (
	"bytes"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/huanggze/x/logrusx"
	"github.com/huanggze/x/otelx"
	"github.com/ory/jsonschema/v3"
)

func newCompiler(schema []byte) (string, *jsonschema.Compiler, error) {
	id := gjson.GetBytes(schema, "$id").String()
	if id == "" {
		id = fmt.Sprintf("%s.json", uuid.Must(uuid.NewV4()).String())
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(id, bytes.NewBuffer(schema)); err != nil {
		return "", nil, errors.WithStack(err)
	}

	// DO NOT REMOVE THIS
	compiler.ExtractAnnotations = true

	if err := otelx.AddConfigSchema(compiler); err != nil {
		return "", nil, err
	}
	if err := logrusx.AddConfigSchema(compiler); err != nil {
		return "", nil, err
	}

	return id, compiler, nil
}
