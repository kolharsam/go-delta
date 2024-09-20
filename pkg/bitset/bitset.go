package bitset

import (
	"fmt"
	"strings"
	"sync"
)

var setError = "failed to set bit at position %d"
var posError = "pos[%d] provided is invalid against the size of the bitset [%d]"

// Bitset is useful for storing bits via the integer type (uint64)
// This forms the base upon which the bloom filter is built on
// Bitset is also ready for async operations
type Bitset struct {
	mtx  sync.RWMutex
	bits []uint64
	size uint64
}

// getIndexPos provides the position of the bit based on
// the bits array in Bitset. the computation is based on
// the underlying data type - uint 64
func getIndexPos(pos uint64) (index, bitPos uint64) {
	index, bitPos = pos/64, pos%64
	return
}

// New creates a new Bitset with the given number of bits
func New(size uint64) *Bitset {
	bits := make([]uint64, (size+63)/64)

	return &Bitset{bits: bits, size: size}
}

// Set sets the bit at the given position to 1
func (b *Bitset) Set(pos uint64) error {
	if pos > b.size {
		return fmt.Errorf(setError, pos)
	}

	index, bitPos := getIndexPos(pos)
	b.mtx.Lock()
	b.bits[index] |= 1 << bitPos
	b.mtx.Unlock()

	return nil
}

// SetN helps set `n` number of positions to `1` at once
func (b *Bitset) SetN(npos ...uint64) error {
	for i := 0; i < len(npos); i++ {
		err := b.Set(npos[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// Remove sets the bit at the given position to 0
func (b *Bitset) Remove(pos uint64) error {
	if pos > b.size {
		return fmt.Errorf(posError, pos, b.size)
	}

	index, bitPos := getIndexPos(pos)
	b.mtx.Lock()
	b.bits[index] &= ^(1 << bitPos)
	b.mtx.Unlock()

	return nil
}

// RemoveN unsets `n` number of positions in the bitset
func (b *Bitset) RemoveN(npos ...uint64) error {
	for i := 0; i < len(npos); i++ {
		err := b.Remove(npos[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// Get returns true if the bit at the given position is 1
func (b *Bitset) Get(pos uint64) (bool, error) {
	if pos > b.size {
		return false, fmt.Errorf(posError, pos, b.size)
	}

	index, bitPos := getIndexPos(pos)
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	return (b.bits[index] & (1 << bitPos)) != 0, nil
}

// GetN helps check whether `n` positions
// of interest are set to 1 in the bitset
func (b *Bitset) GetN(npos ...uint64) (bool, error) {
	var finalRes bool

	for i := 0; i < len(npos); i++ {
		res, err := b.Get(npos[i])
		if err != nil {
			return false, err
		}
		finalRes = finalRes || res
	}

	return finalRes, nil
}

// Count returns the number of bits set to 1 in the entire Bitset
func (b *Bitset) Count() uint64 {
	var count uint64
	b.mtx.RLock()
	for _, v := range b.bits {
		count += count64(v)
	}
	b.mtx.RUnlock()
	return count
}

// count64 is a helper function to count set bits in a uint64
// NOTE: http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
func count64(num uint64) uint64 {
	num -= (num >> 1) & 0x5555555555555555
	num = (num>>2)&0x3333333333333333 + num&0x3333333333333333
	num += num >> 4
	num &= 0x0f0f0f0f0f0f0f0f
	num *= 0x0101010101010101

	return uint64(num >> 56)
}

// Len returns the size of the bitset
func (b *Bitset) Len() uint64 {
	return b.size
}

// Reset sets all of the bits in the bitset to 0
func (b *Bitset) Reset() {
	b.mtx.Lock()
	for i := range b.bits {
		b.bits[i] = 0
	}
	b.mtx.Unlock()
}

// Copy returns a new copy of the current state of the bitset
func (b *Bitset) Copy() *Bitset {
	newBitset := make([]uint64, len(b.bits))
	copy(newBitset, b.bits)
	return &Bitset{
		bits: newBitset,
		size: b.size,
	}
}

// String returns the current bitset in the form of a string
// This is not meant for compute but for debugging or presentations
func (b *Bitset) String() string {
	var bitset strings.Builder
	var i uint64

	for i = 0; i < b.size; i++ {
		currentState, _ := b.Get(uint64(i))
		if currentState {
			bitset.WriteString("1")
		} else {
			bitset.WriteString("0")
		}
	}

	return bitset.String()
}

// TODO: implement union and intersection methods if necessary
