# Documentation

This directory contains learning materials and guides for the try-cart project.

## Available Documents

### 🤖 [Claude Code + GitHub Integration](./CLAUDE_CODE_GITHUB_WORKFLOW.md)
A comprehensive guide on setting up Claude Code with GitHub integration using MCP and custom prompts. This workshop covers:

- **MCP GitHub Setup**: Configuring GitHub API access
- **Custom Prompts**: Creating standardized issue-fixing workflows
- **Automation**: From issue analysis to PR creation in one command
- **Best Practices**: Security, prompt design, and quality gates

**Perfect for**: Teams wanting to standardize AI-assisted development workflows

### 📝 [Validation Workshop](./VALIDATION_WORKSHOP.md)
A comprehensive guide on implementing input validation in Go applications. This workshop covers:

- **Problem Analysis**: Understanding validation requirements
- **Solution Design**: Architecture and design patterns  
- **Implementation**: Step-by-step coding guide
- **Testing Strategy**: Comprehensive test coverage
- **Best Practices**: Do's and don'ts for Go validation

**Perfect for**: Developers learning Go validation patterns, team workshops, code review guidelines

### 🚀 Quick Start Example

```go
// Before: No validation
func (c *Cart) AddProduct(product Product, quantity int64) {
    c.Items[product.ID] = &CartItem{Product: product, Quantity: quantity}
}

// After: Robust validation
func (c *Cart) AddProduct(product Product, quantity int64) error {
    if err := ValidateProduct(product); err != nil {
        return fmt.Errorf("invalid product: %w", err)
    }
    if err := ValidateQuantity(quantity); err != nil {
        return fmt.Errorf("invalid quantity: %w", err)
    }
    c.Items[product.ID] = &CartItem{Product: product, Quantity: quantity}
    return nil
}
```

## Contributing

When adding new documentation:

1. Follow the existing structure and style
2. Include practical code examples
3. Add comprehensive explanations for complex concepts
4. Update this README with links to new documents

## Related

- [Main README](../README.md) - Project overview
- [GitHub Issues](https://github.com/Pkittipat/try-cart/issues) - Current development tasks