package postgresql

import (
	"testing"

	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"
	"github.com/ikotiki/go-rest-api-service-subscriptions/logger"
	"github.com/ikotiki/sqlbuilder"
	"github.com/ikotiki/sqlbuilder/builder"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func Test_parseQueryArgs(t *testing.T) {
	log := logger.InitLoggerByFlag("debug", false)
	sqlBuilder, err := sqlbuilder.NewSQLBuilder("postgres")
	b := postgresSQLBuilder{sqlBuilder}
	if err != nil {
		t.Fatalf("Can't create sqlbuilder: %v", err)
	}

	tests := []struct {
		name  string
		input *storage.QueryArgs
		want  *builder.SelectArguments
	}{
		{
			name: "Ok",
			input: &storage.QueryArgs{
				From: "subscriptions",
				Order: []storage.OrderStruct{
					{
						OrderBy: "name",
						Order:   storage.OrderASC,
					},
					{
						OrderBy: "id",
						Order:   storage.OrderDECS,
					},
				},
				Where: []storage.Where{
					{
						Column:   "name",
						Operator: storage.OpEqual,
						Value:    "test",
					},
					{
						Column:   "id",
						Operator: storage.OpEqual,
						Value:    "1",
					},
				},
				Offset: 10,
				Limit:  10,
			},
			want: &builder.SelectArguments{
				From: builder.Table("subscriptions"),
				Where: []builder.Where{
					{
						Column:   "name",
						Operator: "=",
						Value:    "test",
					},
					{
						Column:   "id",
						Operator: "=",
						Value:    "1",
					},
				},
				OrderBy: []builder.OrderBy{
					{
						Column: "name",
						Order:  "ASC",
					},
					{
						Column: "id",
						Order:  "DESC",
					},
				},
				Limit: builder.Limit{
					Offset: 10,
					Limit:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := b.parseQueryArgs(tt.input)

			log.Debug().Interface("got", got).Msg("got")

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_buildParts(t *testing.T) {
	builder, err := sqlbuilder.NewSQLBuilder("postgres")
	b := postgresSQLBuilder{builder}
	if err != nil {
		t.Fatalf("Can't create sqlbuilder: %v", err)
	}

	type res struct {
		Query     string
		QueryArgs []interface{}
	}
	tests := []struct {
		name  string
		input *storage.QueryArgs
		want  *res
	}{
		{
			name: "Ok",
			input: &storage.QueryArgs{
				From: "subscriptions",
				Order: []storage.OrderStruct{
					{
						OrderBy: "name",
						Order:   storage.OrderASC,
					},
					{
						OrderBy: "id",
						Order:   storage.OrderDECS,
					},
				},
				Where: []storage.Where{
					{
						Column:   "name",
						Operator: storage.OpEqual,
						Value:    "test",
					},
					{
						Column:   "id",
						Operator: storage.OpEqual,
						Value:    1,
					},
				},
				Offset: 10,
				Limit:  10,
			},
			want: &res{
				Query: `SELECT * FROM subscriptions WHERE name = $1 AND id = $2 ORDER BY name ASC, id DESC LIMIT $3 OFFSET $4`,
				QueryArgs: []interface{}{
					"test",
					int(1),
					int64(10),
					int64(10),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := []string{"select", "from", "where", "group_by", "order_by", "limit"}
			query, queryArgs := b.buildParts(parts, tt.input, 1)
			got := &res{
				Query:     query,
				QueryArgs: queryArgs,
			}

			log.Debug().Interface("got", got).Msg("got")

			assert.Equal(t, tt.want, got)
		})
	}
}
