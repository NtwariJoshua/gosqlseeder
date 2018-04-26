package gosqlseeder

import (
	"database/sql"
	"fmt"
	"time"
	"math/rand"

	"sync"
)

type SandClock struct {
	counter int
}
type Seeder struct {
	db       *sql.DB
	seedData []*Seed
}

type Seed struct {
	tableName string
	seedData  []map[string]interface{}
}

type QResult struct {
	status bool
	err    error
	query  string
}

var wg sync.WaitGroup
var status bool = true
var errors error = nil

func CreateSeed(table string, seeds ...map[string]interface{}) *Seed {
	if len(seeds) == 0 {
		panic("Can not create seed with no data")
	}
	return &Seed{
		tableName: table,
		seedData:  seeds,
	}

}
func CreateSeeder(database *sql.DB) *Seeder {
	var seeder = Seeder{
		db: database,
	}
	return &seeder
}

func (s *Seeder) Seed(seeds ...*Seed) (bool, error) {
	if len(seeds) == 0 {
		panic("Can not use empty map as seed")
	}
	s.seedData = seeds
	var status bool
	var err error
	if status, err = s.startSeeding(); err != nil {
		fmt.Println("Error Occured during while seeding:", err)
	}
	return status, err
}

func (s *Seeder) startSeeding() (bool, error) {
	resultsChannel := make(chan QResult)
	completed := make(chan bool)

	for _, seedSlice := range s.seedData {
		sc := &SandClock{}
		oneLineQ := ""
		var querys []string
		for _, seedMap := range seedSlice.seedData {
			mainquery := fmt.Sprintf("INSERT into %v ", seedSlice.tableName)
			queryPortion1 := "("
			queryPortion2 := fmt.Sprintf("Values(")
			totalSize := len(seedMap)
			sc.counter = 0
			for k, v := range seedMap {
				if sc.counter+1 >= totalSize {
					queryPortion1 = fmt.Sprintf(queryPortion1+"%v)", k)
					queryPortion2 = fmt.Sprintf(queryPortion2+"'%v')", v)
				} else {
					queryPortion1 = fmt.Sprintf(queryPortion1+"%v,", k)
					queryPortion2 = fmt.Sprintf(queryPortion2+"'%v',", v)
				}

				sc.counter += 1
			}
			oneLineQ = mainquery+queryPortion1 + queryPortion2 +";"
			querys = append(querys, oneLineQ)

		}

		for _, q := range querys {
			wg.Add(1)
			fmt.Println("Executing: ", q)
			executeQueries(s, q, time.Duration(rand.Int31n(10)), resultsChannel)
		}

	}
	go func() {
		for {
			select {
			case rtn := <-resultsChannel:
				if rtn.status {
					fmt.Println("Successful: ", rtn.query)
				} else {
					fmt.Println(fmt.Sprintf("Error Occured: %v on query: %v", rtn.err, rtn.query))
					status = false
					errors = rtn.err
				}
			case done := <-completed:
				if done {
					break
				}

			default:
				break
			}
		}

	}()

	wg.Wait()
	completed <- true

	return status, errors
}

func executeQueries(s *Seeder, q string, t time.Duration, writeTo chan QResult) {

	go func(q string, conn *sql.DB) {
		defer wg.Done()
		if _, err := s.db.Exec(q); err != nil {
			writeTo <- QResult{false, err, q}
		} else {
			writeTo <- QResult{true, nil, q}
		}
	}(q, s.db)
}
