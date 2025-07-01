# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Product Description**: Added a `Description` field to the `Product` struct to allow for product descriptions.
- **Product-Level Discounts**: Products can now have individual percentage-based discounts, which are applied before cart-wide promotions.
- **In-Memory Cart Repository**: Implemented a complete, thread-safe in-memory cart repository.
- **Cart Repository Interface**: Defined a clear contract for cart repository operations, including `Create`, `GetByID`, `Update`, and `Delete`.
- **Custom Errors**: Introduced specific error types for more robust error handling (e.g., `ErrCartNotFound`, `ErrCartExists`).
- **Comprehensive Unit Tests**: Added extensive tests for the cart repository and product discount logic, including thread-safety checks.
- **Example Usage**: Updated `main.go` to demonstrate the new product discount functionality.

### Changed
- **BREAKING CHANGE**: The `Price` field in the `Product` struct has been changed from `int64` to `float64`. This requires updates to all code that interacts with product prices, including assignments, calculations, and potentially database schemas.
  - `GetDiscountedPrice()` method signature and internal calculations updated to reflect `float64` prices.
- The cart repository now operates in-memory, removing the need for a database connection.

### Fixed
- Corrected a typo in `cartRepository` that prevented compilation.

### Removed
- Removed the unused `database/sql` dependency.
