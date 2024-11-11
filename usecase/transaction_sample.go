package service

import (
	"context"
	"database/sql"
	"fmt"
	"sqlc-test/gen"
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

// UseCase はデータベース操作を行うための構造体です
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

// ExerciseRepository はエクササイズ関連のデータベース操作を定義するインターフェースです
type ExerciseRepository interface {
	CreateExercise(ctx context.Context, exerciseName string) (int64, error)
	CreateSet(ctx context.Context, params gen.CreateSetParams) (gen.GowebappSet, error)
	UpdateSet(ctx context.Context, params gen.UpdateSetParams) (gen.GowebappSet, error)
}

// ExerciseRepositoryImpl は ExerciseRepository インターフェースを実装する構造体です
type ExerciseRepositoryImpl struct {
	queries *gen.Queries
}

// NewExerciseRepository は ExerciseRepositoryImpl の新しいインスタンスを返します
func NewExerciseRepository(queries *gen.Queries) ExerciseRepository {
	return &ExerciseRepositoryImpl{queries: queries}
}

// CreateExercise はエクササイズを作成します
func (repo *ExerciseRepositoryImpl) CreateExercise(ctx context.Context, exerciseName string) (int64, error) {
	return repo.queries.CreateExercise(ctx, exerciseName)
}

// CreateSet はセットを作成します
func (repo *ExerciseRepositoryImpl) CreateSet(ctx context.Context, params gen.CreateSetParams) (gen.GowebappSet, error) {
	return repo.queries.CreateSet(ctx, params)
}

// UpdateSet はセットを更新します
func (repo *ExerciseRepositoryImpl) UpdateSet(ctx context.Context, params gen.UpdateSetParams) (gen.GowebappSet, error) {
	return repo.queries.UpdateSet(ctx, params)
}

// execTx はトランザクションを実行します
func (useCase *UseCase) execTx(ctx context.Context, fn func(*gen.Queries) error) error {
	tx, err := useCase.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	q := useCase.Queries.WithTx(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// AddWorkoutTx はワークアウトを作成するトランザクションを実行します
func (useCase *UseCase) AddWorkoutTx(ctx context.Context, arg CreateWorkoutTxParams) (CreateWorkoutTxResult, error) {
	var result CreateWorkoutTxResult
	err := useCase.execTx(ctx, func(q *gen.Queries) error {
		repo := NewExerciseRepository(q)
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
