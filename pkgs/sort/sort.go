package sort

import (
	"hypermedlab/backend-myblog/models/blog"
	"sort"
)

type blogs []*blog.Blog

func (bs blogs) Len() int {
	return len(bs)
}

func (bs blogs) Less(i, j int) bool {
	return bs[i].CreateAt.After(bs[j].CreateAt)
}

func (bs blogs) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

func Sort(bs blogs) {
	sort.Sort(bs)
}
