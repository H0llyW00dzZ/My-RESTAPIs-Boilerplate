// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package gpg provides functionality for encrypting data using OpenPGP/GPG public keys.
//
// This package includes utilities to create and manage key rings, encrypt files,
// and handle streaming encryption for efficient data transmission or other.
//
// Unlike GPG Proton built on top (fork) the standard library, this package
// is designed with a modern approach, focusing on top of I/O operations.
package gpg
