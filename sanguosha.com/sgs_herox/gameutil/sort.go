package gameutil

type UInt64KV struct {
	Key   uint64
	Value uint64
}
type SortUInt64KV struct {
	List   []UInt64KV
	IsLess bool //ć°çćšć
}

func (ts SortUInt64KV) Len() int {
	return len(ts.List)
}
func (ts SortUInt64KV) Swap(i, j int) {
	ts.List[i], ts.List[j] = ts.List[j], ts.List[i]
}

func (ts SortUInt64KV) Less(i, j int) bool {
	if ts.IsLess {
		return ts.List[i].Value < ts.List[j].Value
	}
	return ts.List[i].Value > ts.List[j].Value
}
