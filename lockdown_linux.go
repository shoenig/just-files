//go:build linux

package main

import (
	"github.com/shoenig/go-landlock"
)

func lockdown(m map[string]string) {
	var paths []*landlock.Path
	for _, v := range m {
		paths = append(paths, landlock.Dir(v, "r"))
	}
	ll := landlock.New(paths...)
	err := ll.Lock(landlock.Mandatory)
	check(err)
}
