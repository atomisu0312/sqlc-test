package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sqlc-test/env"
	"sqlc-test/gen"
	service "sqlc-test/usecase"
	"time"

	_ "github.com/lib/pq" // PostgreSQL ドライバをインポート
)

func main() {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		env.GetAsString("DB_USER", "postgres"),
		env.GetAsString("DB_PASSWORD", "mysecretpassword"),
		env.GetAsString("DB_HOST", "localhost"),
		env.GetAsInt("DB_PORT", 5432),
		env.GetAsString("DB_NAME", "postgres"),
	)

	// Open the database
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}

	// Connectivity check
	if err := db.Ping(); err != nil {
		log.Fatalln("Error from database ping:", err)
	}

	// Create the store
	st := gen.New(db)

	ctx := context.Background()

	_, err = st.CreateUsers(ctx, gen.CreateUsersParams{
		UserName:     "testuser",
		PassWordHash: "hash",
		Name:         "test",
	})
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
	useCase := service.NewUseCase(db)

	// ワークアウトトランザクションを実行
	result, err := useCase.AddWorkoutTx(ctx, workoutParams)
	if err != nil {
		log.Fatalln("Error creating workout transaction:", err)
	}
	// 結果を表示
	fmt.Printf("Workout created with ID: %d\n", result.WorkoutID)

	if err != nil {
		log.Fatalln("Error creating user :", err)
	}

	eid, err := st.CreateExercise(ctx, "Exercise1")

	if err != nil {
		log.Fatalln("Error creating exercise :", err)
	}

	set, err := st.CreateSet(ctx, gen.CreateSetParams{
		ExerciseID: eid,
		Weight:     100,
	})

	if err != nil {
		log.Fatalln("Error updating exercise :", err)
	}

	set, err = st.UpdateSet(ctx, gen.UpdateSetParams{
		ExerciseID: eid,
		SetID:      set.SetID,
		Weight:     2000,
	})

	if err != nil {
		log.Fatalln("Error updating set :", err)
	}

	log.Println("Done!")

	u, err := st.ListUsers(ctx)

	for _, usr := range u {
		fmt.Println(fmt.Sprintf("Name : %s, ID : %d", usr.Name, usr.UserID))
	}
}
