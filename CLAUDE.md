# Project Configuration

## Basic rules

1. Prioritize std go packages to third-party
2. Use idiomatic go patterns
3. Avoid `panic` in situations when simple `result, err := ...` would work
4. Avoid putting too much of the defensive code, like checking each `err` coming from Writers
5. Avoid inline comments, keep them on the function level
   - a too complex function should be refactored instead of putting loads of inline comments
   - not all functions need `/// docs` but most of the non-trivial functions do

## Project Structure

1. Follow the established package layout:
   - `cmd/gitcat/` - CLI application entry point and argument parsing
   - `pkg/` - Reusable library packages with clear responsibilities
   - `pkg/internal/` - Internal utilities not exposed externally
2. Each package should have a single, well-defined purpose
3. Avoid circular dependencies between packages

## Error Handling

1. Use error wrapping with `fmt.Errorf("context: %w", err)` to maintain error chains
2. Log errors at the point where they're handled, not where they're returned
3. For cleanup operations, use silent error handling patterns (e.g., `utils.SilentClose()`)
4. Validate inputs early and return errors with descriptive messages

## Logging

1. Use structured logging with `slog` package
2. Include relevant context as key-value pairs: `log.Info("message", "key", value)`
3. Use appropriate log levels:
   - `Debug` for detailed trace information
   - `Info` for normal operations
   - `Warn` for non-critical issues
   - `Error` for failures
4. Set logger at package level using dependency injection pattern

## Concurrency

1. Use goroutines with proper synchronization (e.g., `sync.WaitGroup`, `sync.Mutex`)
2. Always ensure goroutines complete before returning (use `wg.Wait()`)
3. Protect shared state with mutexes when accessed from multiple goroutines
4. Prefer channels for complex coordination, mutexes for simple shared state

## Resource Management

1. Always use `defer` for cleanup operations (file closes, temp directory removal)
2. Create helper functions for common cleanup patterns
3. For temporary resources, clean up immediately after use or in defer blocks
4. Handle cleanup errors silently when appropriate to avoid masking primary errors

## CLI Design

1. Use `flag` package for argument parsing
2. Provide usage information via custom `fs.Usage` functions
3. Validate arguments after parsing and before execution
4. Support dry-run mode for testing without side effects

## Output Formatting

1. Support multiple output formats (JSON, text) via command-line flags
2. Use proper JSON marshaling for structured output
3. Keep output formatting logic separate from business logic

## Tools

1. Do not make web requests to arbitrary pages, only standard Golang repos and places

## Code Style

1. Follow `go fmt` rules
2. Use method chaining for builder-style APIs (e.g., `NewList().IgnoreDotFiles()`)
3. Return pointers for struct methods that modify state
4. Keep functions focused and single-purpose
