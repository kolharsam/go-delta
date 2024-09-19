package hash

import (
	"fmt"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testcase struct {
	input           string
	expected_output []uint64
}

const STD_ENTROPY uint8 = 8

func TestNew(t *testing.T) {
	h, err := New(3, 100, STD_ENTROPY)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(h.functions))
}

func TestNewWithError(t *testing.T) {
	h, err := New(2, 100, STD_ENTROPY)
	assert.NotNil(t, err)
	assert.Nil(t, h)

	h, err = New(6, 100, STD_ENTROPY)
	assert.NotNil(t, err)
	assert.Nil(t, h)
}

func TestNewWithEntropyError(t *testing.T) {
	h, err := New(3, 100, 0)
	assert.NotNil(t, err)
	assert.Nil(t, h)

	h, err = New(4, 100, STD_ENTROPY+90)
	assert.NotNil(t, err)
	assert.Nil(t, h)
}

func TestNewWithDiffNumbers(t *testing.T) {
	h, err := New(4, 100, STD_ENTROPY)
	assert.NotNil(t, h)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(h.functions))

	h, err = New(5, 100, STD_ENTROPY)
	assert.NotNil(t, h)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(h.functions))
}

func TestNewForFunctionSize(t *testing.T) {
	h, err := New(3, 32, STD_ENTROPY)
	assert.NotNil(t, h)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(h.functions))
	assert.Equal(t, uint64(32), h.filterSize)
}

var testCases = []testcase{
	{
		input:           "foo",
		expected_output: []uint64{27, 15, 16},
	},
	{
		input:           "sam",
		expected_output: []uint64{31, 28, 28},
	},
	{
		input:           "bubblegum",
		expected_output: []uint64{25, 14, 9},
	},
	{
		input:           "foo fighters",
		expected_output: []uint64{14, 9, 18},
	},
	{
		input:           "hash functions are great",
		expected_output: []uint64{9, 4, 28},
	},
	{
		input:           "",
		expected_output: []uint64{13, 20, 29},
	},
}

func TestGetPostionsInFilter(t *testing.T) {
	filterSize := uint64(32)

	h, err := New(3, filterSize, STD_ENTROPY)
	assert.Nil(t, err)
	assert.NotNil(t, h)

	for _, testcase := range testCases {
		actual, err := h.GetPostionsInFilter([]byte(testcase.input))
		assert.Nil(t, err)
		assert.ElementsMatch(t, actual, testcase.expected_output)
	}
}

type MockErrorHash struct{}

func (m MockErrorHash) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("mock error")
}

func (m MockErrorHash) Sum(b []byte) []byte { return nil }
func (m MockErrorHash) Reset()              {}
func (m MockErrorHash) Size() int           { return 0 }
func (m MockErrorHash) BlockSize() int      { return 0 }

func TestGetPostionsInFilterError(t *testing.T) {
	h := &Hash{
		functions:  []hash.Hash{MockErrorHash{}},
		filterSize: 32,
		entropy:    8,
	}

	_, err := h.GetPostionsInFilter([]byte("test"))
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}
