package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"

	"cloud.google.com/go/civil"
)

type S struct {
	bun.BaseModel `bun:"table:test,alias:t"`
	N             int        `json:"n" bun:"n,pk"`
	D             civil.Date `json:"d" bun:"d"`
}

func main() {
	s := S{}
	var err error
	s.D, err = civil.ParseDate("2023-12-27")
	if err != nil {
		panic(err)
	}
	s.N = 42
	fmt.Println(s)

	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())

	ctx := context.Background()
	_, err = db.NewCreateTable().Model((*S)(nil)).Exec(ctx)
	if err != nil {
		panic(err)
	}

	err = db.ResetModel(ctx, (*S)(nil))
	if err != nil {
		panic(err)
	}

	_, err = db.NewInsert().Model(&s).Exec(ctx)
	if err != nil {
		panic(err)
	}

	s = S{}

	err = db.NewSelect().Model(&s).Scan(ctx)
	if err != nil {
		panic(err)
	}

	b, err = json.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

}
