package service

import (
	"context"
	"database/sql"
	"fmt"
	"sqlc-test/gen"
	"sqlc-test/repository"
	"sqlc-test/transaction"
	"time"
)

// CreateWorkoutTxParams はワークアウト作成のパラメータを表します
type CreateWorkoutTxParams struct {
	UserID    int64                    `json:"user_id"`
	StartDate time.Time                `json:"start_date"`
	Sets      []CreateWorkoutSetParams `json:"sets"`
}

// CreateWorkoutSetParams はワークアウトセット作成のパラメータを表します
type CreateWorkoutSetParams struct {
	ExerciseID int64 `json:"exercise_id"`
	Weight     int32 `json:"weight"`
}

// CreateWorkoutTxResult はワークアウト作成の結果を表します
type CreateWorkoutTxResult struct {
	WorkoutID int64 `json:"workout_id"`
}

// ハンドラから直接呼ばれるのがユースケース
type UseCase struct {
	*gen.Queries
	db *sql.DB
}

// NewUseCase は新しい UseCase インスタンスを作成します
func NewUseCase(db *sql.DB) *UseCase {
	return &UseCase{
		db:      db,
		Queries: gen.New(db),
	}
}

// AddWorkoutTx はワークアウトを作成するトランザクションを実行します
func (useCase *UseCase) AddWorkoutTx(ctx context.Context, arg CreateWorkoutTxParams) (CreateWorkoutTxResult, error) {
	var result CreateWorkoutTxResult
	tr := transaction.NewTx(useCase.db)
	err := tr.ExecTx(ctx, func(q *gen.Queries) error {
		repo := repository.NewExerciseRepository(q)
		workout, err := repo.CreateExercise(ctx, "ExerciseDDDD")
		if err != nil {
			return fmt.Errorf("error create workout %w", err)
		}
		for _, set := range arg.Sets {
			setParams := gen.CreateSetParams{
				ExerciseID: set.ExerciseID,
				Weight:     set.Weight,
			}
			// 例として、2番目のセットでエラーをスローする
			// if i == 1 {
			// 	return fmt.Errorf("simulated error for rollback test")
			// }
			_, err := repo.CreateSet(ctx, setParams)
			if err != nil {
				return fmt.Errorf("error create workout set %w", err)
			}
		}
		result.WorkoutID = workout
		return nil
	})
	return result, err
}
