// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

// Package bcrypt provides a secure way to hash and compare passwords using the bcrypt algorithm.
// It is suitable for basic authentication and user management when interacting with a database.
//
// The package offers a Service interface with methods for hashing passwords or other and comparing plaintext
// passwords with stored hashes. It also provides a New function to create a new instance of the
// bcrypt password hashing service.
//
// Example Usage:
//
//  1. Create a new instance of the bcrypt service using the New function:
//
//     bcryptService := bcrypt.New()
//
//  2. Hash a password using the HashPassword method:
//
//     hashedPassword, err := bcryptService.HashPassword(password)
//     if err != nil {
//     // Handle the error
//     }
//
//  3. Store the hashed password in the database along with the user's credentials.
//
//  4. When a user attempts to log in, retrieve the user's record from the database based on the
//     provided username or email.
//
//  5. Compare the user's entered password with the stored hashed password using the ComparePassword method:
//
//     if bcryptService.ComparePassword(enteredPassword, storedHash) {
//     // Password is correct, authenticate the user
//     } else {
//     // Password is incorrect, reject the login attempt
//     }
//
// Note: It is important to store only the hashed passwords in the database and never store plaintext passwords.
// The bcrypt package ensures that the password hashing is secure and resilient against various attacks.
//
// The package uses a default cost factor of 10 for the bcrypt algorithm, which provides a good balance
// between security and performance. However, you can adjust the cost factor by modifying the
// bcrypt.DefaultCost constant in the HashPassword function if needed.
//
// It is recommended to use this package in combination with other security measures, such as HTTPS/TLS
// for secure communication, secure session management (which is already implemented here along with better encryption/decryption that require an 99999999999 cpu for attack), and protection against common vulnerabilities like
// SQL injection and cross-site scripting (XSS), which is already implemented here for XSS protection.
//
// For more information about the bcrypt algorithm and its security properties, refer to the
// official Go documentation: https://pkg.go.dev/golang.org/x/crypto/bcrypt
// REF: https://gowebexamples.com/password-hashing/
package bcrypt
