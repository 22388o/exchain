package state

import (
	amino "github.com/tendermint/go-amino"

	cryptoamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
)

var cdc = amino.NewCodec()

func init() {
	cryptoamino.RegisterAmino(cdc)
}
