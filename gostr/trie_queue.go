package gostr

type trieQueue struct {
	used, front int
	elms        []*Trie
}

func newTrieQueue(capacity int) *trieQueue {
	if capacity < 1 {
		capacity = 1 // never less than one, or the growing won't work
	}

	return &trieQueue{
		used: 0, front: 0,
		elms: make([]*Trie, capacity),
	}
}

func (q *trieQueue) isEmpty() bool {
	return q.used == 0
}

func (q *trieQueue) isFull() bool {
	return q.used == len(q.elms)
}

// Only call this when used=cap!
func (q *trieQueue) grow(newCap int) {
	newElms := make([]*Trie, newCap)
	n := 0

	for i := q.front; i < len(q.elms); i++ {
		newElms[n] = q.elms[i]
		n++
	}

	for i := 0; i < q.front; i++ {
		newElms[n] = q.elms[i]
		n++
	}

	q.elms = newElms
	q.front = 0
}

func (q *trieQueue) enqueue(t *Trie) {
	if q.isFull() {
		q.grow(2 * len(q.elms)) //nolint:gomnd // doubling sizes
	}

	q.elms[(q.front+q.used)%len(q.elms)] = t
	q.used++
}

func (q *trieQueue) dequeue() *Trie {
	t := q.elms[q.front]
	q.used--
	q.front = (q.front + 1) % len(q.elms)

	return t
}
