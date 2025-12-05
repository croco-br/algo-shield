# Contributing to AlgoShield

First off, thank you for considering contributing to AlgoShield! It's people like you that make AlgoShield such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our commitment to providing a welcoming and inspiring community for all.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- Use a clear and descriptive title
- Describe the exact steps to reproduce the problem
- Provide specific examples to demonstrate the steps
- Describe the behavior you observed and what behavior you expected
- Include logs and screenshots if applicable

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- A clear and descriptive title
- A detailed description of the proposed enhancement
- Examples of how the enhancement would be used
- Why this enhancement would be useful to most AlgoShield users

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. Ensure the test suite passes
4. Make sure your code follows the existing style
5. Write a clear commit message

## Development Setup

1. Install Go 1.23 or later
2. Install Node.js 20 or later
3. Install Docker and Docker Compose
4. Clone the repository
5. Run `make deps` to install dependencies

## Running Tests

```bash
make test
```

## Building

```bash
make build
```

## Code Style

### Go

- Follow standard Go conventions
- Use `gofmt` to format your code
- Run `golangci-lint` before submitting

### TypeScript/Svelte

- Use consistent indentation (tabs)
- Follow SvelteKit conventions

## Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters
- Reference issues and pull requests after the first line

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

