package lib

import (
	"fmt"
	"math/rand"
	"strconv"
)

func RandomPortion(soldP []uint8, units string) ([]uint8, error) {
	sold := make(map[uint8]struct{})
	for _, p := range soldP {
		sold[p] = struct{}{}
	}
	all := make(map[uint8]struct{})
	for i := 0; i < 100; i++ {
		all[uint8(i)] = struct{}{}
	}
	var lotP []uint8
	for k := range all {
		_, ok := sold[k]
		if !ok {
			lotP = append(lotP, k)
		}
	}

	u, err := strconv.Atoi(units)
	if err != nil {
		return nil, fmt.Errorf("RandomPortion: %v", err)
	}
	var randP []uint8
	for i := 0; i < u; i++ {
		r := rand.Intn(len(lotP))
		randP = append(randP, lotP[r])
	}
	return randP, nil
}
