package main

import "strings"

type Rgxs []string

func (r *Rgxs) String() string {
	return strings.Join([]string(*r), " ")
}

func (r *Rgxs) Set(value string) error {

	*r = append(*r, value)

	return nil
}
