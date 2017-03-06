package vipcomm

type Int64Set struct {
	m map[int64]bool
}

func NewInt64Set() *Int64Set {
	return &Int64Set{
		m: map[int64]bool{},
	}
}

func (s *Int64Set) Add(item int64) {
	s.m[item] = true
}

func (s *Int64Set) FromList(l []int64) {
	for _, i := range l {
		s.Add(i)
	}
}

func (s *Int64Set) List() (l []int64) {
	for k := range s.m {
		l = append(l, k)
	}
	return
}

func (s *Int64Set) Has(item int64) bool {
	_, ok := s.m[item]
	return ok
}

// 是否包含列表
func (s *Int64Set) CheckSub(l []int64) bool {
	for _, i := range l {
		if _, ok := s.m[i]; !ok {
			return false
		}
	}
	return true
}

// 交集
func (s *Int64Set) Intersection(l []int64) []int64 {
	newl := make([]int64, 0)
	for _, i := range l {
		if _, ok := s.m[i]; ok {
			newl = append(newl, i)
		}
	}
	return newl
}

// 去重
func Int64RmDup(s []int64) {
	if len(s) == 0 {
		return
	}

	set := NewInt64Set()
	set.FromList(s)
	s = set.List()
}
