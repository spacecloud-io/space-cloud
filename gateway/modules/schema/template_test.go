package schema

import (
	"strings"
	"testing"

	"github.com/go-test/deep"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func Test_generateSDL(t *testing.T) {
	type args struct {
		schemaCol model.Collection
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Successful test",
			args: args{
				schemaCol: model.Collection{"table1": model.Fields{
					"col2": &model.FieldType{
						FieldName:       "col2",
						IsAutoIncrement: false,
						Kind:            model.TypeID,
						TypeIDSize:      model.SQLTypeIDSize,
						IsPrimary:       true,
					},
					"col3": &model.FieldType{
						FieldName:       "col2",
						Kind:            model.TypeInteger,
						TypeIDSize:      model.SQLTypeIDSize,
						IsPrimary:       true,
						IsAutoIncrement: true,
					},
					"age": &model.FieldType{
						FieldName:       "age",
						IsAutoIncrement: false,
						Kind:            model.TypeFloat,
					},
					"createdAt": &model.FieldType{
						FieldName:       "createdAt",
						IsAutoIncrement: false,
						Kind:            model.TypeDateTime,
						IsCreatedAt:     true,
					},
					"updatedAt": &model.FieldType{
						FieldName:       "updatedAt",
						IsAutoIncrement: false,
						Kind:            model.TypeDateTime,
						IsUpdatedAt:     true,
					},
					"role": &model.FieldType{
						FieldName:           "role",
						IsAutoIncrement:     false,
						IsFieldTypeRequired: true,
						Kind:                model.TypeID,
						TypeIDSize:          model.SQLTypeIDSize,
						IsDefault:           true,
						Default:             "user",
					},
					"spec": &model.FieldType{
						FieldName:       "spec",
						IsAutoIncrement: false,
						Kind:            model.TypeJSON,
					},
					"first_name": &model.FieldType{
						FieldName:           "first_name",
						IsAutoIncrement:     false,
						IsFieldTypeRequired: true,
						Kind:                model.TypeID,
						TypeIDSize:          model.SQLTypeIDSize,
						IsIndex:             true,
						IndexInfo: &model.TableProperties{
							Group: "user_name",
							Order: 1,
							Sort:  "asc",
						},
					},
					"name": &model.FieldType{
						FieldName:           "name",
						IsAutoIncrement:     false,
						IsFieldTypeRequired: true,
						Kind:                model.TypeID,
						TypeIDSize:          model.SQLTypeIDSize,
						IsIndex:             true,
						IsUnique:            true,
						IndexInfo: &model.TableProperties{
							Group: "user_name",
							Order: 1,
							Sort:  "asc",
						},
					},
					"customer_id": &model.FieldType{
						FieldName:           "customer_id",
						IsAutoIncrement:     false,
						IsFieldTypeRequired: true,
						Kind:                model.TypeID,
						TypeIDSize:          model.SQLTypeIDSize,
						IsForeign:           true,
						JointTable: &model.TableProperties{
							To:             "id",
							Table:          "customer",
							OnDelete:       "CASCADE",
							ConstraintName: "c_tweet_customer_id",
						},
					},
					"order_dates": &model.FieldType{
						FieldName:       "order_dates",
						IsAutoIncrement: false,
						IsList:          true,
						Kind:            model.TypeDateTime,
						IsLinked:        true,
						LinkedTable: &model.TableProperties{
							Table:  "order",
							From:   "id",
							To:     "customer_id",
							Field:  "order_date",
							DBType: "mongo",
						},
					},
				},
				},
			},
			want:    "type  table1 { \n\tage: Float        \n\tcol2: ID @primary @size(value: 50)      \n\tcol3: Integer @primary(autoIncrement:true)      \n\tcreatedAt: DateTime  @createdAt      \n\tcustomer_id: ID!  @size(value: 50)       @foreign(table: customer, field: id ,onDelete: cascade)\n\tfirst_name: ID!  @size(value: 50)    @index(group: \"user_name\", sort: \"asc\", order: 1)   \n\tname: ID!  @size(value: 50)   @unique(group: \"user_name\", order: 1)     \n\torder_dates: DateTime       @link(table: order, from: id, to: customer_id, field: order_date) \n\trole: ID!  @size(value: 50)     @default(value: user)  \n\tspec: JSON        \n\tupdatedAt: DateTime   @updatedAt     \n}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSDL(tt.args.schemaCol)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSDL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// minify string by removing space

			if arr := deep.Equal(strings.Replace(got, " ", "", -1), strings.Replace(tt.want, " ", "", -1)); len(arr) > 0 {
				t.Errorf("generateSDL() differences = %v", arr)
			}
		})
	}
}
