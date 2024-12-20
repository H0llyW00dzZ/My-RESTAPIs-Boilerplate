// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package convert provides utility functions for converting between different units of measurement.
//
// # Compatibility
//
// The current implementation depends on the units of measurement supported by the computer architecture.
// If the architecture is 32-bit, the [int] type will be 32-bit, which may limit the maximum supported unit.
//
// # Supported Units
//
// The package currently supports the following units of measurement:
//
//   - Bytes (B)
//   - Kilobytes (KB) and Kibibytes (KiB)
//   - Megabytes (MB) and Mebibytes (MiB)
//   - Gigabytes (GB) and Gibibytes (GiB)
//   - Terabytes (TB) and Tebibytes (TiB)
//
// # Limitations
//
//   - Exabytes (EB) and larger units are currently not supported as they might cause integer overflow on 32-bit or 64-bit architectures.
//   - The package assumes that the host computer architecture is either 32-bit or 64-bit.
//   - Support for quantum bits or 128-bit architectures is not currently implemented.
//
// # Future Enhancements
//
//   - Support for larger units such as Petabytes (PB), Exabytes (EB), Zettabytes (ZB), and Yottabytes (YB) could be added in the future (might required quantum bits or 128-bit architectures).
//   - Compatibility with quantum bits or 128-bit architectures could be considered as technology evolves.
package convert
