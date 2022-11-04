package slr

import (
	"testing"
	"fmt"
)

func Test_GetTTable(t *testing.T) {
	if value, err := GetTTable("5%", "8"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(value)
	}
}

func Test_GetFTable(t *testing.T) {
	if value, err := GetFTable("1", "8"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(value)
	}
}
