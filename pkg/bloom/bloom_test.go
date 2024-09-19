package bloom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	STD_FILTER_SIZE        uint64 = 1000
	STD_NUM_HASH_FUNCTIONS uint8  = 3
	STD_ENTROPY            uint8  = 8

	TEST_KEY       = []byte("foo")
	TEST_FALSE_KEY = []byte("sam")
)

func TestNew(t *testing.T) {
	b, err := New(STD_FILTER_SIZE, STD_NUM_HASH_FUNCTIONS, STD_ENTROPY)
	assert.Nil(t, err)
	assert.NotNil(t, b)
}

func TestNewWithError(t *testing.T) {
	b, err := New(STD_FILTER_SIZE, STD_NUM_HASH_FUNCTIONS, STD_ENTROPY+99)
	assert.Nil(t, b)
	assert.NotNil(t, err)
}

func TestAddAndCheckKey(t *testing.T) {
	b, err := New(STD_FILTER_SIZE, STD_NUM_HASH_FUNCTIONS, STD_ENTROPY)
	assert.Nil(t, err)
	assert.NotNil(t, b)

	err = b.AddKey(TEST_KEY)
	assert.Nil(t, err)

	present, err := b.CheckKey(TEST_KEY)
	assert.Nil(t, err)
	assert.True(t, present)

	present, err = b.CheckKey(TEST_FALSE_KEY)
	assert.Nil(t, err)
	assert.False(t, present)
}

func TestRemoveKey(t *testing.T) {
	b, err := New(STD_FILTER_SIZE, STD_NUM_HASH_FUNCTIONS, STD_ENTROPY)
	assert.Nil(t, err)
	assert.NotNil(t, b)

	err = b.AddKey(TEST_KEY)
	assert.Nil(t, err)

	present, err := b.CheckKey(TEST_KEY)
	assert.Nil(t, err)
	assert.True(t, present)

	err = b.RemoveKey(TEST_KEY)
	assert.Nil(t, err)

	present, err = b.CheckKey(TEST_KEY)
	assert.Nil(t, err)
	assert.False(t, present)

	err = b.RemoveKey(TEST_FALSE_KEY)
	assert.Nil(t, err)
}
