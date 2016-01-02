package main

import "strings"

type Rgxs []string

func (r *Rgxs) String() string {
	return strings.Join([]string(*r), " ")
}

func (r *Rgxs) Set(value string) error {
	// If we wanted to allow the flag to be set multiple times,
	// accumulating values, we would delete this if statement.
	// That would permit usages such as
	//	-deltaT 10s -deltaT 15s
	// and other combinations.
	//	if len(*r) > 0 {
	//		return errors.New("interval flag already set")
	//	}

	*r = append(*r, value)

	return nil
}
