package main

import (
	"fmt"
	"log"
	"math"
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
	sleep_nano := 10000000 // nano seconds to sleep before each request
	log.Printf("reqps=%d, sleep_nano=%d", loader.reqps, sleep_nano)
	speedup_factor := 1.1
	for {
		curstart := time.Now()
		for i := 0; i < 10; i++ {
			for j := 1; j < loader.reqps/10; j++ {
				time.Sleep(time.Duration(sleep_nano) * time.Nanosecond)

				go func() {
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
		log.Printf("%d requests took %s (%.1f per second)",
			loader.reqps,
			curelapsed,
			math.Round(float64(loader.reqps)/curelapsed.Seconds()))

		incr_elapsed := time.Since(last_increment)

		if curelapsed.Seconds() > 1 {
			log.Printf("PREV speed-up factor %.1f", speedup_factor)
			speedup_factor *= 1.1
			log.Printf("speed-up factor %.1f", speedup_factor)
		}

		if curelapsed.Seconds() < 1 {
			log.Printf("PREV speed-up factor %.1f", speedup_factor)
			speedup_factor *= 0.9
			log.Printf("speed-up factor %.1f", speedup_factor)
		}
		sleep_nano = int(math.Round(float64(1000000000/loader.reqps) / speedup_factor))

		if incr_elapsed.Seconds() > loadperiod {
			last_increment = time.Now()
			loader.reqps += increment
		}
	}
	return nil
}
