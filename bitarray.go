/*
BitArray for Golang.
*/
package bitarray

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type BitArray struct {
	lenpad int
	length int
	bytes  []byte
}

const (
	_BytesPW int = 8
	BitsPW   int = _BytesPW * 8
)

var msbmask = [8]byte{0xFF, 0xFE, 0xFC, 0xF8, 0xF0, 0xE0, 0xC0, 0x80}
var lsbmask = [8]byte{0x01, 0x03, 0x07, 0x0F, 0x1F, 0x3F, 0x7F, 0xFF}

// var count = [16]int{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4}

func nbytes(length int) int {
	return (((length + 7) & (^7)) / 8) // (length/8) + ((length%8)? 1:0)
}

func nwords(length int) int {
	return (((length + BitsPW - 1) & (^(BitsPW - 1))) / BitsPW)
}

func bytes2word(bs []byte) uint64 {
	var n uint64
	buf := bytes.NewBuffer(bs)
	err := binary.Read(buf, binary.BigEndian, &n)
	if err != nil {
		panic(err)
	}
	return n
}

// https://en.wikipedia.org/wiki/Hamming_weight
func countbits64(n uint64) int {
	n -= (n >> 1) & 0x5555555555555555
	n = (n & 0x3333333333333333) + ((n >> 2) & 0x3333333333333333)
	n = (n + (n >> 4)) & 0x0f0f0f0f0f0f0f0f
	return int((n * 0x0101010101010101) >> 56)
}

// New create a new BitArray with length(bits).
func New(length int) *BitArray {
	lenpad := nwords(length) * _BytesPW
	return &BitArray{
		lenpad: lenpad,
		length: length,
		bytes:  make([]byte, lenpad, lenpad),
	}
}

// Len return the length of the BitArray.
func (bits *BitArray) Len() int {
	return bits.length
}

// Count return the count of bit 1.
func (bits *BitArray) Count() int {
	length := 0

	for i := 0; i < bits.lenpad; i += _BytesPW {
		w := bytes2word(bits.bytes[i : i+_BytesPW])
		length += countbits64(w)
	}

	return length
}

func (bits *BitArray) indexOutOfRange(idx int) error {
	if idx < 0 || idx >= bits.length {
		msg := fmt.Sprintf("index %d out of range [%d, %d)", idx, 0, bits.length)
		return errors.New(msg)
	}
	return nil
}

// Get return the bit by index n.
// If index out of range [0, BitArray.Len()), return error.
func (bits *BitArray) Get(n int) (int, error) {
	if err := bits.indexOutOfRange(n); err != nil {
		return 0, err
	}
	return int((bits.bytes[n/8] >> byte(n%8)) & 1), nil
}

// Put set the nth bit with 0/1, and return the old value of nth bit.
// If index out of range [0, BitArray.Len()), return error.
func (bits *BitArray) Put(n int, bit int) (int, error) {
	if err := bits.indexOutOfRange(n); err != nil {
		return 0, err
	}
	prev, _ := bits.Get(n)

	if bit == 1 {
		bits.bytes[n/8] |= 1 << byte(n%8)
	} else {
		bits.bytes[n/8] &= ^(1 << byte(n%8))
	}

	return prev, nil
}

// Set the value of all bits to 1, which index range between low and high,
// both low and high included.
// low must less than high, and low/high cannot out of range [0, BitArray.Len()).
func (bits *BitArray) Set(low int, high int) error {
	if low > high {
		msg := fmt.Sprintf("low %d should less than high %d", low, high)
		return errors.New(msg)
	}
	for _, idx := range []int{low, high} {
		if err := bits.indexOutOfRange(idx); err != nil {
			return err
		}
	}

	lb, hb := low/8, high/8

	if lb < hb {
		bits.bytes[lb] |= msbmask[low%8]
		for i := lb + 1; i < hb; i++ {
			bits.bytes[i] = 0xFF
		}
		bits.bytes[hb] |= lsbmask[high%8]
	} else {
		bits.bytes[lb] |= (msbmask[low%8] & lsbmask[high%8])
	}

	return nil
}

// Clear set the value of all bits to 0, which index range between low and high,
// both low and high included.
// low must less than high, and low/high cannot out of range [0, BitArray.Len()).
func (bits *BitArray) Clear(low int, high int) error {
	if low > high {
		msg := fmt.Sprintf("low %d should less than high %d", low, high)
		return errors.New(msg)
	}
	for _, idx := range []int{low, high} {
		if err := bits.indexOutOfRange(idx); err != nil {
			return err
		}
	}

	lb, hb := low/8, high/8

	if lb < hb {
		bits.bytes[lb] &= ^msbmask[low%8]
		for i := lb + 1; i < hb; i++ {
			bits.bytes[i] = 0
		}
		bits.bytes[hb] &= ^lsbmask[high%8]
	} else {
		bits.bytes[lb] &= ^(msbmask[low%8] & lsbmask[high%8])
	}

	return nil
}

// Not flips the value of all bits, which index range between low and high,
// both low and high included.
// low must less than high, and low/high cannot out of range [0, BitArray.Len()).
func (bits *BitArray) Not(low int, high int) error {
	if low > high {
		msg := fmt.Sprintf("low %d should less than high %d", low, high)
		return errors.New(msg)
	}
	for _, idx := range []int{low, high} {
		if err := bits.indexOutOfRange(idx); err != nil {
			return err
		}
	}

	lb, hb := low/8, high/8

	if lb < hb {
		bits.bytes[lb] ^= msbmask[low%8]
		for i := lb + 1; i < hb; i++ {
			bits.bytes[i] ^= 0xFF
		}
		bits.bytes[hb] ^= lsbmask[high%8]
	} else {
		bits.bytes[lb] ^= (msbmask[low%8] & lsbmask[high%8])
	}

	return nil
}

// Eq check whether the BitArray is equal to another BitArray.
// If length isn't same, return false.
func (bits *BitArray) Eq(obits *BitArray) bool {
	if bits.length != obits.length {
		return false
	}

	for i := 0; i < bits.lenpad; i += _BytesPW {
		wself := bytes2word(bits.bytes[i : i+_BytesPW])
		wother := bytes2word(obits.bytes[i : i+_BytesPW])
		if wself != wother {
			return false
		}
	}
	return true
}

// Leq check whether the BitArray is the subset of the another.
// If length isn't same, return false.
func (bits *BitArray) Leq(obits *BitArray) bool {
	if bits.length != obits.length {
		return false
	}
	for i := 0; i < bits.lenpad; i += _BytesPW {
		wself := bytes2word(bits.bytes[i : i+_BytesPW])
		wother := bytes2word(obits.bytes[i : i+_BytesPW])
		if (wself & ^wother) != 0 {
			return false
		}
	}
	return true
}

// Lt check whether the BitArray is the proper subset of the another.
// If length isn't same, return false.
func (bits *BitArray) Lt(obits *BitArray) bool {
	if bits.length != obits.length {
		return false
	}
	lt := 0
	for i := 0; i < bits.lenpad; i += _BytesPW {
		wself := bytes2word(bits.bytes[i : i+_BytesPW])
		wother := bytes2word(obits.bytes[i : i+_BytesPW])
		if (wself & ^wother) != 0 {
			return false
		} else if wself != wother { // at least one word does not equal
			lt |= 1
		}
	}
	if lt == 0 {
		return false
	}
	return true
}

// Convert the BitArray to a array of integers, and return.
func (bits *BitArray) ToArray() []int {
	ints := make([]int, bits.length, bits.length)

	for i := 0; i < bits.length; i++ {
		ints[i], _ = bits.Get(i)
	}

	return ints
}
