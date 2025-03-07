# Guidelines for Shellcanary Discover CLI

## Build & Run Commands
- Build: `go build -o discover`
- Run: `./discover` or `go run main.go`
- Test: `go test ./...`
- Single test: `go test ./path/to/package -run TestName`
- Lint: `golint ./...`
- Vet: `go vet ./...`

## Code Style Guidelines
- **Imports**: Standard Go style (stdlib first, then third-party, then internal)
- **Formatting**: Use `gofmt` or `go fmt ./...`
- **Types**: Prefer explicit types; use interfaces for abstraction
- **Error Handling**: Check errors explicitly; use `fmt.Errorf("context: %w", err)` for context
- **Naming**: 
  - PascalCase for exported functions/types (GetDockerComposeProjects)
  - camelCase for unexported functions/variables
  - Use descriptive, concise names
- **Comments**: All exported functions must have doc comments that begin with the function name
- **Function Structure**: Single responsibility; limit function size (< 50 lines preferred)
- **Project Structure**: Maintain separation of concerns (agents, models, UI)
- **Dependencies**: Limited use of external packages (promptui for UI)
- **Testing**: Write table-driven tests where appropriate