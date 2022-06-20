package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const startload = 100 // initial concurrent queries
const loadperiod = 5  // number of seconds to preserve current load
const increment = 100 // the additional number of queries to add to the previous load
const timeout = 1000  // in milliseconds, if a request exceeds this, exit the test

type Loader struct {
	svc   *Service
	reqps int // requests per seconds ("load")
}

func (loader *Loader) Run() error {
	loader.reqps = startload
	last_increment := time.Now()
	for {
		curstart := time.Now()
		for i := 0; i < 10; i++ {
			for j := 1; j < loader.reqps/10; j++ {
				//wg.Add(1)
				time.Sleep(10 * time.Millisecond)
				go func() {
					//		defer wg.Done()
					usernum := rand.Intn(2000)
					start := time.Now()
					err := loader.svc.Query(fmt.Sprintf("demo%d", usernum))
					elapsed := time.Since(start)
					if elapsed.Milliseconds() > int64(timeout) {
						log.Fatalf("query exeeded %d milliseconds at %d queries per second",
							timeout, loader.reqps+i*10+j)
					}
					if err != nil {
						log.Fatalf("query returned error at %d\n %v", loader.reqps+i*10+j, err)
					}
				}()
			}
		}
		curelapsed := time.Since(curstart)
		log.Printf("%d requests took %s", loader.reqps, curelapsed)

		incr_elapsed := time.Since(last_increment)
		log.Printf("%s since last increment", incr_elapsed)

		if incr_elapsed.Seconds() > loadperiod {
			last_increment = time.Now()
			loader.reqps += increment
		}
	}
	return nil
}
