package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const startload = 100      // initial query per seconds
const loadperiod = 5       // number of seconds to preserve current load
const increment = 100      // the additional number of queries to add to the previous load
const secondpart = 1 / 100 // hundredth of second

type Loader struct {
	svc   *Service
	reqps int // requests per seconds ("load")
}

func (loader *Loader) Run() error {
	loader.reqps = startload
	curperiod := 0
	for {
		gstart := time.Now()
		var wg sync.WaitGroup
		for i := 1; i < loader.reqps; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				usernum := rand.Intn(2000)
				start := time.Now()
				err := loader.svc.Query(fmt.Sprintf("demo%d", usernum))
				elapsed := time.Since(start)
				if elapsed.Milliseconds() > int64(1000) {
					log.Fatalf("query exeeded 1 second at %d queries per second", i)
				}
				if err != nil {
					log.Fatalf("query returned error at %d\n %v", i, err)
				}
			}()
		}
		wg.Wait()
		gelapsed := time.Since(gstart)
		log.Printf("%d requests took %s", loader.reqps, gelapsed)

		curperiod += 1
		if curperiod == loadperiod {
			curperiod = 0
			loader.reqps += increment
		}
	}
	return nil
}
