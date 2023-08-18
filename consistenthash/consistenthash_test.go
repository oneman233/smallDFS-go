package consistenthash

import (
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	hash := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	// 对应虚拟节点是：20 21 22 40 41 42 60 61 62
	hash.Add("4", "6", "2")

	testCases := map[string]string{
		"2":  "2",
		"35": "4",
		"41": "4",
		"18": "2",
		"59": "6",
		"62": "6",
		"63": "2",
		"80": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, expect %s, get %s", k, v, hash.Get(k))
		}
	}

	// 对应的虚拟节点是：80 81 82
	hash.Add("8")

	if hash.Get("80") != "8" {
		t.Errorf("Add 8 failed")
	}
}
