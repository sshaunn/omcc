package util

import (
	"fmt"
	"math/big"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange/bitget"
)

func SumMoneyDecimals(volumeList []bitget.CustomerVolume) (*big.Float, error) {
	sum := new(big.Float)

	for _, volume := range volumeList {
		// Create a new big.Rat for each float string
		f := new(big.Float)
		// Set the value of the big.Rat from the string float
		_, err := fmt.Sscan(volume.Volume, f)
		if err != nil {
			return nil, err
		}

		// Add the current big.Rat to the sum
		sum.Add(sum, f)
	}

	return sum, nil
}
