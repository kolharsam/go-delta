package bitset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkGetValue(t *testing.T, b *Bitset, pos uint64, value bool) {
	g, err := b.Get(pos)
	assert.Nil(t, err)
	assert.Equal(t, value, g)
}

func TestNew(t *testing.T) {
	b := New(100)
	assert.GreaterOrEqual(t, b.Len()*64, uint64(100))
}

func TestSetAndGet(t *testing.T) {
	b := New(100)
	b.Set(50)

	checkGetValue(t, b, 50, true)
	checkGetValue(t, b, 51, false)

	_, err := b.Get(9000)
	assert.NotNil(t, err)
}

func TestSetN(t *testing.T) {
	b := New(100)
	b.SetN(10, 20, 30, 45)

	checkGetValue(t, b, 10, true)
	checkGetValue(t, b, 20, true)
	checkGetValue(t, b, 30, true)
	checkGetValue(t, b, 45, true)
	checkGetValue(t, b, 90, false)
}

func TestSetNWithError(t *testing.T) {
	b := New(100)
	err := b.SetN(10, 20, 30, 45, 122)

	assert.NotNil(t, err)
	checkGetValue(t, b, 10, true)
	checkGetValue(t, b, 20, true)
	checkGetValue(t, b, 30, true)
	checkGetValue(t, b, 45, true)

	_, err = b.Get(122)
	assert.NotNil(t, err)
}

func TestGetN(t *testing.T) {
	b := New(100)
	err := b.SetN(10, 20, 30, 45)

	assert.Nil(t, err)
	checkGetValue(t, b, 10, true)
	checkGetValue(t, b, 20, true)
	checkGetValue(t, b, 30, true)
	checkGetValue(t, b, 45, true)

	present, err := b.GetN(10, 20, 30)
	assert.Nil(t, err)
	assert.True(t, present)

	present, err = b.GetN(20, 30, 34)
	assert.Nil(t, err)
	assert.True(t, present)

	_, err = b.GetN(20, 30, 900)
	assert.NotNil(t, err)
}

func TestRemove(t *testing.T) {
	b := New(100)
	b.Set(50)

	checkGetValue(t, b, 50, true)
	checkGetValue(t, b, 51, false)

	b.Remove(50)

	checkGetValue(t, b, 50, false)
	checkGetValue(t, b, 51, false)

	err := b.Remove(90000)
	assert.NotNil(t, err)
}

func TestRemoveN(t *testing.T) {
	b := New(100)
	b.Set(50)
	b.Set(55)

	checkGetValue(t, b, 50, true)
	checkGetValue(t, b, 51, false)
	checkGetValue(t, b, 55, true)

	b.RemoveN(50, 55)

	checkGetValue(t, b, 50, false)
	checkGetValue(t, b, 51, false)
	checkGetValue(t, b, 55, false)

	err := b.RemoveN(908008, 2323)
	assert.NotNil(t, err)
}

func TestCount(t *testing.T) {
	b := New(100)

	b.Set(10)
	checkGetValue(t, b, 10, true)

	b.Set(20)
	checkGetValue(t, b, 20, true)

	b.Set(30)
	checkGetValue(t, b, 30, true)

	assert.Equal(t, uint64(3), b.Count())
}

func TestReset(t *testing.T) {
	b := New(100)

	b.Set(10)
	checkGetValue(t, b, 10, true)

	b.Set(20)
	checkGetValue(t, b, 20, true)

	assert.Equal(t, uint64(2), b.Count())

	b.Reset()

	assert.Equal(t, uint64(0), b.Count())
}

func TestCopy(t *testing.T) {
	b1 := New(100)

	b1.Set(10)
	checkGetValue(t, b1, 10, true)

	b1.Set(20)
	checkGetValue(t, b1, 20, true)

	b2 := b1.Copy()
	assert.Equal(t, b1.Count(), b2.Count())

	b2.Set(30)
	checkGetValue(t, b2, 30, true)
	checkGetValue(t, b1, 30, false)
}

func TestString(t *testing.T) {
	b := New(10)
	b.Set(1)
	b.Set(3)
	b.Set(5)

	expected := "0101010000"

	if b.String() != expected {
		t.Errorf("String() returned %s, expected %s", b.String(), expected)
	}
}

func TestLargeBitset(t *testing.T) {
	b := New(1000000)
	err := b.Set(500000)

	assert.Nil(t, err)
	assert.Equal(t, uint64(1), b.Count())
}

func TestPeripheryOfBitArrays(t *testing.T) {
	b := New(64)
	b.Set(63)
	checkGetValue(t, b, 63, true)

	b = New(65)
	b.Set(64)
	checkGetValue(t, b, 64, true)
}

func TestErrorCase(t *testing.T) {
	b := New(100)
	err := b.Set(1000000)

	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), b.Count())
}
