package main

import (
	"fmt"
	"math/rand"
)

type Loader struct {
	svc *Service
}

func (loader *Loader) Run() error {
	for {
		usernum := rand.Intn(2000)
		err := loader.svc.Query(fmt.Sprintf("demo%d", usernum))
		if err != nil {
			return err
		}
	}
	return nil
}
