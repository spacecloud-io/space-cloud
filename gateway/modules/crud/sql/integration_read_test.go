// +build integration

package sql

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestSQL_Read(t *testing.T) {
	var distinctCol = "amount"
	var col int64 = 10
	type args struct {
		ctx context.Context
		col string
		req *model.ReadRequest
	}
	tests := []struct {
		name           string
		args           args
		wantCount      int64
		wantErr        bool
		wantReadResult []interface{}
	}{
		{
			name: "Simple read with no select & where clause",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Operation: utils.All,
				},
			},
			wantCount: 20,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		// equal operator
		{
			name: "Read where id equals 1",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{"$eq": "1"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,

			wantReadResult: []interface{}{map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(float32(12.3))}},
		},
		{
			name: "Read where order_date equals 2001-11-01T14:29:36Z",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"order_date": map[string]interface{}{"$eq": "2001-11-01 14:29:36"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,

			wantReadResult: []interface{}{map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)}},
		},
		{
			name: "Read where amount equals 10",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"amount": map[string]interface{}{"$eq": 10},
					},
					Operation: utils.All,
				},
			},
			wantCount: 2,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
			},
		},
		{
			name: "Read where stars equals 4.5",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"stars": map[string]interface{}{"$eq": 4.5},
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		{
			name: "Read where is_prime is true",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"is_prime": map[string]interface{}{"$eq": 1},
					},
					Operation: utils.All,
				},
			},
			wantCount: 9,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
			},
		},
		{
			name: "Read where product_id equals fridge",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"product_id": map[string]interface{}{"$eq": "fridge"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
			},
		},
		// implicit equal operator
		{
			name: "Read where id equals 1",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"id": "1",
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,

			wantReadResult: []interface{}{map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(float32(12.3))}},
		},
		{
			name: "Read where order_date equals 2001-11-01T14:29:36Z",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"order_date": "2001-11-01 14:29:36",
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,

			wantReadResult: []interface{}{map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)}},
		},
		{
			name: "Read where amount equals 10",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"amount": 10,
					},
					Operation: utils.All,
				},
			},
			wantCount: 2,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
			},
		},
		{
			name: "Read where stars equals 4.5",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"stars": 4.5,
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		{
			name: "Read where is_prime is true",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"is_prime": 1,
					},
					Operation: utils.All,
				},
			},
			wantCount: 9,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
			},
		},
		{
			name: "Read where product_id equals fridge",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"product_id": "fridge",
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
			},
		},
		// not equal operator
		{
			name: "Read where id not equal 1",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{"$ne": "1"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 19,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		{
			name: "Read where order_date not equal 2001-11-01T14:29:36Z",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"order_date": map[string]interface{}{"$ne": "2001-11-01 14:29:36"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 19,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		{
			name: "Read where amount not equal 10",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"amount": map[string]interface{}{"$ne": 10},
					},
					Operation: utils.All,
				},
			},
			wantCount: 18,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		{
			name: "Read where stars not equals 4.5",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"stars": map[string]interface{}{"$ne": 4.5},
					},
					Operation: utils.All,
				},
			},
			wantCount: 19,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
			},
		},
		{
			name: "Read where is_prime is not false",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"is_prime": map[string]interface{}{"$ne": 0},
					},
					Operation: utils.All,
				},
			},
			wantCount: 9,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
			},
		},
		{
			name: "Read where product_id not equal fridge",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"product_id": map[string]interface{}{"$ne": "fridge"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 19,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		// greater than operator
		{
			name: "Read where id greater than 19",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{"$gt": "19"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 9,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		{
			name: "Read where order_date greater than 2045-11-25T14:29:36Z",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"order_date": map[string]interface{}{"$gt": "2050-11-25 14:29:36"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)}},
		},
		{
			name: "Read where amount greater than 97",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"amount": map[string]interface{}{"$gt": 97},
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
			},
		},
		{
			name: "Read where stars greater than 1500.5",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"stars": map[string]interface{}{"$gt": 1500.5},
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
			},
		},
		// doesnt work {
		// 	name: "Read where is_prime is true",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "orders",
		// 		req: &model.ReadRequest{
		// 			Find: map[string]interface{}{
		// 				"is_prime": map[string]interface{}{"$gt": 5},
		// 			},
		// 			Operation: utils.All,
		// 		},
		// 	},
		// 	wantCount: 9,
		// 	wantErr:   false,
		// 	wantReadResult: []interface{}{
		// 		map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
		// 		map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
		// 		map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
		// 		map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
		// 		map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
		// 		map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
		// 		map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
		// 		map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
		// 		map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
		// 	},
		// },
		{
			name: "Read where product_id greater than shoes",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"product_id": map[string]interface{}{"$gt": "shoes"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 3,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
			},
		},
		// greater than equal to operator
		{
			name: "Read where id greater than equal to 19",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{"$gte": "19"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 10,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		{
			name: "Read where order_date greater than equal to 2045-11-25T14:29:36Z",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"order_date": map[string]interface{}{"$gte": "2050-11-25 14:29:36"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)}},
		},
		{
			name: "Read where amount greater than equal to 97",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"amount": map[string]interface{}{"$gte": 97},
					},
					Operation: utils.All,
				},
			},
			wantCount: 2,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
			},
		},
		{
			name: "Read where stars greater than equal to 1500.5",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"stars": map[string]interface{}{"$gte": 1500.5},
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
			},
		},
		{
			name: "Read where product_id greater than equal to shoes",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"product_id": map[string]interface{}{"$gte": "shoes"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 4,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
			},
		},
		// less than operator
		{
			name: "Read where id less than 2",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{"$lt": "2"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 11,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
			},
		},
		{
			name: "Read where order_date less than 2002-11-01T14:29:36Z",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"order_date": map[string]interface{}{"$lt": "2002-11-01 14:29:36"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 2,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
			},
		},
		{
			name: "Read where amount less than 11",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"amount": map[string]interface{}{"$lt": 11},
					},
					Operation: utils.All,
				},
			},
			wantCount: 2,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
			},
		},
		{
			name: "Read where stars less thans 4.5",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"stars": map[string]interface{}{"$lt": 4.5},
					},
					Operation: utils.All,
				},
			},
			wantCount: 3,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
			},
		},
		{
			name: "Read where product_id less than fridge",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"product_id": map[string]interface{}{"$lt": "books"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 2,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
			},
		},
		// less than equal to operator
		{
			name: "Read where id less than equal to 2",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{"$lte": "2"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 12,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
			},
		},
		{
			name: "Read where order_date less than equal to 2002-11-01T14:29:36Z",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"order_date": map[string]interface{}{"$lte": "2002-11-01 14:29:36"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 2,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
			},
		},
		{
			name: "Read where amount less than equal to 19",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"amount": map[string]interface{}{"$lte": 19},
					},
					Operation: utils.All,
				},
			},
			wantCount: 6,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
			},
		},
		{
			name: "Read where stars less than equal tos 4.5",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"stars": map[string]interface{}{"$lte": 4.5},
					},
					Operation: utils.All,
				},
			},
			wantCount: 4,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		{
			name: "Read where product_id less than equal to fridge",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"product_id": map[string]interface{}{"$lte": "books"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 3,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
			},
		},
		// in operator
		{
			name: "Read where id with in operator [1,2,19]",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{"$in": []interface{}{"1", "2", "19"}},
					},
					Operation: utils.All,
				},
			},
			wantCount: 3,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
			},
		},
		{
			name: "Read where order_date with in operator [2001-11-01T14:29:36Z,2002-11-05T14:29:36Z]]",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"order_date": map[string]interface{}{"$in": []interface{}{"2001-11-01 14:29:36", "2002-11-05 14:29:36"}},
					},
					Operation: utils.All,
				},
			},
			wantCount: 2,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
			},
		},
		{
			name: "Read where amount with in operator [10,19,14]",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"amount": map[string]interface{}{"$in": []interface{}{10, 19, 14}},
					},
					Operation: utils.All,
				},
			},
			wantCount: 6,
			wantErr:   false,

			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
			},
		},
		// not working
		// {
		// 	name: "Read where stars with in operator [4.5,1.54]",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "orders",
		// 		req: &model.ReadRequest{
		// 			Find: map[string]interface{}{
		// 				"stars": map[string]interface{}{"$in": []interface{}{4.5, 1.54}},
		// 			},
		// 			Operation: utils.All,
		// 		},
		// 	},
		// 	wantCount: 2,
		// 	wantErr:   false,
		// 	wantReadResult: []interface{}{
		// 		map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": `{"city": "chennai", "pinCode": 40560014}`, "stars": float32(1.54)},
		// 		map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
		// 	},
		// },
		{
			name: "Read where product_id with in operator [books,bed,basket]",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"product_id": map[string]interface{}{"$in": []interface{}{"books", "bed", "basket"}},
					},
					Operation: utils.All,
				},
			},
			wantCount: 3,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
			},
		},
		// not in operator
		{
			name: "Read where id with not in operator [1,2,3]",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"id": map[string]interface{}{"$nin": []interface{}{"1", "2", "3"}},
					},
					Operation: utils.All,
				},
			},
			wantCount: 17,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		{
			name: "Read where order_date with not in operator [2001-11-01T14:29:36Z,2002-11-05T14:29:36Z]]",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"order_date": map[string]interface{}{"$nin": []interface{}{"2001-11-01 14:29:36", "2002-11-05 14:29:36"}},
					},
					Operation: utils.All,
				},
			},
			wantCount: 18,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)}},
		},
		{
			name: "Read where amount with not in operator [37,97]",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"amount": map[string]interface{}{"$nin": []interface{}{37, 97}},
					},
					Operation: utils.All,
				},
			},
			wantCount: 17,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
			},
		},
		// not working
		// {
		// 	name: "Read where stars with not in operator [4.5,1.54]",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "orders",
		// 		req: &model.ReadRequest{
		// 			Find: map[string]interface{}{
		// 				"stars": map[string]interface{}{"$nin": []interface{}{4.5, 1.54}},
		// 			},
		// 			Operation: utils.All,
		// 		},
		// 	},
		// 	wantCount: 18,
		// 	wantErr:   false,
		// 	wantReadResult: []interface{}{
		// 		map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
		// 		map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
		// 		map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
		// 		map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
		// 		map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
		// 		map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
		// 		map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
		// 		map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
		// 		map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
		// 		map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
		// 		map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
		// 		map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
		// 		map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
		// 		map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
		// 		map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
		// 		map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
		// 		map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
		// 		map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
		// 	},
		// },
		// contains operator
		{
			name: "Read where address contains city mumbai",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"address": map[string]interface{}{"$contains": map[string]interface{}{"city": "newyork"}},
					},
					Operation: utils.All,
				},
			},
			wantCount: 1,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
			},
		},
		// regex operator
		{
			name: "Read where product_id with regex",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"product_id": map[string]interface{}{"$regex": "^j"},
					},
					Operation: utils.All,
				},
			},
			wantCount: 3,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)}},
		},
		// or operator
		{
			name: "Read where using or",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Find: map[string]interface{}{
						"$or": []interface{}{
							map[string]interface{}{"product_id": map[string]interface{}{"$regex": "^j"}},
							map[string]interface{}{"address": map[string]interface{}{"$contains": map[string]interface{}{"city": "newyork"}}},
						},
					},
					Operation: utils.All,
				},
			},
			wantCount: 4,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)}},
		},
		// // aggregate operator
		// {
		// 	name: "Read aggregate sum",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "orders",
		// 		req: &model.ReadRequest{
		// 			Aggregate: map[string][]string{
		// 				"sum": {"amount"},
		// 			},
		// 			Operation: utils.All,
		// 		},
		// 	},
		// 	wantCount: 1,
		// 	wantErr:   false,
		// 	wantReadResult: []interface{}{
		// 		// TODO: WHY THE FLOAT VALUE IS LIKE THIS
		// 		map[string]interface{}{"aggregate": map[string]interface{}{"sum": map[string]interface{}{"amount": int64(1088)}}},
		// 	},
		// },
		// {
		// 	name: "Read aggregate count",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "orders",
		// 		req: &model.ReadRequest{
		// 			Aggregate: map[string][]string{
		// 				"count": {"amount"},
		// 			},
		// 			Operation: utils.All,
		// 		},
		// 	},
		// 	wantCount: 1,
		// 	wantErr:   false,
		// 	wantReadResult: []interface{}{
		// 		map[string]interface{}{"aggregate": map[string]interface{}{"count": map[string]interface{}{"amount": int64(20)}}},
		// 	},
		// },
		// {
		// 	name: "Read aggregate avg",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "orders",
		// 		req: &model.ReadRequest{
		// 			Aggregate: map[string][]string{
		// 				"avg": {"amount"},
		// 			},
		// 			Operation: utils.All,
		// 		},
		// 	},
		// 	wantCount: 1,
		// 	wantErr:   false,
		// 	wantReadResult: []interface{}{
		// 		map[string]interface{}{"aggregate": map[string]interface{}{"avg": map[string]interface{}{"amount": 20}}},
		// 	},
		// },
		// {
		// 	name: "Read aggregate max",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "orders",
		// 		req: &model.ReadRequest{
		// 			Aggregate: map[string][]string{
		// 				"max": {"amount"},
		// 			},
		// 			Operation: utils.All,
		// 		},
		// 	},
		// 	wantCount: 1,
		// 	wantErr:   false,
		// 	wantReadResult: []interface{}{
		// 		map[string]interface{}{"aggregate": map[string]interface{}{"max": map[string]interface{}{"amount": 20}}},
		// 	},
		// },
		// {
		// 	name: "Read aggregate min",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "orders",
		// 		req: &model.ReadRequest{
		// 			Aggregate: map[string][]string{
		// 				"min": {"amount"},
		// 			},
		// 			Operation: utils.All,
		// 		},
		// 	},
		// 	wantCount: 1,
		// 	wantErr:   false,
		// 	wantReadResult: []interface{}{
		// 		map[string]interface{}{"aggregate": map[string]interface{}{"min": map[string]interface{}{"amount": 20}}},
		// 	},
		// },
		// misc
		{
			name: "Read with operation select with limit",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Operation: utils.All,
					Options:   &model.ReadOptions{Limit: &col, Select: map[string]int32{"id": 1}},
				},
			},
			wantCount: 10,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1"},
				map[string]interface{}{"id": "10"},
				map[string]interface{}{"id": "11"},
				map[string]interface{}{"id": "12"},
				map[string]interface{}{"id": "13"},
				map[string]interface{}{"id": "14"},
				map[string]interface{}{"id": "15"},
				map[string]interface{}{"id": "16"},
				map[string]interface{}{"id": "17"},
				map[string]interface{}{"id": "18"},
			},
		},
		{
			name: "Read with operation select & limit & sort ascending",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Operation: utils.All,
					Options:   &model.ReadOptions{Sort: []string{"id"}, Select: map[string]int32{"id": 1}, Limit: &col},
				},
			},
			wantCount: 10,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1"},
				map[string]interface{}{"id": "10"},
				map[string]interface{}{"id": "11"},
				map[string]interface{}{"id": "12"},
				map[string]interface{}{"id": "13"},
				map[string]interface{}{"id": "14"},
				map[string]interface{}{"id": "15"},
				map[string]interface{}{"id": "16"},
				map[string]interface{}{"id": "17"},
				map[string]interface{}{"id": "18"},
			},
		},
		{
			name: "Read with operation limit & sort descending",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Operation: utils.All,
					Options:   &model.ReadOptions{Sort: []string{"-id"}, Limit: &col},
				},
			},
			wantCount: 10,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
				map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
				map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
				map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
				map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
				map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
				map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
				map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
				map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
			},
		},

		{
			name: "Read with operation limit to 10 rows",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Operation: utils.All,
					Options:   &model.ReadOptions{Limit: &col},
				},
			},
			wantCount: 10,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
				map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
				map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
				map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
				map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
				map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
				map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
				map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
				map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
				map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
			},
		},
		// doesnt work for mysql as limit also required
		// {
		// 	name: "Read with operation skip 10 rows",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		col: "orders",
		// 		req: &model.ReadRequest{
		// 			Operation: utils.All,
		// 			Options:   &model.ReadOptions{Skip: &col},
		// 		},
		// 	},
		// 	wantCount: 10,
		// 	wantErr:   false,
		// 	wantReadResult: []interface{}{
		// 		map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
		// 		map[string]interface{}{"id": "2", "order_date": "2001-11-12T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "shoes", "address": (`{"city": "newyork", "pinCode": 4003014}`), "stars": float32(51.3)},
		// 		map[string]interface{}{"id": "3", "order_date": "2002-11-05T14:29:36Z", "amount": int64(52), "is_prime": int64(1), "product_id": "fridge", "address": (`{"city": "amsterdam", "pinCode": 4200014}`), "stars": float32(1.37)},
		// 		map[string]interface{}{"id": "4", "order_date": "2002-11-02T14:29:36Z", "amount": int64(95), "is_prime": int64(0), "product_id": "door", "address": (`{"city": "pune", "pinCode": 4000134}`), "stars": float32(41.3)},
		// 		map[string]interface{}{"id": "5", "order_date": "2004-11-03T14:29:36Z", "amount": int64(79), "is_prime": int64(0), "product_id": "basket", "address": (`{"city": "hyderabad", "pinCode": 4030014}`), "stars": float32(21.3)},
		// 		map[string]interface{}{"id": "6", "order_date": "2004-11-05T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "books", "address": (`{"city": "bangalore", "pinCode": 400014}`), "stars": float32(81.3)},
		// 		map[string]interface{}{"id": "7", "order_date": "2006-11-03T14:29:36Z", "amount": int64(85), "is_prime": int64(0), "product_id": "cover", "address": (`{"city": "surat", "pinCode": 4000134}`), "stars": float32(81.3)},
		// 		map[string]interface{}{"id": "8", "order_date": "2006-11-06T14:29:36Z", "amount": int64(97), "is_prime": int64(0), "product_id": "sheets", "address": (`{"city": "ahemdabad", "pinCode": 40450014}`), "stars": float32(122.3)},
		// 		map[string]interface{}{"id": "9", "order_date": "2008-11-21T14:29:36Z", "amount": int64(94), "is_prime": int64(0), "product_id": "bed", "address": (`{"city": "venice", "pinCode": 4000154}`), "stars": float32(111.3)},
		// 		map[string]interface{}{"id": "10", "order_date": "2008-11-13T14:29:36Z", "amount": int64(93), "is_prime": int64(1), "product_id": "sofa", "address": (`{"city": "california", "pinCode": 40006514}`), "stars": float32(1.96)},
		// 		map[string]interface{}{"id": "11", "order_date": "2050-11-13T14:29:36Z", "amount": int64(91), "is_prime": int64(1), "product_id": "pillow", "address": (`{"city": "chennai", "pinCode": 40560014}`), "stars": float32(1.54)},
		// 		map[string]interface{}{"id": "12", "order_date": "2050-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(1), "product_id": "mat", "address": (`{"city": "berlin", "pinCode": 40001654}`), "stars": float32(1134.3)},
		// 		map[string]interface{}{"id": "13", "order_date": "2045-11-15T14:29:36Z", "amount": int64(19), "is_prime": int64(1), "product_id": "juice", "address": (`{"city": "moscow", "pinCode": 40005714}`), "stars": float32(451.3)},
		// 		map[string]interface{}{"id": "14", "order_date": "2045-11-25T14:29:36Z", "amount": int64(14), "is_prime": int64(1), "product_id": "mixer", "address": (`{"city": "paris", "pinCode": 40005614}`), "stars": float32(761.433)},
		// 		map[string]interface{}{"id": "15", "order_date": "2080-11-12T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "grinder", "address": (`{"city": "vein", "pinCode": 400056714}`), "stars": float32(1435.3)},
		// 		map[string]interface{}{"id": "16", "order_date": "2016-11-10T14:29:36Z", "amount": int64(98), "is_prime": int64(1), "product_id": "washing", "address": (`{"city": "islamabad", "pinCode": 40530014}`), "stars": float32(131.3)},
		// 		map[string]interface{}{"id": "17", "order_date": "2026-11-05T14:29:36Z", "amount": int64(28), "is_prime": int64(0), "product_id": "powder", "address": (`{"city": "dhaka", "pinCode": 400014}`), "stars": float32(1567.3)},
		// 		map[string]interface{}{"id": "18", "order_date": "2015-11-05T14:29:36Z", "amount": int64(26), "is_prime": int64(0), "product_id": "cake", "address": (`{"city": "los-angeles", "pinCode": 40001434}`), "stars": float32(71.3)},
		// 		map[string]interface{}{"id": "19", "order_date": "2015-11-05T14:29:36Z", "amount": int64(37), "is_prime": int64(0), "product_id": "jam", "address": (`{"city": "mumbai", "pinCode": 40001445}`), "stars": float32(451.3)},
		// 		map[string]interface{}{"id": "20", "order_date": "2010-11-05T14:29:36Z", "amount": int64(19), "is_prime": int64(0), "product_id": "jeans", "address": `{"city": "mumbai", "pinCode": 40002314}`, "stars": float32(4.5)},
		// 	},
		// },
		{
			name: "Read with op type distinct",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Operation: utils.Distinct,
					Options:   &model.ReadOptions{Distinct: &distinctCol},
				},
			},
			wantCount: 15,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"amount": int64(10)},
				map[string]interface{}{"amount": int64(52)},
				map[string]interface{}{"amount": int64(95)},
				map[string]interface{}{"amount": int64(79)},
				map[string]interface{}{"amount": int64(85)},
				map[string]interface{}{"amount": int64(97)},
				map[string]interface{}{"amount": int64(94)},
				map[string]interface{}{"amount": int64(93)},
				map[string]interface{}{"amount": int64(91)},
				map[string]interface{}{"amount": int64(37)},
				map[string]interface{}{"amount": int64(19)},
				map[string]interface{}{"amount": int64(14)},
				map[string]interface{}{"amount": int64(98)},
				map[string]interface{}{"amount": int64(28)},
				map[string]interface{}{"amount": int64(26)},
				map[string]interface{}{"amount": int64(37)},
			},
		},
		{
			name: "Read with op type distinct",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Operation: utils.Distinct,
					Options:   &model.ReadOptions{Distinct: &distinctCol},
				},
			},
			wantCount: 15,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"amount": int64(10)},
				map[string]interface{}{"amount": int64(52)},
				map[string]interface{}{"amount": int64(95)},
				map[string]interface{}{"amount": int64(79)},
				map[string]interface{}{"amount": int64(85)},
				map[string]interface{}{"amount": int64(97)},
				map[string]interface{}{"amount": int64(94)},
				map[string]interface{}{"amount": int64(93)},
				map[string]interface{}{"amount": int64(91)},
				map[string]interface{}{"amount": int64(37)},
				map[string]interface{}{"amount": int64(19)},
				map[string]interface{}{"amount": int64(14)},
				map[string]interface{}{"amount": int64(98)},
				map[string]interface{}{"amount": int64(28)},
				map[string]interface{}{"amount": int64(26)},
				map[string]interface{}{"amount": int64(37)},
			},
		},
		{
			name: "Read with op type count",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Operation: utils.Count,
				},
			},
			wantCount: 20,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"count": 20},
			},
		},
		{
			name: "Simple read with no select & where clause op type one",
			args: args{
				ctx: context.Background(),
				col: "orders",
				req: &model.ReadRequest{
					Operation: utils.One,
				},
			},
			wantCount: 1,
			wantErr:   false,
			wantReadResult: []interface{}{
				map[string]interface{}{"id": "1", "order_date": "2001-11-01T14:29:36Z", "amount": int64(10), "is_prime": int64(1), "product_id": "smart-phone", "address": (`{"city": "mumbai", "pinCode": 400014}`), "stars": float32(12.3)},
			},
		},
	}
	db, err := Init(utils.DBType(*dbType), true, *connection, "myproject")
	if err != nil {
		t.Fatal("Couldn't establishing connection with database", dbType)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, gotReadResult, gotErr := db.Read(tt.args.ctx, tt.args.col, tt.args.req)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("Read() count mismatch got = %v, want %v", gotCount, tt.wantCount)
			}

			for _, wantReadResult := range tt.wantReadResult {
				switch v := gotReadResult.(type) {
				case map[string]interface{}:
					if !reflect.DeepEqual(gotReadResult, wantReadResult) {
						t.Errorf("Read() mismatch value got %v\n want  %v", gotReadResult, wantReadResult)
					}
				case []interface{}:
					isFound := false
					// t.Log("got", v)
					for _, gotReadMap := range v {
						// for wantKey, wantValue := range wantReadResult.(map[string]interface{}) {
						// 	gotValue, ok := gotReadMap.(map[string]interface{})[wantKey]
						// 	if !ok {
						// 		t.Errorf("Read() missing field %v in result", wantKey)
						// 	}
						if cmp.Equal(gotReadMap, wantReadResult) {
							isFound = true
							// t.Errorf("Read() value mismatch got %v type %v want %v type %v", gotValue, reflect.TypeOf(gotValue), wantValue, reflect.TypeOf(wantValue))
						}
						// }
					}
					if !isFound {
						t.Log("Got", v)
						t.Errorf("Read() want value not found  %v", wantReadResult)
					}
				}
			}
		})
	}
}
