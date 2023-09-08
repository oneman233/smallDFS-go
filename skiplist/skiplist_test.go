package skiplist

import (
	"testing"
)

func TestSkipList_Get(t *testing.T) {
	list := NewSkipList(5)
	list.Put("1", 1)
	list.Put("1", 2)
	list.Put("3", 3)
	list.Show()
	if v, err := list.Get("1"); v.(int) != 2 || err != nil {
		t.Errorf("wrong value %d", v.(int))
	}
	if v, err := list.Get("2"); err == nil || v != nil {
		t.Error("key not existed")
	}
}

func TestSkipList_Delete(t *testing.T) {
	list := NewSkipList(5)
	list.Put("1", 1)
	list.Put("1", 2)
	list.Put("3", 3)
	list.Put("5", 5)
	list.Delete("4")
	list.Delete("3")
	if v, err := list.Get("3"); v != nil || err == nil {
		t.Error("delete 3 failed")
	}
	if v, err := list.Get("5"); v == nil || err != nil {
		t.Error("get 5 failed")
	}
}

func TestSkipList_LowerBound(t *testing.T) {
	list := NewSkipList(5)
	list.Put("1", 1)
	list.Put("1", 2)
	list.Put("3", 3)
	list.Put("5", 5)
	list.Delete("4")
	list.Delete("3")
	list.Put("7", 7)
	// 1:2 5:5 7:7
	v, err := list.Get(list.LowerBound("6"))
	if v.(int) != 5 || err != nil {
		t.Errorf("lower bound 6 get %d", v)
	}
	v, err = list.Get(list.LowerBound("5"))
	if v.(int) != 5 || err != nil {
		t.Errorf("lowe bound 5 get %d", v)
	}
	v, err = list.Get(list.LowerBound("4"))
	if v.(int) != 2 || err != nil {
		t.Errorf("lower bound 4 get %d", v)
	}

}
