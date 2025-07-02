# Claude Code + GitHub Integration Workshop

This guide demonstrates how to set up Claude Code with GitHub integration using MCP (Model Context Protocol) and custom prompts for standardized issue fixing workflows.

## 📚 Table of Contents

- [Overview](#overview)
- [Setup Configuration](#setup-configuration)
- [Custom Prompts](#custom-prompts)
- [Workflow Demo](#workflow-demo)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

Claude Code can be configured to seamlessly integrate with GitHub for automated issue resolution. This workshop shows how we set up:

1. **MCP GitHub Server** for direct GitHub API access
2. **Custom prompts** for standardized workflows
3. **Automated issue fixing** from identification to PR creation

### What You'll Learn

- How to configure Claude Code with GitHub MCP
- Creating custom prompts for repetitive workflows
- Automating the entire issue-to-PR process
- Best practices for AI-assisted development

## Setup Configuration

### 1. Claude Code Configuration File

Create or edit `.claude_config.json` in your project root:

```json
{
  "mcpServers": {
    "github": {
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "your_github_token_here"
      }
    }
  },
  "prompts": {
    "fix-github-issue": "Please analyze and fix the GitHub issue: $ARGUMENTS. \nFollow these steps:\n1. Use `MCP Github Server` to get the issue details\n2. Understand the problem described in the issue\n3. Search the codebase for relevant files\n4. Implement the necessary changes to fix the issue\n5. Write and run tests to verify the fix\n6. Ensure code passes linting and type checking\n7. Create a descriptive commit message\n8. Push and create a PR\nRemember to use the MCP Github for all GitHub-related tasks."
  }
}
```

### 2. GitHub Token Setup

**Step 1**: Create a GitHub Personal Access Token
- Go to GitHub Settings → Developer settings → Personal access tokens
- Generate a new token with these permissions:
  - `repo` (full repository access)
  - `issues` (read/write issues)
  - `pull_requests` (create/manage PRs)

**Step 2**: Add token to configuration
```json
{
  "mcpServers": {
    "github": {
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "ghp_your_token_here"
      }
    }
  }
}
```

### 3. MCP Server Configuration

The MCP GitHub server enables Claude Code to:
- ✅ Fetch issue details
- ✅ Create branches and commits
- ✅ Push code changes
- ✅ Create pull requests
- ✅ Manage repository operations

**Alternative Docker Setup** (if available):
```json
{
  "mcpServers": {
    "github": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "GITHUB_PERSONAL_ACCESS_TOKEN",
        "mcp/github"
      ],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "your_token_here"
      }
    }
  }
}
```

## Custom Prompts

### The `fix-github-issue` Prompt

Our custom prompt standardizes the entire issue-fixing workflow:

```json
{
  "prompts": {
    "fix-github-issue": "Please analyze and fix the GitHub issue: $ARGUMENTS. \nFollow these steps:\n1. Use `MCP Github Server` to get the issue details\n2. Understand the problem described in the issue\n3. Search the codebase for relevant files\n4. Implement the necessary changes to fix the issue\n5. Write and run tests to verify the fix\n6. Ensure code passes linting and type checking\n7. Create a descriptive commit message\n8. Push and create a PR\nRemember to use the MCP Github for all GitHub-related tasks."
  }
}
```

### Usage Pattern

```bash
# In Claude Code CLI
/fix-github-issue 1
```

This single command triggers a complete workflow that:
1. 🔍 Fetches issue #1 details
2. 📖 Analyzes the problem
3. 🔎 Searches relevant code files
4. ⚡ Implements the solution
5. 🧪 Writes and runs tests
6. ✅ Validates with linting/type checking
7. 📝 Creates descriptive commit
8. 🚀 Pushes and creates PR

### Advanced Prompt Customization

You can create specialized prompts for different types of issues:

```json
{
  "prompts": {
    "fix-github-issue": "Standard issue fixing workflow...",
    "fix-bug": "Focus on bug investigation and fixing: $ARGUMENTS...",
    "add-feature": "Implement new feature from issue: $ARGUMENTS...",
    "refactor-code": "Refactor code based on issue: $ARGUMENTS...",
    "add-tests": "Add comprehensive tests for issue: $ARGUMENTS..."
  }
}
```

## Workflow Demo

### Our Issue #1 Example

**Command Used:**
```bash
/fix-github-issue 1
```

**What Happened:**

#### 1. Issue Analysis
```bash
# Claude Code automatically:
curl -H "Authorization: token ghp_..." \
  https://api.github.com/repos/pkittipat/try-cart/issues/1
```

**Issue Details Retrieved:**
- Title: "Add input validation for cart operations"
- Requirements: Validate quantities, product IDs, error handling
- Acceptance criteria: Specific validation rules

#### 2. Codebase Analysis
```bash
# Automatically searched for relevant files
find . -name "*.go" | grep -E "(cart|product)"
```

**Files Identified:**
- `internal/domain/cart/cart.go` - Main implementation
- `internal/domain/cart/cart_test.go` - Test file
- `internal/domain/cart/product.go` - Product model

#### 3. Implementation
**Added validation functions:**
```go
func ValidateProduct(product Product) error {
    if strings.TrimSpace(product.ID) == "" {
        return errors.New("product ID cannot be empty")
    }
    // ... more validation
}

func ValidateQuantity(quantity int64) error {
    if quantity <= 0 {
        return errors.New("quantity must be a positive integer")
    }
    return nil
}
```

#### 4. Testing
**Generated comprehensive tests:**
- 30+ test cases covering all validation scenarios
- Updated existing tests for backward compatibility
- Verified all tests pass

#### 5. Quality Assurance
```bash
# Automatically ran
go test ./internal/domain/cart/... -v
go vet ./...
go build ./...
```

#### 6. Git Operations
```bash
# Automatically executed
git checkout -b feature/add-cart-validation
git add .
git commit -m "Add input validation for cart operations..."
git push -u origin feature/add-cart-validation
```

#### 7. PR Creation
```bash
# Automatically created PR via GitHub API
curl -X POST \
  -H "Authorization: token ghp_..." \
  https://api.github.com/repos/pkittipat/try-cart/pulls \
  -d '{"title":"Add input validation for cart operations","body":"...","head":"feature/add-cart-validation","base":"main"}'
```

**Result:** [PR #2](https://github.com/Pkittipat/try-cart/pull/2) created automatically!

## Best Practices

### 1. Token Security

#### ✅ DO: Use Environment Variables
```bash
export GITHUB_PERSONAL_ACCESS_TOKEN="your_token"
```

#### ✅ DO: Limit Token Scope
- Only grant necessary permissions
- Use fine-grained tokens when possible

#### ❌ DON'T: Commit Tokens to Git
```bash
# Add to .gitignore
.claude_config.json
```

### 2. Prompt Design

#### ✅ DO: Be Specific
```json
{
  "fix-github-issue": "Please analyze and fix the GitHub issue: $ARGUMENTS. \nFollow these steps:\n1. Use `MCP Github Server` to get the issue details..."
}
```

#### ✅ DO: Include Validation Steps
- Always include testing requirements
- Specify linting and type checking
- Require descriptive commit messages

#### ❌ DON'T: Make Prompts Too Generic
```json
{
  "fix-issue": "Fix the issue: $ARGUMENTS"  // Too vague
}
```

### 3. Workflow Standards

#### ✅ DO: Follow Consistent Branch Naming
```bash
feature/add-cart-validation
bugfix/fix-quantity-validation
refactor/improve-error-handling
```

#### ✅ DO: Use Descriptive Commit Messages
```
Add input validation for cart operations

- Add validation for product quantities (must be positive integers)
- Validate product IDs (must not be empty/nil)
- Add proper error handling and return meaningful error messages
- Ensure cart operations fail gracefully with invalid input
- Comprehensive unit tests cover all validation scenarios

Fixes #1
```

### 4. Quality Gates

Always include these steps in your prompts:
```
5. Write and run tests to verify the fix
6. Ensure code passes linting and type checking
7. Create a descriptive commit message
```

## Troubleshooting

### Common Issues

#### 1. "MCP Server Not Found"
**Problem**: GitHub MCP server not properly configured
**Solution**: 
```json
{
  "mcpServers": {
    "github": {
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "your_token_here"
      }
    }
  }
}
```

#### 2. "Permission Denied" Errors
**Problem**: GitHub token lacks required permissions
**Solution**: Regenerate token with `repo`, `issues`, and `pull_requests` scopes

#### 3. "Command Not Found: /fix-github-issue"
**Problem**: Custom prompt not properly defined
**Solution**: Verify `.claude_config.json` syntax and restart Claude Code

#### 4. API Rate Limiting
**Problem**: Too many GitHub API calls
**Solution**: 
- Use authenticated requests (increases rate limit)
- Add delays between operations if needed

### Debug Configuration

Add debug information to your config:
```json
{
  "mcpServers": {
    "github": {
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "your_token",
        "DEBUG": "true"
      }
    }
  }
}
```

## Advanced Configurations

### Multiple GitHub Accounts
```json
{
  "mcpServers": {
    "github-work": {
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "work_token"
      }
    },
    "github-personal": {
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "personal_token"
      }
    }
  }
}
```

### Project-Specific Prompts
```json
{
  "prompts": {
    "fix-go-issue": "Fix Go-specific issue: $ARGUMENTS\n- Run go vet\n- Run go test\n- Check gofmt",
    "fix-js-issue": "Fix JavaScript issue: $ARGUMENTS\n- Run eslint\n- Run npm test\n- Check prettier",
    "fix-security-issue": "Fix security issue: $ARGUMENTS\n- Run security scan\n- Update dependencies\n- Add security tests"
  }
}
```

## Conclusion

This GitHub + Claude Code integration workflow provides:

1. **🚀 Speed**: From issue to PR in minutes, not hours
2. **🎯 Consistency**: Standardized approach to issue fixing
3. **✅ Quality**: Built-in testing and validation steps
4. **📖 Documentation**: Automatic commit messages and PR descriptions
5. **🔄 Repeatability**: Same high-quality process every time

### Key Takeaways

- **MCP enables direct GitHub API access** from Claude Code
- **Custom prompts standardize workflows** and ensure consistency
- **Single command execution** can handle complex multi-step processes
- **Quality gates ensure** robust solutions every time

This approach transforms issue fixing from a manual, error-prone process into a reliable, automated workflow that maintains high code quality standards.

---

*This documentation demonstrates the Claude Code + GitHub integration used to fix issue #1 in the try-cart project.*