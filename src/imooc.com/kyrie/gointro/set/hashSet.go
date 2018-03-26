package set

import (
	"bytes"
	"fmt"
)

type HashSet struct {
	m map[interface{}] bool
}

func NewHashSet() *HashSet {
	return &HashSet{m: make(map[interface{}] bool)}
}

func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e)
}
// 清空的话，直接赋值为空，之前的值等GC回收
func (set *HashSet) Clear() {
	set.m = make(map[interface{}] bool)
}

func (set *HashSet) Contains(e interface{}) bool  {
	return set.m[e]
}

func (set *HashSet) Len() int {
	return len(set.m)
}

func (set *HashSet) Same(other Set) bool {
	if other == nil {
		return false
	}
	if set.Len() != other.Len() {
		return false
	}
	for key := range set.m {
		if !other.Contains(key) {
			return false
		}
	}

	return true
}

// 迭代过程中m的数量可能增加也可能减少
// todo m的值不是并发安全的，还需要加上互斥量
func (set *HashSet) Elements() []interface{}  {
	initialLen := len(set.m)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for key := range set.m {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else {
			snapshot = append(snapshot, key)
		}
		actualLen++
	}

	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}

	return snapshot
}

//使用bytes.Buffer类型值作为结果值的缓冲区，这样可以避免因string类型值的拼接造成的内存空间
func (set *HashSet) String() string {
	var buf bytes.Buffer
	buf.WriteString("HashSet{")
	first := true
	for key := range set.m {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	buf.WriteString("}")
	return buf.String()
}


