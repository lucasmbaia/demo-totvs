package request

import (
	"testing"
	"fmt"
)

type DataTeste struct {
	Firewall string	`json:"Firewall"`
}

func Test_Request(t *testing.T) {
	c, _ := NewClient()
	var data []DataTeste

	if err := c.Get("demo-nats-000001", &data, nil); err != nil {
		t.Fatal(err)
	}

	fmt.Println(data)
}
