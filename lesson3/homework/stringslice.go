package main

import "strings"

type StringSlice []string

func (ss *StringSlice) String() string {
	return strings.Join(*ss, ", ")
}

func (ss *StringSlice) Set(value string) error {
	*ss = strings.Split(value, ",")
	return nil
}
