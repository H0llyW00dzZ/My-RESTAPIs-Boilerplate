// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package keyrotation provides a key manager for rotating encryption keys periodically.
//
// IMPORTANT: This package is still a work in progress (WIP) and should be used with caution.
// The reason for including this package here is to avoid creating another repository for experimental purposes related to cryptography.
// Implementation and API may change in future versions without backward compatibility.
// There are no guarantees, so use it only for educational, testing, and experimental purposes related to cryptography (e.g., ciphers, secure systems).
//
// Also note that Go is suitable for cryptographic systems due to its memory safety, garbage collection features, and powerful standard library.
// It allows for both high-level and low-level programming, unlike languages like C or C++.
// However, the ecosystem for Go cryptography is still relatively small compared to other (e.g., cloud, web).
//
// So do not use this package in production environments (e.g., for encrypting sensitive data). Instead, use it for educational purposes, which is a better fit.
//
// The keyrotation package allows you to manage the rotation of encryption keys used for
// AES and ChaCha20-Poly1305 encryption. It provides a KeyManager struct that handles the
// periodic rotation of keys based on a specified rotation interval.
//
// Key Features:
//   - Periodic key rotation: The KeyManager automatically rotates the encryption keys at the
//     specified interval, ensuring that the keys are regularly updated for enhanced security.
//   - Thread-safe access: The package uses synchronization primitives (mutex) to ensure
//     thread-safe access to the current encryption keys.
//   - Customizable rotation interval: You can specify the desired key rotation interval when
//     creating a new instance of KeyManager.
//   - Error handling: The package provides error types (ErrInvalidKeySize and ErrKeyRotationFailed)
//     to handle and communicate errors related to key management and rotation.
//
// Usage:
//  1. Create a new instance of KeyManager using the NewKeyManager function, providing the
//     initial AES and ChaCha20-Poly1305 keys and the desired rotation interval.
//  2. Use the GetCurrentKeys method to retrieve the current encryption keys whenever needed.
//  3. Call the StartKeyRotation method to start the background key rotation process.
//
// Example:
//
//	initialAESKey := make([]byte, keyrotation.KeySize)
//	initialChaChaKey := make([]byte, keyrotation.KeySize)
//	rotationInterval := 30 * time.Minute
//
//	km, err := keyrotation.NewKeyManager(initialAESKey, initialChaChaKey, rotationInterval)
//	if err != nil {
//		// Handle error
//	}
//
//	go km.StartKeyRotation()
//
//	// Retrieve current keys
//	currentAESKey, currentChaChaKey := km.GetCurrentKeys()
//
// Note:
//   - The package assumes that the initial AES and ChaCha20-Poly1305 keys are securely generated
//     and provided by the user. It is the responsibility of the user to ensure the security and
//     randomness of the initial keys.
//   - The key rotation interval should be chosen carefully based on the specific security
//     requirements of your application. A shorter interval provides better security but may
//     impact performance, while a longer interval reduces the rotation overhead but may
//     increase the window of vulnerability if a key is compromised.
//   - The package uses the crypto/rand package for generating random keys during rotation.
//     Ensure that your system has a reliable source of randomness for secure key generation.
package keyrotation
