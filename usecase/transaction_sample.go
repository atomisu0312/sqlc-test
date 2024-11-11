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
type useCase struct {
	db *sql.DB
}

type UseCase interface {
	AddWorkoutTx(ctx context.Context, arg CreateWorkoutTxParams) (CreateWorkoutTxResult, error)
}

// NewUseCase は新しい UseCase インスタンスを作成します
func NewUseCase(db *sql.DB) UseCase {
	return &useCase{
		db: db,
	}
}

// AddWorkoutTx はワークアウトを作成するトランザクションを実行します
func (useCase *useCase) AddWorkoutTx(ctx context.Context, arg CreateWorkoutTxParams) (CreateWorkoutTxResult, error) {
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
