package database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	//. "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/dalpengida/portfolio-api-go/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	// Import this so we don't have to use qm.Limit etc.
)

var ctx = context.Background()

// username:password@protocol(address)/dbname?param=value
func dsn(host, dbName, username, password string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, dbName)
}

func Test(t *testing.T) {
	dsn := dsn("localhost:3306", "portfolio", "portfolio", "test")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}

	boil.SetDB(db)

	exist, err := models.AccountExists(ctx, db, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)

	account, err := models.Accounts().One(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(account)

	accounts, err := models.Accounts().All(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%v", accounts)

	// account := models.Account(
	// 	Select("user_id", "username"),
	// 	Where("user_id = ?", 1),
	// 	Limit(1),
	// ).All(ctx, db)

	// err = account.Insert(ctx, db)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// fmt.Println(account)

}

// func AA() error {
// 	dsn := dsn("localhost:3306", "portfolio", "portfolio", "test")
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return err
// 	}

// 	// boil.SetDB(db)
// 	// users, err := models.Acc().AllG(ctx)

// 	// // Query all users
// 	// users, err := models.Users().All(ctx, db)

// 	// // Panic-able if you like to code that way (--add-panic-variants to enable)
// 	// users := models.Users().AllP(db)

// 	// // More complex query
// 	// users, err := models.Users(Where("age > ?", 30), Limit(5), Offset(6)).All(ctx, db)

// 	// // Ultra complex query
// 	// users, err := models.Users(
// 	// 	Select("id", "name"),
// 	// 	InnerJoin("credit_cards c on c.user_id = users.id"),
// 	// 	Where("age > ?", 30),
// 	// 	AndIn("c.kind in ?", "visa", "mastercard"),
// 	// 	Or("email like ?", `%aol.com%`),
// 	// 	GroupBy("id", "name"),
// 	// 	Having("count(c.id) > ?", 2),
// 	// 	Limit(5),
// 	// 	Offset(6),
// 	// ).All(ctx, db)

// 	// // Use any "boil.Executor" implementation (*sql.DB, *sql.Tx, data-dog mock db)
// 	// // for any query.
// 	// tx, err := db.BeginTx(ctx, nil)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// users, err := models.Users().All(ctx, tx)

// 	// // Relationships
// 	// user, err := models.Users().One(ctx, db)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// movies, err := user.FavoriteMovies().All(ctx, db)

// 	// // Eager loading
// 	// users, err := models.Users(Load("FavoriteMovies")).All(ctx, db)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// fmt.Println(len(users.R.FavoriteMovies))

// }
