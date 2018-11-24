package dkgnode

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

//Validates transactions to be delivered to the BFT. is the master switch for all tx
//TODO: create variables for types here and in bftrpc.go
func (app *ABCIApp) ValidateBFTTx(tx []byte) (bool, *[]common.KVPair, error) {
	if bytes.Compare(tx[:len([]byte("mug00"))], []byte("mug00")) != 0 {
		return false, nil, errors.New("Tx signature is not mug")
	}
	var tags []common.KVPair
	txNoSig := tx[len([]byte("mug00")):]

	switch txNoSig[0] {
	case byte(1): //PubpolyTx
		pubPolyTx := DefaultBFTTxWrapper{&PubPolyBFTTx{}}
		err := pubPolyTx.DecodeBFTTx(txNoSig)
		if err != nil {
			return false, nil, err
		}
		//verify correct epoch
		return true, nil, nil
		//verify share index has not yet been submitted for epoch

	case byte(2): // EpochTx
		EpochTx := DefaultBFTTxWrapper{&EpochBFFTx{}}
		err := EpochTx.DecodeBFTTx(txNoSig)
		if err != nil {
			return false, nil, err
		}
		fmt.Println("ATTATCHING TAGS")
		tags = []common.KVPair{
			{Key: []byte("epoch"), Value: []byte{byte(1)}},
		}
		return true, &tags, nil

	}
	return false, &tags, errors.New("Tx type not recognised")
}