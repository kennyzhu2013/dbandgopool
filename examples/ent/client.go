package main

import (
	"context"
	"fmt"
	"github.com/panjf2000/ants/v2/examples/ent/ent"
	"github.com/panjf2000/ants/v2/examples/ent/ent/car"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 解决golang：unsupported Scan, storing driver.Value type []uint8 into type *time.Time: &parseTime参数.
	client, err := ent.Open("mysql", "root:L0v1!@#$@tcp(10.153.90.12:3306)/mid_bussops_prod?charset=utf8&parseTime=true")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	ctx := context.Background()
	user,_ := CreateCars(ctx, client)
	if err = QueryCars(ctx, user); err != nil {
		log.Printf("failed QueryCars: %v, ctx:%v\n", err, ctx.Err())
	}

	_ = QueryCarUsers(ctx, user)
}

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	// Create a new car with model "Tesla".
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", tesla)

	// Create a new car with model "Ford".
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", ford)

	// Create a new user, and add it the 2 cars.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)
	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed QueryCars all: %w", err)
	}
	log.Println("returned cars:", cars)

	// What about filtering specific cars.
	ford, err := a8m.QueryCars().
		Where(car.Model("Ford")).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println(ford)
	return nil
}

func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	// Query the inverse edge.
	for _, ca := range cars {
		owner, err := ca.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q owner: %w", ca.Model, err)
		}
		log.Printf("car %q owner: %q\n", ca.Model, owner.Name)
	}
	return nil
}
//
//func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
//	u, err := client.User.
//		Create().
//		SetAge(30).
//		SetName("a8m").
//		Save(ctx)
//	if err != nil {
//		return nil, fmt.Errorf("failed creating user: %w", err)
//	}
//	log.Println("user was created: ", u)
//	return u, nil
//}
//
//func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
//	u, err := client.User.
//		Query().
//		Where(user.Name("a8m")).
//		// `Only` fails if no user found,
//		// or more than 1 user returned.
//		Only(ctx)
//	if err != nil {
//		return nil, fmt.Errorf("failed querying user: %w", err)
//	}
//	log.Println("user returned: ", u)
//	return u, nil
//}
