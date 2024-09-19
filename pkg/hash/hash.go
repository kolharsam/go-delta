package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"math/big"
)

type Hash struct {
	functions  []hash.Hash
	filterSize uint64
	entropy    uint8
}

func New(numFunctions uint8, filterSize uint64, entropy uint8) (*Hash, error) {
	if numFunctions < 3 || numFunctions > 5 {
		return nil, fmt.Errorf("failed to init hashes since [3-5] hash functions have to be configured")
	}

	if entropy <= 0 || entropy >= 20 {
		return nil, fmt.Errorf("failed to init hashes since entropy bytes has to be between [0-20] (non-inclusive)")
	}

	var functions []hash.Hash

	s1 := sha1.New()        // output => 20 bytes
	s256 := sha256.New()    // output => 32 bytes
	s512 := sha512.New()    // output => 64 bytes
	s384 := sha512.New384() // output => 48 bytes
	s224 := sha256.New224() // output => 28 bytes

	if numFunctions >= 3 {
		functions = append(functions, s1)
		functions = append(functions, s256)
		functions = append(functions, s512)
	}

	if numFunctions >= 4 {
		functions = append(functions, s224)
	}

	if numFunctions == 5 {
		functions = append(functions, s384)
	}

	return &Hash{functions, filterSize, entropy}, nil
}

func (h *Hash) GetPostionsInFilter(key []byte) ([]uint64, error) {
	results := make([]uint64, 0)

	for _, hashFunction := range h.functions {
		hashFunction.Reset()
		_, err := hashFunction.Write(key)
		if err != nil {
			return nil, err
		}
		v := hashFunction.Sum(nil)
		results = append(results, hashToPosition(v, h.filterSize, h.entropy))
	}

	return results, nil
}

func hashToPosition(hash []byte, filterSize uint64, entropyBytes uint8) uint64 {
	hashInt := new(big.Int).SetBytes(hash[:entropyBytes])

	return new(big.Int).Mod(hashInt, big.NewInt(int64(filterSize))).Uint64()
}
