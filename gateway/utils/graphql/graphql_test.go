package graphql_test

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils/graphql"
)

type tests struct {
	name             string
	args             args
	crudMockArgs     []mockArgs
	functionMockArgs []mockArgs
	authMockArgs     []mockArgs
	schemaMockArgs   []mockArgs
	wantResult       interface{}
	wantErr          bool
}
type mockArgs struct {
	method         string
	args           []interface{}
	paramsReturned []interface{}
}
type args struct {
	req   *model.GraphQLRequest
	token string
}

func TestModule_ExecGraphQLQuery(t *testing.T) {
	tests := make([]tests, 0)
	tests = append(tests, queryTestCases...)
	tests = append(tests, mutationTestCases...)
	tests = append(tests, upsertTestCases...)
	tests = append(tests, deleteTestCases...)
	tests = append(tests, transactionTestCases...)
	tests = append(tests, functionTestCases...)
	tests = append(tests, prepareQueryTestCases...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCrud := mockGraphQLCrudInterface{}
			for _, m := range tt.crudMockArgs {
				mockCrud.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			mockAuth := mockGraphQLAuthInterface{}
			for _, m := range tt.authMockArgs {
				mockAuth.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			mockFunction := mockGraphQLFunctionInterface{}
			for _, m := range tt.functionMockArgs {
				mockFunction.On(m.method, m.args...).Return(m.paramsReturned...)
			}
			mockSchema := mockGraphQLSchemaInterface{}
			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			graph := graphql.New(&mockAuth, &mockCrud, &mockFunction, &mockSchema)

			ctx := context.Background()

			var wg sync.WaitGroup
			var testOp interface{}
			var testErr error
			wg.Add(1)
			graph.ExecGraphQLQuery(ctx, tt.args.req, tt.args.token, func(op interface{}, err error) {
				defer wg.Done()
				testOp = op
				testErr = err
			})
			wg.Wait()
			if (testErr != nil) != tt.wantErr {
				t.Errorf("ExecGraphQLQuery() got error %v want error %v", testErr, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.wantResult, testOp) {
				t.Errorf("ExecGraphQLQuery() got result %v want result %v", testOp, tt.wantResult)
			}

			if !tt.wantErr {
				mockCrud.AssertExpectations(t)
				mockSchema.AssertExpectations(t)
				mockFunction.AssertExpectations(t)
				mockAuth.AssertExpectations(t)
			}
		})
	}
}
