package csv

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncoder(t *testing.T) {
	focus := os.Getenv("FOCUS")

	examples := []struct {
		Name   string
		Input  interface{}
		Output string
	}{
		{
			Name:  "slice-string",
			Input: []string{"hello", "there"},
		},
		{
			Name:  "slice-interface",
			Input: []interface{}{"hello", "there", 123, 10.4},
		},
		{
			Name: "simple-struct",
			Input: struct {
				Name     string
				Greeting string
				Score    int `render:"points,omitempty"`
				boop     bool
			}{"Obi Wan", "Hello there", 9001, false},
		},
		{
			Name: "slice-simple-struct",
			Input: []struct {
				Name     string
				Greeting string
				Score    int `render:"points,omitempty"`
				boop     bool
			}{
				{Name: "Obi Wan", Greeting: "Hello there", Score: 9001},
				{Name: "Obi Wan", Greeting: "Hello there"},
			},
		},
		{
			Name: "simple-embedded-struct",
			Input: struct {
				Name     string
				Greeting string
				IgnoreMe string `render:"-"`
				Score    int    `render:"points,omitempty"`
				Rank     struct {
					Class string
					Title string
				} `render:",noprefix"`
			}{
				Name:     "Obi Wan",
				Greeting: "Hello there",
				IgnoreMe: "please",
				Score:    9001,
				Rank: struct {
					Class string
					Title string
				}{
					Title: "Master",
					Class: "Jedi",
				},
			},
		},
		{
			Name: "struct-with-slice",
			Input: struct {
				Name     string
				Greeting string
				Planets  []string
				Scores   []int
				Grades   []float32
			}{
				Name:     "Obi Wan",
				Greeting: "Hello there",
				Planets:  []string{"tatooine", "naboo"},
				Scores:   []int{9001, 1337, 1138},
				Grades:   []float32{12.2, 3.14, 8.67},
			},
		},
		{
			Name: "embedded-struct-with-slice",
			Input: struct {
				Name string
				Rank struct {
					Class   string
					Title   string
					Aliases []string
				} `render:",noprefix"`
			}{
				Name: "Obi Wan",
				Rank: struct {
					Class   string
					Title   string
					Aliases []string
				}{
					Title:   "Master",
					Class:   "Jedi",
					Aliases: []string{"ben", "obi-wan", "master"},
				},
			},
		},
		{
			Name: "simple-map",
			Input: map[string]interface{}{
				"name":     "kenobi",
				"greeting": "Hello There",
			},
		},
		{
			Name: "slice-map",
			Input: []map[string]interface{}{
				{
					"Name":     "Obi-Wan",
					"Greeting": "Hello There",
				},
				{
					"Name":     "Mark",
					"Greeting": "Oh hi",
				},
			},
		},
		{
			Name: "slice-map-missing-keys",
			Input: []map[string]interface{}{
				{
					"Name": "Obi-Wan",
				},
				{
					"Greeting": "Oh hi",
				},
			},
		},
	}

	for _, example := range examples {
		example := example
		t.Run(example.Name, func(t *testing.T) {
			if focus != "" && example.Name != focus {
				t.Skipf("example '%s' is out of focus", example.Name)
			}
			enc := &Encoder{}
			err := enc.Encode(example.Input)
			if err != nil {
				t.Errorf("unexpected error '%s'", err)
			}

			expected, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", example.Name))
			if err != nil {
				t.Fatalf("could not read expected testdata file: 'testdata/%s'", example.Name)
			}
			if !bytes.Equal(enc.Bytes(), expected) {
				t.Errorf("output is incorrect\nactual:\n%s\n\nexpected:\n%s", string(enc.Bytes()), string(expected))
			}
		})
	}
}
