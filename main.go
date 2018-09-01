package main

import (
	"github.com/larwef/ki/config"
	"github.com/larwef/ki/config/persistence"
	"log"
	"time"
)

func main() {
	log.Println("Starting...")

	conf := config.Config{
		Id:           "someId",
		Name:         "someOtherName",
		LastModified: time.Now(),
		Properties:   []byte(`{"num":6.13,"strs":["a","b"]}`),
	}

	local := persistence.NewLocal("test/")

	err := local.Store(conf)
	if err != nil {
		log.Fatal(err)
	}

	conf2, err := local.Retrieve("someId")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(conf2)

	log.Println("Exiting application.")
}
