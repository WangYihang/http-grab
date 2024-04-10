package model

import (
	"encoding/json"
	"math/big"
)

type BigIntWrapper struct {
	BigInt *big.Int
}

func (b *BigIntWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.BigInt.String())
}

func (b *BigIntWrapper) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	b.BigInt = new(big.Int)
	b.BigInt.SetString(str, 10)
	return nil
}
