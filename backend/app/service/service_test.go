package service

import (
	"fmt"
	"testing"
)

func TestGetIdString(t *testing.T) {
	id, err := GetIdString(`ObjectID("5d05cda973f46b9cc4e2107b)`)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(id)
}
