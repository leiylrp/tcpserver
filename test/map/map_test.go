package main

import "testing"

func TestMapShareMemory(t *testing.T) {
	m1 := map[string]int{}
	m2 := m1
	m1["a"] = 1
	t.Log(m1, len(m1))
	// map[a:1] 1
	t.Log(m2, len(m2))
	// map[a:1]
}


func TestSliceFn(t *testing.T) {
	m := map[string]int{}
	t.Log(m, len(m))
	// map[a:1]
	m["a"] = 1
	t.Log(m, len(m))
	mapAppend(m, "b", 2)
	t.Log(m, len(m))
	// map[a:1 b:2] 2
}

func mapAppend(m map[string]int, key string, val int) {
	m[key] = val
}

