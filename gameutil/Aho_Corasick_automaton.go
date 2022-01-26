package gameutil

// AC自动机 字符串多模 匹配
// 接收一个字符流，并检查这些字符的后缀是否是字符串数组 words 中的一个字符串

// Aho–Corasick Algorithm
type StreamChecker struct {
	t           *Trie
	nameToPhone map[string]string
}

func Constructor(words []string) StreamChecker {
	t := &Trie{}

	for _, w := range words {
		t.Insert(w)
	}
	t.Build()

	return StreamChecker{t: t}
}
func (p *StreamChecker) SetMap(nToP map[string]string) {
	p.nameToPhone = nToP
}

func (p *StreamChecker) Query(s string) []string {
	var ret []string
	now := p.t
	for _, v := range s {
		if v < 'a' || v > 'z' {
			continue
		}
		//if now.child[v-'a'] != nil {
		//	now = now.child[v-'a']
		//} else {
		//now = now.fail
		now = now.child[v-'a']
		//}
		if now.isEnd {
			for _, name := range now.name {
				ret = append(ret, p.nameToPhone[name])
			}
		}
	}
	return ret
}

type Trie struct {
	child [26]*Trie
	isEnd bool
	name  []string
	fail  *Trie
}

func (t *Trie) Insert(word string) {
	for i := range word {
		if word[i] < 'a' || word[i] > 'z' {
			continue
		}
		idx := word[i] - 'a'
		if t.child[idx] == nil {
			t.child[idx] = &Trie{}
		}
		t = t.child[idx]
	}
	t.isEnd = true
	t.name = append(t.name, word)
}

func (t *Trie) Build() {
	t.fail = t
	q := []*Trie{}

	for i := 0; i < 26; i++ {
		if t.child[i] != nil {
			t.child[i].fail = t
			q = append(q, t.child[i])
		} else {
			t.child[i] = t
		}
	}

	for len(q) > 0 {
		t := q[0]
		q = q[1:]
		t.isEnd = t.isEnd || t.fail.isEnd
		for i := 0; i < 26; i++ {
			if t.child[i] != nil {
				t.child[i].fail = t.fail.child[i]
				q = append(q, t.child[i])
			} else {
				t.child[i] = t.fail.child[i]
			}
		}
	}
}
