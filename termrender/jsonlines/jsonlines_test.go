package jsonlines

import "fmt"

func ExampleMarshal() {
	type Address struct {
		Street string
		City   string
		State  string
		Zip    int `json:"ZipCode"`
	}

	type Customer struct {
		Name    string
		Age     int `json:",omitempty"`
		Address Address
	}

	customers := []Customer{
		{
			Name: "Obi-Wan Kenobi",
			Address: Address{
				Street: "2000 Ultimate Way",
				City:   "Weston",
				State:  "FL",
				Zip:    33326,
			},
		},
		{
			Name: "Chewbacca",
			Age:  200,
			Address: Address{
				Street: "2000 Ultimate Way",
				City:   "Weston",
				State:  "FL",
				Zip:    33326,
			},
		},
	}

	out, _ := Marshal(customers)

	fmt.Println(string(out))

	// Output:
	// {"Name":"Obi-Wan Kenobi","Address":{"Street":"2000 Ultimate Way","City":"Weston","State":"FL","ZipCode":33326}}
	// {"Name":"Chewbacca","Age":200,"Address":{"Street":"2000 Ultimate Way","City":"Weston","State":"FL","ZipCode":33326}}

}
