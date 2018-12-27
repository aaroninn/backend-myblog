package sort

import (
	"hypermedlab/myblog/models/blog"
	"sort"
	"time"
)

type blogs []*blog.Blog

func (bs blogs) Len() int {
	return len(bs)
}

func (bs blogs) Less(i, j int) bool {
	t1, err := time.Parse("2006-01-02 15:04:05", bs[i].CreateAt.String())
	if err != nil {
		return false
	}
	t2, err := time.Parse("2006-01-02 15:04:05", bs[j].CreateAt.String())
	if err != nil {
		return false
	}

	return t1.Before(t2)
}

func (bs blogs) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

func Sort(bs blogs) {
	sort.Sort(bs)
}
