package bloom

import (
	"github.com/kolharsam/go-delta/pkg/bitset"
	"github.com/kolharsam/go-delta/pkg/hash"
)

type Bloom struct {
	bitset *bitset.Bitset
	hash   *hash.Hash
}

func New(filterSize uint64, numHashFunctions uint8, entropy uint8) (*Bloom, error) {
	bitset := bitset.New(filterSize)
	hashFunctions, err := hash.New(numHashFunctions, filterSize, entropy)

	if err != nil {
		return nil, err
	}

	return &Bloom{
		bitset: bitset,
		hash:   hashFunctions,
	}, nil
}

func (b *Bloom) AddKey(key []byte) error {
	positions, err := b.hash.GetPostionsInFilter(key)
	if err != nil {
		return err
	}

	err = b.bitset.SetN(positions...)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bloom) CheckKey(key []byte) (bool, error) {
	positions, err := b.hash.GetPostionsInFilter(key)
	if err != nil {
		return false, err
	}

	present, err := b.bitset.GetN(positions...)
	if err != nil {
		return false, err
	}

	return present, nil
}

func (b *Bloom) RemoveKey(key []byte) error {
	positions, err := b.hash.GetPostionsInFilter(key)
	if err != nil {
		return err
	}

	err = b.bitset.RemoveN(positions...)
	if err != nil {
		return err
	}

	return nil
}
