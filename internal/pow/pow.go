package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
)

var ErrCalcProofHash = errors.New("cannot to calculate proof hash")

const nonceSize = 8

func CalcProof(difficulty byte, data []byte) (proofNonce, proofHash []byte, err error) {
	nonceOffset := len(data)
	buf := make([]byte, nonceOffset+nonceSize)
	copy(buf, data)
	// buf = | data | nonce
	var hash [32]byte

	var nonce uint64
	for nonce < math.MaxUint64 {
		binary.BigEndian.PutUint64(buf[nonceOffset:], nonce)

		hash = sha256.Sum256(buf)

		if leadingZerosCount(hash[:]) >= difficulty {
			proofNonce = buf[nonceOffset:]
			proofHash = hash[:]
			return
		} else {
			nonce++
		}
	}

	err = ErrCalcProofHash
	return
}

func CheckProof(difficulty byte, data []byte, proofNonce []byte) bool {
	nonceOffset := len(data)
	buf := make([]byte, nonceOffset+nonceSize)
	copy(buf, data)
	copy(buf[nonceOffset:], proofNonce)
	// buf = | data | nonce

	return CheckBufProof(difficulty, buf)
}

func CheckBufProof(difficulty byte, buf []byte) bool {
	// buf = | data | nonce
	hash := sha256.Sum256(buf)

	return leadingZerosCount(hash[:]) >= difficulty
}
