package pow

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcAndCheckPoW(t *testing.T) {
	var difficulty byte = 10

	data := make([]byte, 64)
	rand.Read(data)

	// calc proof
	nonce, _, err := CalcProof(difficulty, data)
	assert.NoError(t, err)

	// check valid proof
	isValid := CheckProof(difficulty, data, nonce)
	assert.True(t, isValid)

	// if invalid data
	data[0], data[1] = data[1], data[0]
	isValid = CheckProof(difficulty, data, nonce)
	assert.False(t, isValid)
	data[0], data[1] = data[1], data[0]

	// is invalid proof nonce
	rand.Read(nonce)
	isValid = CheckProof(difficulty, data, nonce)
	assert.False(t, isValid)
}

func BenchmarkCalculateProof(b *testing.B) {
	data := make([]byte, 64)
	rand.Read(data)

	var difficulty byte = 25

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalcProof(difficulty, data)
	}
}

func BenchmarkCheckProof(b *testing.B) {
	data := make([]byte, 64)
	rand.Read(data)

	var difficulty byte = 10

	nonce, _, _ := CalcProof(difficulty, data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckProof(difficulty, data, nonce)
	}
}

func BenchmarkCheckBufProof(b *testing.B) {
	const tokenSize = 64
	data := make([]byte, tokenSize+nonceSize)
	rand.Read(data[:tokenSize])

	var difficulty byte = 10

	nonce, _, _ := CalcProof(difficulty, data)
	copy(data[tokenSize:], nonce)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckBufProof(difficulty, data)
	}
}
