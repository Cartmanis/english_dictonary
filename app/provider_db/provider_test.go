package provider_db

import (
	"fmt"
	"testing"
)

// You will be using this Trainer type later in the program
type Word struct {
	En            string
	Ru            string
	Transcription string
}

func TestMongoClient_InsertOne(t *testing.T) {
	client, err := NewStoreContext(10)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()
	word := &Word{
		En:            "late",
		Ru:            "поздний",
		Transcription: "лэйт",
	}
	id, err := client.InsertOne(word)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("id: ", id)
}
