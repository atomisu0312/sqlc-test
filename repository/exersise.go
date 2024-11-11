package repository

import (
	"context"
	"sqlc-test/gen"
)

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
