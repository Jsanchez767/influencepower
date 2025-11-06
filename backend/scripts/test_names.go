package main

import (
	"fmt"
	"log"

	"github.com/Jsanchez767/InfluencePower/backend/cityapi"
)

func main() {
	client := cityapi.NewClient()
	
	persons, err := client.GetPersons()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Total persons: %d\n\n", len(persons))
	fmt.Println("First 20 names:")
	for i, p := range persons {
		if i < 20 {
			fmt.Printf("%d. %s\n", i+1, p.PersonFullName)
		}
	}
}
