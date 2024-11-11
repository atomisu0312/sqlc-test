package main

import (
	"context"
	"fmt"
	"log"
	"sqlc-test/config"
	service "sqlc-test/usecase"
	"time"

	_ "github.com/lib/pq" // PostgreSQL ドライバをインポート
)

func main() {

	// Open the database
	db, err := config.NewDbConnection()

	if err != nil {
		panic(err)
	}

	// Connectivity check
	if err := db.DB.Ping(); err != nil {
		log.Fatalln("Error from database ping:", err)
	}

	ctx := context.Background()

	// ワークアウトトランザクションのパラメータを設定
	workoutParams := service.CreateWorkoutTxParams{
		UserID:    1,
		StartDate: time.Now(),
		Sets: []service.CreateWorkoutSetParams{
			{ExerciseID: 1, Weight: 100},
			{ExerciseID: 2, Weight: 200},
		},
	}

	// UseCaseのインスタンスを作成
	useCase := service.NewUseCase(db.DB)

	// ワークアウトトランザクションを実行
	result, err := useCase.AddWorkoutTx(ctx, workoutParams)

	if err != nil {
		log.Fatalln("Error creating workout transaction:", err)
	}

	// 結果を表示
	fmt.Printf("Workout created with ID: %d\n", result.WorkoutID)
}
