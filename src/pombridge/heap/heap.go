package heap

import "fmt"

type Heapable interface {
	Priority() int
}

type Heap struct {
	max  int
	tree []Heapable
}

func New() *Heap {
	return &Heap{0, make([]Heapable, 50)}
}

func (h *Heap) String() string {
	var res string
	for i := 1; i <= h.max; i++ {
		res = fmt.Sprintf("%s %d", res, h.tree[i].Priority())
	}
	return fmt.Sprintf("[%s ]", res)
}

func (h *Heap) Push(v Heapable) {
	if h.max+1 >= len(h.tree) {
		h.resize()
	}
	h.max++
	h.tree[h.max] = v
	h.up(h.max)
}

func (h *Heap) Pop() Heapable {
	if h.max == 0 {
		return nil
	}

	min := h.tree[1] // root
	h.tree[1] = h.tree[h.max]
	h.tree[h.max] = nil
	h.max--
	h.down(1)

	return min
}

func (h *Heap) Top() Heapable {
	return h.tree[1]
}

func (h *Heap) Empty() bool {
	return h.max == 0
}

func (h *Heap) Len() int {
	return h.max
}

func (h *Heap) highPrioritySon(first, second int) int {
	if first > h.max || h.tree[first] == nil {
		return -1
	} else if h.tree[second] == nil {
		return first
	} else if h.tree[first].Priority() > h.tree[second].Priority() {
		return second
	}

	return first
}

func (h *Heap) up(pos int) {
	if pos == 1 {
		return
	}

	if h.tree[pos].Priority() < h.tree[pos/2].Priority() {
		h.swap(pos/2, pos)
		h.up(pos / 2)
	}
}

func (h *Heap) down(father int) {
	son := h.highPrioritySon(father*2, (father*2)+1)

	if son != -1 {
		if h.tree[father].Priority() > h.tree[son].Priority() {
			h.swap(father, son)
			h.down(son)
		}
	}
}

func (h *Heap) swap(i, j int) {
	h.tree[j], h.tree[i] = h.tree[i], h.tree[j]
}

func (h *Heap) resize() {
	h.tree = append(h.tree, make([]Heapable, len(h.tree))...)
}
