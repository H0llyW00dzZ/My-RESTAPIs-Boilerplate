// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package stream

// calculateAESNonceCapacity calculates the nonce capacity for AES-CTR based on the overhead.
// It takes the overhead as input and returns the calculated nonce capacity.
// The nonce capacity is calculated as the sum of the AES block size and the overhead.
// If the calculated capacity is less than the AES block size, it is set to the AES block size.
// Otherwise, an additional capacity percentage is added to the calculated capacity.
func (s *Stream) calculateAESNonceCapacity(overhead int) int {
	capacity := s.aesBlock.BlockSize() + overhead
	if capacity < s.aesBlock.BlockSize() {
		capacity = s.aesBlock.BlockSize()
	} else {
		capacity += int(float64(capacity) * s.customizeNonce.AESNonceCapacity)
	}
	return capacity
}

// calculateChachaNonceCapacity calculates the nonce capacity for XChaCha20-Poly1305 based on the nonce size and overhead.
// It takes the nonce size and overhead as input and returns the calculated nonce capacity.
// The nonce capacity is calculated as the sum of the nonce size and the overhead.
// If the calculated capacity is less than the nonce size, it is set to the nonce size.
// Otherwise, an additional capacity percentage is added to the calculated capacity.
func (s *Stream) calculateChachaNonceCapacity(nonceSize, overhead int) int {
	capacity := nonceSize + overhead
	if capacity < nonceSize {
		capacity = nonceSize
	} else {
		capacity += int(float64(capacity) * s.customizeNonce.ChachaNonceCapacity)
	}
	return capacity
}
