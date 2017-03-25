package bitarray

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestErr(t *testing.T) {
	fmt.Println("Test: Error")

	bits := New(64)
	idxs := []int{-1, -2, -3444, 65, 66, 67, 657}

	for _, idx := range idxs {

		if _, err := bits.Get(idx); err == nil {
			t.Fatalf("expect Get index %d out of range\n", idx)
		}

		if _, err := bits.Put(idx, 1); err == nil {
			t.Fatalf("expect Put index %d out of range\n", idx)
		}

		if err := bits.Set(idx, idx+1); err == nil {
			t.Fatalf("expect Set index %d out of range\n", idx)
		}

		if err := bits.Clear(idx, idx+1); err == nil {
			t.Fatalf("expect Clear index %d out of range\n", idx)
		}

		if err := bits.Not(idx, idx+1); err == nil {
			t.Fatalf("expect Not index %d out of range\n", idx)
		}
	}

	if err := bits.Set(13, 12); err == nil {
		t.Fatalf("Set(13, 12): low should less than high\n")
	}

	if err := bits.Clear(13, 12); err == nil {
		t.Fatalf("Set(13, 12): low should less than high\n")
	}

	if err := bits.Not(13, 12); err == nil {
		t.Fatalf("Set(13, 12): low should less than high\n")
	}
}

func TestLen(t *testing.T) {
	fmt.Println("Test: Len")

	bits := New(3)
	if n := bits.Len(); n != 3 {
		t.Fatalf("expect length: %d, got : %d\n", 3, n)
	}
}

func TestBasic(t *testing.T) {
	fmt.Println("Test: Get/Put")

	bits := New(65)

	for i := 0; i < bits.Len(); i++ {
		if bit, _ := bits.Get(i); bit != 0 {
			t.Fatalf("expect bit 0")
		}
		bits.Put(i, 1)
		if bit, _ := bits.Get(i); bit != 1 {
			t.Fatalf("expect bit 1, got: %d\n", bit)
		}
	}
}

func TestCount(t *testing.T) {
	fmt.Println("Test: Count")

	bits := New(6401)

	for i := 0; i < bits.Len(); i++ {
		bits.Put(i, 1)
		if c := bits.Count(); c != i+1 {
			t.Fatalf("expect count of bit 1: %d, got: %d \n", i+1, c)
		}
	}
}

func randIdx(start int, max int, ifZero int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	end := r.Intn((max-start)/2) + start
	if end == start {
		end += ifZero
	}
	return end
}

func TestSet(t *testing.T) {
	fmt.Println("Test: Set")

	bits := New(6400)
	count := 0

	if c := bits.Count(); c != count {
		t.Fatalf("expect count of bit 1: %d, got: %d\n", count, c)
	}

	for i, n := 0, bits.Len(); i < n-2; {
		end := randIdx(i, n, 2)
		bits.Set(i, end-1)
		count += end - i
		if c := bits.Count(); c != count {
			t.Fatalf("expect count of bit 1: %d, got: %d\n", count, c)
		}
		i += end
	}
}

func TestClear(t *testing.T) {
	fmt.Println("Test: Clear")

	bits := New(6401)
	n := bits.Len()
	count := n
	bits.Set(0, count-1)

	if c := bits.Count(); c != count {
		t.Fatalf("expect count of bit 1: %d, got: %d\n", count, c)
	}

	for i := 0; i < n-2; {
		end := randIdx(i, n, 2)
		bits.Clear(i, end-1)
		count -= end - i
		if c := bits.Count(); c != count {
			t.Fatalf("expect count of bit 1: %d, got: %d\n", count, c)
		}
		i += end
	}

	bits.Clear(0, n-1)
	if c := bits.Count(); c != 0 {
		t.Fatalf("expect count of bit 1: 0, got: %d\n", c)
	}
}

func TestNot(t *testing.T) {
	fmt.Println("Test: Not")

	bits := New(6463)
	n := bits.Len()
	count := 0

	for i := 0; i < n-3; {
		start := randIdx(i, n, 1)
		end := randIdx(start, n, 2)
		bits.Set(start, end-1)
		count += start - i
		bits.Not(i, end-1)
		if c := bits.Count(); c != count {
			t.Fatalf("expect count of bit 1: %d, got: %d\n", count, c)
		}
		i += end
	}

	bits.Not(0, n-1)
	if c := bits.Count(); c != n-count {
		t.Fatalf("expect count of bit 1: %d, got: %d\n", count, c)
	}
}

func TestEq(t *testing.T) {
	fmt.Println("Test: Eq")

	n := 64
	bits1 := New(n)
	bits2 := New(n)

	for i := 0; i < n; i++ {
		bit := i % 2
		bits1.Put(i, bit)
		bits2.Put(i, bit)
		if !bits1.Eq(bits2) {
			t.Fatalf("expect equal, got %v != %v\n", bits1.ToArray(), bits2.ToArray())
		}
	}
}

func _testLtOrEq(t *testing.T, eq bool) {
	n := 6400
	bits1 := New(n)
	bits2 := New(n)
	idxs := []int{}

	for i := 0; i < n-2; i++ {
		idx := randIdx(i, n, 2)
		idxs = append(idxs, idx-1)
		bits1.Put(idx-1, 1)
		i += idx
	}

	for _, i := range idxs[:len(idxs)-1] {
		bits2.Put(i, 1)
		if !bits2.Leq(bits1) {
			t.Fatalf("expect equal, got %v != %v\n", bits1.ToArray(), bits2.ToArray())
		}
	}

	if eq {
		bits2.Put(idxs[len(idxs)-1], 1)
		if !bits2.Leq(bits1) {
			t.Fatalf("expect equal, got %v != %v\n", bits1.ToArray(), bits2.ToArray())
		}
	}
}

func TestLt(t *testing.T) {
	fmt.Println("Test: Lt")

	_testLtOrEq(t, false)
}

func TestLeq(t *testing.T) {
	fmt.Println("Test: Leq")

	_testLtOrEq(t, true)
}
