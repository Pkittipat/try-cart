# Input Validation Workshop: Building Robust Go Applications

This document provides a comprehensive guide on implementing input validation in Go applications, based on our implementation of cart validation for the try-cart project.

## 📚 Table of Contents

- [Overview](#overview)
- [Problem Statement](#problem-statement)
- [Solution Architecture](#solution-architecture)
- [Implementation Details](#implementation-details)
- [Testing Strategy](#testing-strategy)
- [Best Practices](#best-practices)
- [Lessons Learned](#lessons-learned)

## Overview

This workshop demonstrates how to add robust input validation to a Go application using a shopping cart system as an example. We'll cover validation patterns, error handling, and comprehensive testing strategies.

### What You'll Learn

- How to design validation functions in Go
- Proper error handling and meaningful error messages
- Testing validation scenarios comprehensively
- Best practices for input validation in domain models

## Problem Statement

The original shopping cart implementation had a critical security and reliability issue:

```go
// BEFORE: No validation - accepts any input
func (c *Cart) AddProduct(product Product, quantity int64) {
    if item, ok := c.Items[product.ID]; ok {
        item.Quantity += quantity
        return
    }
    c.Items[product.ID] = &CartItem{Product: product, Quantity: quantity}
}
```

### Issues with the Original Code

1. **No quantity validation**: Could accept negative or zero quantities
2. **No product ID validation**: Could accept empty or nil product IDs
3. **No price validation**: Could accept negative prices
4. **No error handling**: No way to report validation failures
5. **No discount validation**: Could accept invalid discount percentages

### Real-World Impact

- Negative quantities could break inventory calculations
- Empty product IDs could cause data integrity issues
- Invalid prices could result in financial losses
- Poor user experience with unclear error states

## Solution Architecture

Our solution follows the **fail-fast principle** and implements validation at the domain layer:

```
┌─────────────────────────────────────┐
│           API Layer                 │
│  (HTTP handlers, input parsing)     │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│         Service Layer               │
│    (Business logic orchestration)  │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│         Domain Layer                │
│  *** VALIDATION HAPPENS HERE ***    │
│     (Core business rules)           │
└─────────────────────────────────────┘
```

### Why Validate at Domain Layer?

1. **Single responsibility**: Domain objects know their own validation rules
2. **Consistency**: Same validation logic regardless of entry point (API, CLI, etc.)
3. **Testability**: Easy to unit test without external dependencies
4. **Maintainability**: Changes to business rules are centralized

## Implementation Details

### Step 1: Update Method Signature

```go
// BEFORE: No error handling
func (c *Cart) AddProduct(product Product, quantity int64)

// AFTER: Returns error for validation failures
func (c *Cart) AddProduct(product Product, quantity int64) error
```

### Step 2: Create Validation Functions

We created separate validation functions following the **single responsibility principle**:

```go
// ValidateProduct validates product data
func ValidateProduct(product Product) error {
    if strings.TrimSpace(product.ID) == "" {
        return errors.New("product ID cannot be empty")
    }

    if product.Price < 0 {
        return errors.New("product price cannot be negative")
    }

    if !product.ValidateDiscount() {
        return errors.New("product discount must be between 0 and 100")
    }

    return nil
}

// ValidateQuantity validates quantity value
func ValidateQuantity(quantity int64) error {
    if quantity <= 0 {
        return errors.New("quantity must be a positive integer")
    }

    return nil
}
```

### Step 3: Implement Validation Logic

```go
func (c *Cart) AddProduct(product Product, quantity int64) error {
    // Validate product first
    if err := ValidateProduct(product); err != nil {
        return fmt.Errorf("invalid product: %w", err)
    }

    // Validate quantity
    if err := ValidateQuantity(quantity); err != nil {
        return fmt.Errorf("invalid quantity: %w", err)
    }

    // Only proceed if validation passes
    if item, ok := c.Items[product.ID]; ok {
        item.Quantity += quantity
        return nil
    }

    c.Items[product.ID] = &CartItem{Product: product, Quantity: quantity}
    return nil
}
```

### Key Design Decisions

1. **Error Wrapping**: Using `fmt.Errorf` with `%w` for proper error context
2. **Early Returns**: Fail fast on validation errors
3. **Meaningful Messages**: Clear, actionable error messages
4. **Separation of Concerns**: Separate functions for different validation types

## Testing Strategy

### Comprehensive Test Coverage

We implemented multiple layers of testing:

#### 1. Unit Tests for Individual Validators

```go
func TestValidateProduct(t *testing.T) {
    tests := []struct {
        name        string
        product     Product
        expectError bool
        errorMsg    string
    }{
        {
            name: "valid product",
            product: Product{
                ID:       "valid-id",
                Price:    100.00,
                Discount: 10,
            },
            expectError: false,
        },
        {
            name: "empty product ID",
            product: Product{
                ID:    "",
                Price: 100.00,
            },
            expectError: true,
            errorMsg:    "product ID cannot be empty",
        },
        // ... more test cases
    }
    // Test implementation...
}
```

#### 2. Integration Tests for AddProduct

```go
func TestCart_AddProduct_Validation(t *testing.T) {
    tests := []struct {
        name        string
        product     Product
        quantity    int64
        expectError bool
        errorMsg    string
    }{
        {
            name: "valid product and quantity",
            product: Product{
                ID:    "1",
                Price: 10.00,
            },
            quantity:    2,
            expectError: false,
        },
        {
            name: "negative quantity",
            product: Product{
                ID:    "1",
                Price: 10.00,
            },
            quantity:    -5,
            expectError: true,
            errorMsg:    "invalid quantity: quantity must be a positive integer",
        },
        // ... more test cases
    }
    // Test implementation...
}
```

#### 3. Backward Compatibility Tests

We updated all existing tests to handle the new error return:

```go
// BEFORE: No error handling
cart.AddProduct(product, quantity)

// AFTER: Proper error handling
err := cart.AddProduct(product, quantity)
assert.NoError(t, err)
```

### Test Categories Covered

- ✅ **Happy Path**: Valid inputs should succeed
- ✅ **Edge Cases**: Boundary values (0, negative numbers)
- ✅ **Invalid Inputs**: Empty strings, nil values
- ✅ **Error Messages**: Verify specific error text
- ✅ **Backward Compatibility**: Existing functionality preserved

## Best Practices

### 1. Validation Design Patterns

#### ✅ DO: Validate Early and Often
```go
func (c *Cart) AddProduct(product Product, quantity int64) error {
    // Validate immediately at method entry
    if err := ValidateProduct(product); err != nil {
        return fmt.Errorf("invalid product: %w", err)
    }
    // ... rest of logic
}
```

#### ❌ DON'T: Silent Failures
```go
func (c *Cart) AddProduct(product Product, quantity int64) {
    if quantity <= 0 {
        return // Silent failure - caller doesn't know why it failed
    }
}
```

### 2. Error Message Guidelines

#### ✅ DO: Provide Actionable Messages
```go
return errors.New("quantity must be a positive integer")
```

#### ❌ DON'T: Vague Error Messages
```go
return errors.New("invalid input")
```

### 3. Function Design

#### ✅ DO: Single Responsibility
```go
// Each function validates one concern
func ValidateProduct(product Product) error { /* ... */ }
func ValidateQuantity(quantity int64) error { /* ... */ }
```

#### ❌ DON'T: Mixed Responsibilities
```go
func ValidateEverything(product Product, quantity int64, user User) error {
    // Too many responsibilities in one function
}
```

### 4. Testing Patterns

#### ✅ DO: Table-Driven Tests
```go
tests := []struct {
    name        string
    input       Input
    expectError bool
    errorMsg    string
}{
    // Multiple test cases in structured format
}
```

#### ✅ DO: Test Error Messages
```go
if tt.expectError {
    assert.Error(t, err)
    assert.Equal(t, tt.errorMsg, err.Error())
}
```

## Lessons Learned

### 1. Start with Domain Validation

**Why**: Domain models know their own business rules best
**How**: Implement validation in the domain layer, not at the API layer

### 2. Error Context is Critical

**Before**:
```go
return errors.New("invalid")
```

**After**:
```go
return fmt.Errorf("invalid product: %w", err)
```

**Impact**: Debugging becomes much easier with proper error context

### 3. Comprehensive Testing Pays Off

- We wrote 30+ test cases covering all validation scenarios
- Found edge cases during test writing that we hadn't considered
- Tests serve as documentation for expected behavior

### 4. Backward Compatibility Matters

- Updated all existing tests to handle new error returns
- Ensured existing functionality still works
- Made changes incrementally to avoid breaking existing code

### 5. Tool Integration

We integrated multiple tools for quality assurance:

```bash
# Run tests
go test ./internal/domain/cart/... -v

# Check for issues
go vet ./...

# Build verification
go build ./...
```

## Implementation Checklist

When implementing validation in your Go applications:

### Planning Phase
- [ ] Identify all input validation requirements
- [ ] Define error message standards
- [ ] Plan backward compatibility strategy

### Implementation Phase
- [ ] Create separate validation functions
- [ ] Update method signatures to return errors
- [ ] Implement meaningful error messages
- [ ] Add proper error wrapping

### Testing Phase
- [ ] Write unit tests for each validator
- [ ] Write integration tests for main functions
- [ ] Update existing tests for new error handling
- [ ] Test edge cases and boundary conditions

### Quality Assurance
- [ ] Run `go vet` for static analysis
- [ ] Ensure all tests pass
- [ ] Verify code builds successfully
- [ ] Check test coverage

## Conclusion

This workshop demonstrated how to implement robust input validation in Go applications. Key takeaways:

1. **Validate early and fail fast** - Catch errors at the domain boundary
2. **Provide meaningful error messages** - Help users understand what went wrong
3. **Test comprehensively** - Cover happy paths, edge cases, and error conditions
4. **Design for maintainability** - Separate concerns and use clear function names

The result is a more robust, reliable, and maintainable codebase that provides clear feedback when things go wrong.

### Further Reading

- [Go Error Handling Best Practices](https://blog.golang.org/error-handling-and-go)
- [Table-Driven Tests in Go](https://github.com/golang/go/wiki/TableDrivenTests)
- [Domain-Driven Design in Go](https://medium.com/@hatajoe/clean-architecture-in-go-4030f11ec1b1)

---

*This documentation was created as part of the cart validation workshop for the try-cart project.*