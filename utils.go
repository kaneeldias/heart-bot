package main

import "math/rand"

func GetRandomFromList(list []string) string {
	if len(list) == 0 {
		return ""
	}
	return list[rand.Intn(len(list))]
}
