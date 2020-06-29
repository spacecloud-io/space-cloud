package crud

import (
	"context"
	"testing"
	"time"
)

func TestModule_createBatch(t *testing.T) {
	type fields struct {
		batchMapTableToChan batchMap
	}
	type args struct {
		project string
		dbAlias string
		col     string
		doc     interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "No value to be inserted",
			fields: fields{
				batchMapTableToChan: map[string]map[string]map[string]batchChannels{
					"myproject": {
						"db": {
							"customers": {
								request: make(batchRequestChan, 20),
								closeC:  make(chan struct{}),
							},
						},
					},
				},
			},
			args: args{
				project: "myproject",
				dbAlias: "db",
				col:     "customers",
				doc:     []interface{}{},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Single value inserted",
			fields: fields{
				batchMapTableToChan: map[string]map[string]map[string]batchChannels{
					"myproject": {
						"db": {
							"customers": {
								request: make(batchRequestChan, 20),
								closeC:  make(chan struct{}),
							},
						},
					},
				},
			},
			args: args{
				project: "myproject",
				dbAlias: "db",
				col:     "customers",
				doc: map[string]interface{}{
					"id":   "1",
					"name": "John",
					"age":  20,
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Multiple values inserted",
			fields: fields{
				batchMapTableToChan: map[string]map[string]map[string]batchChannels{
					"myproject": {
						"db": {
							"customers": {
								request: make(batchRequestChan, 20),
								closeC:  make(chan struct{}),
							},
						},
					},
				},
			},
			args: args{
				project: "myproject",
				dbAlias: "db",
				col:     "customers",
				doc: []interface{}{
					map[string]interface{}{
						"id":   "1",
						"name": "John",
						"age":  20,
					},
					map[string]interface{}{
						"id":   "2",
						"name": "Sam",
						"age":  30,
					},
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "Channel not found provided collection",
			fields: fields{
				batchMapTableToChan: map[string]map[string]map[string]batchChannels{
					"myproject": {
						"db": {
							"orders": {
								request: make(batchRequestChan, 20),
								closeC:  make(chan struct{}),
							},
						},
					},
				},
			},
			args: args{
				project: "myproject",
				dbAlias: "db",
				col:     "customers",
				doc: []interface{}{
					map[string]interface{}{
						"id":   "1",
						"name": "John",
						"age":  20,
					},
					map[string]interface{}{
						"id":   "2",
						"name": "Sam",
						"age":  30,
					},
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				batchMapTableToChan: tt.fields.batchMapTableToChan,
			}
			go func() {
				_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				v := <-tt.fields.batchMapTableToChan[tt.args.project][tt.args.dbAlias][tt.args.col].request
				v.response <- batchResponse{err: nil}
			}()
			got, err := m.createBatch(tt.args.project, tt.args.dbAlias, tt.args.col, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("createBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createBatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPreparedQueryKey(t *testing.T) {
	type args struct {
		dbAlias string
		id      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Correct value",
			args: args{
				dbAlias: "db",
				id:      "prepare1",
			},
			want: "db--prepare1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPreparedQueryKey(tt.args.dbAlias, tt.args.id); got != tt.want {
				t.Errorf("getPreparedQueryKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
