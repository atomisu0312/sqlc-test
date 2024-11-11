package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"sqlc-test/gen"
)

type Store struct {
	*gen.Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db, Queries: gen.New(db)}
}

// fn func()の引数はrepositoriesという表記にしておけば、どのサービスからも使える
// そのためにはサービスにおいてrepositoriesという変数を用意する必要がある
func (store *Store) execTx(ctx context.Context, fn func(*gen.Queries) error) error {
	// read-committedでトランザクションを開始する
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	q := store.Queries.WithTx(tx) // `WithTx` メソッドを使用してトランザクションを設定 //ここをコンポーネントに対するメソッドに変える
	err = fn(q)                   //ここで外から渡した関数を実行する。
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
