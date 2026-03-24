# Contributing to SAT Jade Backend

First off, thank you for taking the time to contribute! It is people like you who make the developer community such a great place.

## Getting Started

1. **Fork the repository** and create your branch from `main`.
2. **Set up your environment** following the instructions in the `README.md`.
3. Ensure you have **Go 1.24.1+** installed locally if you aren't using Docker for development.

## Style Guidelines

- **Standard Go:** Follow [Effective Go](https://golang.org/doc/effective_go) and use `gofmt` before committing.
- **Commit Messages:** Use descriptive, imperative titles (e.g., `feat: add AI feedback timeout` instead of `fixed stuff`).
- **Performance:** Since this handles high-concurrency mock exams, avoid unnecessary allocations in hot paths.

## Testing

We value stability. Before submitting a PR:
- Ensure all existing tests pass: `go test ./...`
- Add new tests for any features or bug fixes.
- If you're modifying the database schema, include the necessary migration files.

## Pull Request Process

1. Link any relevant issues in your PR description.
2. Ensure your code builds successfully in the CI environment.
3. Once the maintainers approve your changes, they will be merged into `main`.

## License

By contributing, you agree that your contributions will be licensed under the project's **License**.
