package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"sqlc-test/gen"
)

type Tx struct {
	db *sql.DB
}

func NewTx(db *sql.DB) *Tx {
	return &Tx{db: db}
}

// fn func()の引数はrepositoriesという表記にしておけば、どのサービスからも使える
// そのためにはサービスにおいてrepositoriesという変数を用意する必要がある
func (tx *Tx) ExecTx(ctx context.Context, fn func(*gen.Queries) error) error {
	// read-committedでトランザクションを開始する
	txx, err := tx.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	q := gen.New(tx.db).WithTx(txx) //
	err = fn(q)                     //ここで外から渡した関数を実行する。
	if err != nil {
		if rbErr := txx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return txx.Commit()
}
