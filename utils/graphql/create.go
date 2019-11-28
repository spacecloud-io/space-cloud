package graphql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/schema"
	"github.com/spaceuptech/space-cloud/utils"
)

func (graph *Module) generateWriteReq(ctx context.Context, field *ast.Field, token string, store map[string]interface{}) ([]model.AllRequest, []interface{}, error) {
	dbType, err := GetDBType(field)
	if err != nil {
		return nil, nil, err
	}

	col := strings.TrimPrefix(field.Name.Value, "insert_")

	docs, err := extractDocs(field.Arguments, store)
	if err != nil {
		return nil, nil, err
	}

	reqs, returningDocs, err := graph.processNestedFields(docs, dbType, col)
	if err != nil {
		return nil, nil, err
	}

	// Check if the requests are authorised
	for _, req := range reqs {
		r := &model.CreateRequest{Document: req.Document, Operation: req.Operation}
		_, err = graph.auth.IsCreateOpAuthorised(ctx, graph.project, dbType, req.Col, token, r)
		if err != nil {
			return nil, nil, err
		}
	}

	return reqs, returningDocs, nil
}

func (graph *Module) prepareDocs(doc map[string]interface{}, schemaFields schema.SchemaFields) {
	// FieldIDs is the array of fields for which an unique id needs to be generated. These will only be done for those
	// fields which have the type ID.
	// FieldDates is the array of fields for which the current time needs to be set.
	fieldIDs := make([]string, 0)
	fieldDates := make([]string, 0)

	for fieldName, fieldSchema := range schemaFields {
		if fieldSchema.Kind == schema.TypeID {
			fieldIDs = append(fieldIDs, fieldName)
		}

		if fieldSchema.IsCreatedAt || fieldSchema.IsUpdatedAt {
			fieldDates = append(fieldDates, fieldName)
		}
	}

	// Set a new id for all those field ids which do not have the field set already
	for _, field := range fieldIDs {
		if _, p := doc[field]; !p {
			doc[field] = ksuid.New().String()
		}
	}

	// Set the current time for all fieldDates
	for _, field := range fieldDates {
		doc[field] = time.Now().UTC()
	}
}

func copyDoc(doc map[string]interface{}) map[string]interface{} {
	newDoc := make(map[string]interface{}, len(doc))
	for k, v := range doc {
		newDoc[k] = v
	}

	return newDoc
}

func (graph *Module) processNestedFields(docs []interface{}, dbType, col string) ([]model.AllRequest, []interface{}, error) {
	createRequests := make([]model.AllRequest, 0)
	afterRequests := make([]model.AllRequest, 0)

	// Check if we can the schema for this collection
	schemaFields, p := graph.auth.Schema.GetSchema(dbType, col)
	if !p {
		// Return the docs as is if no schema is available
		return []model.AllRequest{{Type: string(utils.Create), Col: col, Operation: utils.All, Document: docs}}, docs, nil
	}

	returningDocs := make([]interface{}, len(docs))

	for i, docTemp := range docs {

		// Each document is actually an object
		doc := docTemp.(map[string]interface{})
		graph.prepareDocs(doc, schemaFields)

		newDoc := copyDoc(doc)

		// Iterate over each field of the document to see if has any linked fields that are present
		for fieldName, fieldValue := range doc {
			fieldSchema, p := schemaFields[fieldName]
			if !p || !fieldSchema.IsLinked {
				// Simply ignore if the field does not have a corresponding schemaFields or it isn't linked
				continue
			}

			// We are here means that the field is actually a linked value

			fromFieldSchema, p := schemaFields[fieldSchema.LinkedTable.From]
			if !p {
				// Ignore if the `from` key is not present in the schema
				continue
			}

			// Ignore if the field key is present in the linked table config. We don't support such operations.
			if fieldSchema.LinkedTable.Field != "" {
				continue
			}

			if fromFieldSchema.IsPrimary {
				// Ignore if the from field doesn't exist in the document
				if _, p := doc[fromFieldSchema.FieldName]; !p {
					continue
				}
			}

			// Lets populate an array of linked docs
			var linkedDocs []interface{}
			if !fieldSchema.IsList {
				temp, ok := fieldValue.(map[string]interface{})
				if !ok {
					return nil, nil, fmt.Errorf("invalid format provided for linked field %s - wanted object got array", fieldName)
				}

				linkedDocs = []interface{}{temp}
			} else {
				temp, ok := fieldValue.([]interface{})
				if !ok {
					return nil, nil, fmt.Errorf("invalid format provided for linked field %s - wanted array got object", fieldName)
				}
				linkedDocs = temp
			}

			// Iterate over each linked doc
			for _, linkedDocTemp := range linkedDocs {
				// Each document is actually an object
				linkedDoc := linkedDocTemp.(map[string]interface{})

				graph.prepareDocs(linkedDoc, schemaFields)

				// Check if the `from` field is a primary key. If it is, we need to set that value in the `to` field
				// of the nested value. If it is not a primary key, we'll have to set it with the value of the `to`
				// field of the nested value
				if fromFieldSchema.IsPrimary {
					linkedDoc[fieldSchema.LinkedTable.To] = doc[fieldSchema.LinkedTable.From]
				} else {
					// The nested docs need to be inserted first in this case
					doc[fieldSchema.LinkedTable.From] = linkedDoc[fieldSchema.LinkedTable.To]
				}
			}

			linkedCreateRequests, returningLinkedDocs, err := graph.processNestedFields(linkedDocs, dbType, fieldSchema.LinkedTable.Table)
			if err != nil {
				return nil, nil, err
			}

			if fromFieldSchema.IsPrimary {
				// It the from field is primary, it means that the nested docs need to be inserted after the parent docs have been inserted
				afterRequests = append(afterRequests, linkedCreateRequests...)
			} else {
				// If the from field is not primary, it means that the nested docs need to be inserted before the parent docs
				createRequests = append(createRequests, linkedCreateRequests...)
			}

			// Delete the nested field. The schema module would throw an error otherwise
			delete(doc, fieldName)

			newDoc[fieldName] = returningLinkedDocs
		}
		returningDocs[i] = newDoc
	}
	createRequests = append(createRequests, model.AllRequest{Type: string(utils.Create), Col: col, Operation: utils.All, Document: docs})
	return append(createRequests, afterRequests...), returningDocs, nil
}

func extractDocs(args []*ast.Argument, store utils.M) ([]interface{}, error) {
	for _, v := range args {
		switch v.Name.Value {
		case "docs":
			temp, err := ParseValue(v.Value, store)
			if err != nil {
				return nil, err
			}
			return temp.([]interface{}), nil
		}
	}

	return []interface{}{}, nil
}
