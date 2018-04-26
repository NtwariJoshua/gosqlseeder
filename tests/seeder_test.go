package tests

import (
	"testing"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/NtwariJoshua/gosqlseeder"
	"fmt"
)

var db *sql.DB
var err error
var testSeeder *gosqlseeder.Seeder
var testSeed *gosqlseeder.Seed
func init() {
	db, err = sql.Open("sqlite3", "./testdatabase")
	if err != nil {
		panic(err)
	}
}
func TestSeederCreation(t *testing.T) {
	testSeeder = gosqlseeder.CreateSeeder(db)
	if testSeeder == nil {
		t.Error(fmt.Sprintf("Test Failed the return value is not a Seeder Object: %v", err))
	}
}

func TestSeedCreation(t *testing.T){
	testSeed = gosqlseeder.CreateSeed("foo",map[string]interface{}{
		"name":"John Doe",
		"email":"john@doe.com",
	},map[string]interface{}{
		"name":"Jane Doe",
		"email":"jane@doe.com",
	})
	if testSeed == nil{
		t.Error(fmt.Sprintf("Test Failed the return value is not a Seed Object: %v", err))
	}
}

func TestSeedingProcess(t *testing.T)  {
	defer db.Close()
	_,err := testSeeder.Seed(testSeed)
	if err != nil{
		t.Error(fmt.Sprintf("Test Failed Seeding Process Failed: %v", err))
	}
}
