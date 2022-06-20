package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const startload = 100 // initial query per seconds
const loadperiod = 5  // number of seconds to preserve current load
const increment = 100 // the additional number of queries to add to the previous load

type Loader struct {
	svc   *Service
	reqps int // requests per seconds ("load")
}

func (loader *Loader) Run() error {
	loader.reqps = startload
	curperiod := 0
	for {
		gstart := time.Now()
		for i := 0; i < 10; i++ {
			//curstart := time.Now()
			//var wg sync.WaitGroup
			for j := 1; j < loader.reqps/10; j++ {
				//wg.Add(1)
				go func() {
					//		defer wg.Done()
					usernum := rand.Intn(2000)
					start := time.Now()
					err := loader.svc.Query(fmt.Sprintf("demo%d", usernum))
					elapsed := time.Since(start)
					if elapsed.Milliseconds() > int64(1000) {
						log.Fatalf("query exeeded 1 second at %d queries per second", loader.reqps+j)
					}
					if err != nil {
						log.Fatalf("query returned error at %d\n %v", j, err)
					}
					time.Sleep(10 * time.Millisecond)
				}()
			}
			//wg.Wait()
			//curelapsed := time.Since(curstart)
			//log.Printf("elapsed %s at %d chunk %d", curelapsed, loader.reqps, i)
		}
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
