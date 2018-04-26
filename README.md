# gosqlseeder

A seeding package to help initialize data in the database!

## how to use

* 1. go get github.com/NtwariJoshua/gosqlseeder
* 2. import "github.com/NtwariJoshua/gosqlseeder"
## Example Code

```
package main

import(
  "github.com/NtwariJoshua/gosqlseeder"
  "database/sql"
)

var db *sql.DB
var err error

func init(){
 db,err = sql.Open("mysql","connection string")
 if err != nil {
		panic(err)
	}
}

func main(){
  defer db.Close()
  seeder := gosqlseeder.CreateSeeder(db)
  //create your seed data as map[string]interface{}
  data := map[string]{}{
    "FirstName":"John",
    "LastName":"Doe"
  }
  seed = gosqlseeder.CreateSeed("tableName",&data)
  _,err := seeder.Seed(&seed)
  if err != nil{
    panic(err)
  }

}
```
