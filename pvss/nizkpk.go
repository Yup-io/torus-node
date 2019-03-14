// This implements the NIZKPK presented in Distributed Key Generation in the Wild
// with the aims of increasing asynchronicity within AVSS
package pvss

import (
	"math/big"

	"github.com/torusresearch/torus-public/common"
	"github.com/torusresearch/torus-public/secp256k1"
)

// Generates NIZK Proof with the aims of increasing asynchronicity within AVSS
func GenerateNIZKPK(s big.Int, r big.Int) (big.Int, big.Int, big.Int) {
	// create randomness
	v1 := *RandomBigInt()
	v2 := *RandomBigInt()
	t1 := common.BigIntToPoint(secp256k1.Curve.ScalarBaseMult(v1.Bytes()))
	t2 := common.BigIntToPoint(secp256k1.Curve.ScalarMult(&secp256k1.H.X, &secp256k1.H.Y, v2.Bytes()))
	gs := common.BigIntToPoint(secp256k1.Curve.ScalarBaseMult(s.Bytes()))
	gshr := secp256k1.Curve.Add(gs.X, gs.Y, secp256k1.Curve.ScalarMult(&secp256k1.H.X, &secp256k1.H.Y, r.Bytes()))

	//prepare bytestring for hashing
	bytesToBeHashed := append(secp256k1.G.X.Bytes(), secp256k1.H.X.Bytes()...)
	bytesToBeHashed = append(bytesToBeHashed, gs.X.Bytes()...)
	bytesToBeHashed = append(bytesToBeHashed, gshr.X.Bytes()...)
	bytesToBeHashed = append(bytesToBeHashed, t1.X.Bytes()...)
	bytesToBeHashed = append(bytesToBeHashed, t2.X.Bytes()...)

	c := new(big.Int)
	c.SetBytes(secp256k1.Keccak256(bytesToBeHashed))

	u1 := new(big.Int)
	u1.Set(&v1)
	cs := new(big.Int)
	cs.Mul(c, &s)
	cs.Mod(cs, secp256k1.GeneratorOrder)
	u1.Sub(u1, cs)
	u1.Mod(u1, secp256k1.GeneratorOrder)

	u2 := new(big.Int)
	u2.Set(&v2)
	cr := new(big.Int)
	cr.Mul(c, &r)
	cr.Mod(cr, secp256k1.GeneratorOrder)
	u2.Sub(u2, cr)
	u2.Mod(u2, secp256k1.GeneratorOrder)

	return *c, *u1, *u2

}

// func VerifyNIZKPK(c, u1, u2 big.Int, commitmentMatrixs [][]common.Point, commitmentMatrixsr [][]common.Point) {

// 	t1prime := common.BigIntToPoint(secp256k1.Curve.ScalarBaseMult(u1.Bytes()))
// }
