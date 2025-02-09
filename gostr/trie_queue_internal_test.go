package gostr

import "testing"

func TestTrieQueue(t *testing.T) { //nolint:funlen // test functions can be long, damned you!
	// this is just dummy data...
	var (
		n1    = new(Trie)
		n2    = new(Trie)
		n3    = new(Trie)
		n4    = new(Trie)
		n5    = new(Trie)
		n6    = new(Trie)
		n7    = new(Trie)
		queue = newTrieQueue(0)
	)

	if !queue.isEmpty() {
		t.Fatal("queue should be empty")
	}

	if queue.isFull() {
		t.Fatal("queue should not be full")
	}

	queue.enqueue(n1)

	if !queue.isFull() {
		t.Fatal("queue should be full now")
	}

	queue.enqueue(n2)

	if len(queue.elms) != 2 {
		t.Errorf("we expected the queue to have two elements, but it has %d", len(queue.elms))
	}

	if queue.used != 2 {
		t.Errorf("we expected the queue to have two elements, but it has %d", queue.used)
	}

	if !queue.isFull() {
		t.Fatal("queue should be full again (now with cap 2)")
	}

	queue.enqueue(n3)

	if queue.used != 3 {
		t.Errorf("The queue should have three elments, but it has %d", queue.used)
	}

	if len(queue.elms) != 4 {
		t.Errorf("expected cap on 5 but it is %d", len(queue.elms))
	}

	var n *Trie

	n = queue.dequeue()
	if n != n1 {
		t.Errorf("unexpected dequeue: %p", n)
	}

	if queue.used != 2 {
		t.Errorf("The queue should have two elments, but it has %d", queue.used)
	}

	n = queue.dequeue()
	if n != n2 {
		t.Errorf("unexpected dequeue: %p", n)
	}

	if queue.used != 1 {
		t.Errorf("The queue should have one elment, but it has %d", queue.used)
	}

	// push it up to capacity again... cap is 4 and we have used 1, so we need to add 4
	queue.enqueue(n4)
	queue.enqueue(n5)
	queue.enqueue(n6)
	queue.enqueue(n7)

	n = queue.dequeue()
	if n != n3 {
		t.Errorf("unexpected dequeue: %p", n)
	}

	n = queue.dequeue()
	if n != n4 {
		t.Errorf("unexpected dequeue: %p", n)
	}

	n = queue.dequeue()
	if n != n5 {
		t.Errorf("unexpected dequeue: %p", n)
	}

	n = queue.dequeue()
	if n != n6 {
		t.Errorf("unexpected dequeue: %p", n)
	}

	n = queue.dequeue()
	if n != n7 {
		t.Errorf("unexpected dequeue: %p", n)
	}

	if queue.used != 0 {
		t.Errorf("The queue should have zero elments, but it has %d", queue.used)
	}

	if !queue.isEmpty() {
		t.Error("The queue should be empty, but it isn't.")
	}
}
