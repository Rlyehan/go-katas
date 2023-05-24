package main

import (
	"hash/fnv"
	"math"
)

type BloomFilter struct {
	bitset 	[]bool
	k 		uint
}

func NewBloomFilter(size uint, k uint) *BloomFilter {
	return &BloomFilter{make([]bool, size), k}
}

func (bf *BloomFilter) Add(item string) {
	indices := bf.hashValues(item)
	for _, index := range indices {
		bf.bitset[index] = true
	}
}

func (bf *BloomFilter) Test(item string) bool {
	indices := bf.hashValues(item)
	for _, index := range indices {
		if !bf.bitset[index] {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) hashValues(item string) []uint {
	indeces := make([]uint, bf.k)
	h := fnv.New64a()
	h.Write([]byte(item))
	hashValue := int(h.Sum64())

	for i := uint(0); i < bf.k; i++ {
		indeces[i] = uint(math.Abs(float64(hashValue + int(i)*hashValue))) % uint(len(bf.bitset))
	}

	return indeces
}

func main() {
	bf := NewBloomFilter(1000, 3)
	bf.Add("testitem")

	if bf.Test("testitem") {
		println("testitem is possibly in the set")
	} else {
		println("testitem is definitely not in the set")
	}

	if bf.Test("otheritem") {
		println("otheritem is possibly in the set")
	} else {
		println("otheritem is definitely not in the set")	
	}
}