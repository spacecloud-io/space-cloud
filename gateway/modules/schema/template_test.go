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
		{
			name: "Successful test",
			args: args{
				schemaCol: model.Collection{"table1": model.Fields{
					"col2": &model.FieldType{
						FieldName:  "col2",
						Kind:       model.TypeID,
						TypeIDSize: model.DefaultCharacterSize,
						IsPrimary:  true,
					},
					"col3": &model.FieldType{
						FieldName:       "col3",
						Kind:            model.TypeInteger,
						TypeIDSize:      model.DefaultCharacterSize,
						IsPrimary:       true,
						IsAutoIncrement: true,
						PrimaryKeyInfo: &model.TableProperties{
							Order: 2,
						},
					},
					"col4": &model.FieldType{
						FieldName:  "col4",
						Kind:       model.TypeInteger,
						TypeIDSize: model.DefaultCharacterSize,
						IsPrimary:  true,
						PrimaryKeyInfo: &model.TableProperties{
							Order: 1,
						},
					},
					"col5": &model.FieldType{
						FieldName:  "col5",
						Kind:       model.TypeInteger,
						TypeIDSize: model.DefaultCharacterSize,
						IsPrimary:  true,
					},
					"amount": &model.FieldType{
						FieldName:  "amount",
						Kind:       model.TypeBigInteger,
						TypeIDSize: model.DefaultCharacterSize,
					},
					"coolDownInterval": &model.FieldType{
						FieldName:  "coolDownInterval",
						Kind:       model.TypeSmallInteger,
						TypeIDSize: model.DefaultCharacterSize,
					},
					"age": &model.FieldType{
						FieldName: "age",
						Kind:      model.TypeFloat,
					},
					"createdAt": &model.FieldType{
						FieldName:   "createdAt",
						Kind:        model.TypeDateTime,
						IsCreatedAt: true,
					},
					"updatedAt": &model.FieldType{
						FieldName:   "updatedAt",
						Kind:        model.TypeDateTime,
						IsUpdatedAt: true,
					},
					"role": &model.FieldType{
						FieldName:           "role",
						IsFieldTypeRequired: true,
						Kind:                model.TypeID,
						TypeIDSize:          model.DefaultCharacterSize,
						IsDefault:           true,
						Default:             "user",
					},
					"spec": &model.FieldType{
						FieldName: "spec",
						Kind:      model.TypeJSON,
					},
					"first_name": &model.FieldType{
						FieldName:           "first_name",
						IsFieldTypeRequired: true,
						Kind:                model.TypeID,
						TypeIDSize:          model.DefaultCharacterSize,
						IndexInfo: []*model.TableProperties{{
							IsIndex: true,
							Group:   "user_name",
							Order:   1,
							Sort:    "asc",
						}},
					},
					"name": &model.FieldType{
						FieldName:           "name",
						IsFieldTypeRequired: true,
						Kind:                model.TypeID,
						TypeIDSize:          model.DefaultCharacterSize,
						IndexInfo: []*model.TableProperties{{
							IsIndex:  false,
							IsUnique: true,
							Group:    "user_name",
							Order:    1,
							Sort:     "asc",
						}},
					},
					"customer_id": &model.FieldType{
						FieldName:           "customer_id",
						IsFieldTypeRequired: true,
						Kind:                model.TypeID,
						TypeIDSize:          model.DefaultCharacterSize,
						IsForeign:           true,
						JointTable: &model.TableProperties{
							To:             "id",
							Table:          "customer",
							OnDelete:       "CASCADE",
							ConstraintName: "c_tweet_customer_id",
						},
					},
					"order_dates": &model.FieldType{
						FieldName: "order_dates",
						IsList:    true,
						Kind:      model.TypeDateTime,
						IsLinked:  true,
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
			want: "type  table1 { " +
				"\n\tcol2: ID @primary @size(value: 100)" +
				"\n\tcol3: Integer @primary(order:2)@autoIncrement" +
				"\n\tcol4: Integer @primary(order:1)" +
				"\n\tcol5: Integer @primary" +
				"\n\tage: Float" +
				"\n\tamount: BigInteger" +
				"\n\tcoolDownInterval: SmallInteger" +
				"\n\tcreatedAt: DateTime  @createdAt" +
				"\n\trole: ID!  @size(value: 100)     @default(value: user)" +
				"\n\tspec: JSON" +
				"\n\tupdatedAt: DateTime   @updatedAt" +
				"\n\tfirst_name: ID!  @size(value: 100)    @index(group: \"user_name\", sort: \"asc\", order: 1)" +
				"\n\tname: ID!  @size(value: 100)   @unique(group: \"user_name\", sort: \"asc\",  order: 1)" +
				"\n\tcustomer_id: ID!  @size(value: 100) @foreign(table: customer, field: id ,onDelete: cascade)" +
				"\n\torder_dates: [DateTime]       @link(table: \"order\", from: \"id\", to: \"customer_id\", db:\"mongo\", field: \"order_date\")" +
				"\n}",
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
