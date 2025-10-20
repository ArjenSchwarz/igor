# Go Modernization Requirements

## Introduction

This specification outlines the requirements for modernizing the Igor Slack bot codebase to align with current Go best practices, improve security, and enhance maintainability. The Igor project is a Slack slash command bot framework written in Go that runs on AWS Lambda or as a Docker container. Over time, the codebase has accumulated deprecated patterns, outdated dependencies, and infrastructure configurations that need to be updated.

This modernization effort focuses on three primary areas:
1. **Code Modernization**: Replacing deprecated Go patterns and improving code quality
2. **Security Improvements**: Updating HTTP URLs to HTTPS and fixing resource management issues
3. **Infrastructure Modernization**: Migrating from deprecated build tools to modern CI/CD platforms

## Requirements

### 1. Deprecated Go API Replacement

**User Story**: As a maintainer, I want all deprecated Go standard library APIs replaced with their modern equivalents, so that the codebase remains compatible with future Go versions and follows current best practices.

**Acceptance Criteria**:

1.1. **io/ioutil Replacement**
   - The system SHALL replace all instances of `ioutil.ReadFile()` with `os.ReadFile()`
   - The system SHALL replace all instances of `ioutil.ReadDir()` with `os.ReadDir()`
   - The system SHALL remove the `io/ioutil` import from all files where it's no longer needed
   - The system SHALL maintain identical functionality after replacement
   - The system SHALL handle errors in the same manner as the original implementation

1.2. **strings.Replace Modernization**
   - The system SHALL replace all instances of `strings.Replace(str, old, new, -1)` with `strings.ReplaceAll(str, old, new)`
   - The system SHALL verify that all replacement operations produce identical results
   - The system SHALL update at minimum the following files:
     - `config/config.go`
     - `slack/response.go`
     - `plugins/status.go`
     - Any other files containing this pattern

1.3. **goquery API Updates**
   - The system SHALL replace all instances of `goquery.NewDocument(url)` with `goquery.NewDocumentFromReader()` or equivalent modern API
   - The system SHALL properly handle HTTP response bodies when using the new API
   - The system SHALL update at minimum the following files:
     - `plugins/status.go` (lines 278, 323, 364, 381 and surrounding context)

### 2. Dependency Management Modernization

**User Story**: As a developer, I want the project to use Go modules for dependency management, so that dependencies are explicitly tracked, versioned, and reproducible across environments.

**Acceptance Criteria**:

2.1. **Go Modules Initialization**
   - The system SHALL create a `go.mod` file with appropriate module path
   - The system SHALL create a `go.sum` file containing cryptographic checksums of dependencies
   - The system SHALL include all current dependencies with their versions:
     - `github.com/aws/aws-lambda-go`
     - `github.com/aws/aws-sdk-go`
     - `gopkg.in/yaml.v2`
     - `github.com/PuerkitoBio/goquery`
     - All transitive dependencies

2.2. **Dependency Version Selection**
   - The system SHALL use stable, non-deprecated versions of all dependencies
   - The system SHALL document any breaking changes requiring code modifications
   - The system SHALL verify that all tests pass with the selected versions

### 3. Security Enhancements

**User Story**: As a security-conscious operator, I want all external API communications to use HTTPS, so that data in transit is encrypted and protected from interception.

**Acceptance Criteria**:

3.1. **HTTP to HTTPS Migration**
   - The system SHALL replace `http://api.openweathermap.org` with `https://api.openweathermap.org` in `plugins/weather.go`
   - The system SHALL replace `http://xkcd.com/` with `https://xkcd.com/` in `plugins/xkcd.go`
   - The system SHALL verify that all HTTPS endpoints are functional
   - The system SHALL maintain backwards compatibility with existing plugin configurations

3.2. **HTTP Client Timeout Configuration**
   - The system SHALL add timeout configurations to all HTTP clients
   - The system SHALL use context-aware HTTP requests where appropriate
   - The system SHALL prevent indefinite hangs on slow external services
   - The system SHALL configure reasonable timeout values (e.g., 30 seconds for API calls)

### 4. Resource Management Improvements

**User Story**: As a reliability engineer, I want proper resource cleanup in all error paths, so that the application doesn't leak file descriptors or memory.

**Acceptance Criteria**:

4.1. **HTTP Response Body Cleanup**
   - The system SHALL ensure `defer resp.Body.Close()` is called only after nil checks
   - The system SHALL fix the resource leak in `plugins/status.go` line 195
   - The system SHALL audit all HTTP request handlers for similar issues
   - The system SHALL ensure response bodies are closed even in error paths

4.2. **Error Handling Consistency**
   - The system SHALL uncomment and properly handle errors in `plugins/status.go` lines 99-100
   - The system SHALL ensure all error cases are logged or returned appropriately
   - The system SHALL not silently ignore errors that could indicate problems

### 5. Global State Refactoring

**User Story**: As a developer writing tests, I want configuration to be dependency-injected rather than stored in global variables, so that tests can run in parallel without interference and are easier to write.

**Acceptance Criteria**:

5.1. **Config Package Refactoring**
   - The system SHALL refactor `config.go` to eliminate global variables:
     - `configHolder`
     - `configFile`
     - `jsonConfig`
     - `fallbackLanguage`
   - The system SHALL create a Config struct/interface that can be passed to components
   - The system SHALL maintain backwards compatibility with existing plugin interfaces
   - The system SHALL ensure thread-safety for concurrent access

5.2. **Main Package Cleanup**
   - The system SHALL refactor the global `servervar` flag in `main.go`
   - The system SHALL use proper flag parsing and configuration structures
   - The system SHALL maintain existing CLI behavior

### 6. CI/CD Modernization

**User Story**: As a maintainer, I want the project to use a modern, supported CI/CD platform, so that automated builds and tests continue to function reliably.

**Acceptance Criteria**:

6.1. **GitHub Actions Migration**
   - The system SHALL create a `.github/workflows` directory
   - The system SHALL create a CI workflow that:
     - Runs on pull requests and pushes to main branch
     - Tests against multiple Go versions (e.g., 1.21, 1.22, latest)
     - Runs `go test ./...`
     - Runs `go build`
     - Checks code formatting with `gofmt`
   - The system SHALL create a release workflow that:
     - Builds binaries for multiple platforms (linux/darwin, amd64/arm64)
     - Creates Docker images
     - Publishes releases to GitHub

6.2. **Wercker Removal**
   - The system SHALL remove or archive `wercker.yml`
   - The system SHALL update README.md to remove Wercker references
   - The system SHALL document the new CI/CD process in the README

### 7. Build Process Modernization

**User Story**: As a contributor, I want a simple and standard build process, so that I can easily build and test the project locally.

**Acceptance Criteria**:

7.1. **Build Script Updates**
   - The system SHALL update `bin/build.sh` to remove Wercker-specific commands
   - The system SHALL use standard Go build commands
   - The system SHALL support cross-compilation for Lambda deployment
   - The system SHALL maintain Docker build functionality

7.2. **Makefile Creation (Optional Enhancement)**
   - The system MAY create a Makefile with common targets:
     - `make build`: Build the binary
     - `make test`: Run tests
     - `make lint`: Run linters
     - `make docker`: Build Docker image
     - `make clean`: Clean build artifacts

### 8. Code Quality Improvements

**User Story**: As a maintainer, I want TODO comments to be resolved or converted to tracked issues, so that technical debt is properly managed.

**Acceptance Criteria**:

8.1. **TODO Resolution**
   - The system SHALL address or document the following TODOs:
     - `handler.go` line 39: "TODO clean this up"
     - `slack/response.go` lines 55, 65: "TODO validation"
     - `plugins/plugins.go` line 32: Error handling improvements
   - The system SHALL either:
     - Fix the TODO item, OR
     - Create a GitHub issue and reference it in the comment, OR
     - Remove the TODO if no longer applicable

8.2. **Concurrent Goroutine Safety**
   - The system SHALL review goroutine usage in:
     - `plugins/help.go` lines 111-127
     - `plugins/status.go` lines 94-106
   - The system SHALL add panic recovery where appropriate
   - The system SHALL consider adding timeout contexts to prevent goroutine leaks
   - The system SHALL ensure proper channel buffer sizing or select with timeout

## Technical Constraints

- **Go Version**: The modernization SHALL target Go 1.21 or later as the minimum supported version
- **Backwards Compatibility**: The modernization SHALL maintain compatibility with existing configuration files (YAML/JSON)
- **Plugin Interface**: The modernization SHALL NOT break the existing plugin interface without providing migration guidance
- **AWS Lambda**: The modernization SHALL ensure continued compatibility with AWS Lambda runtime
- **Docker**: The modernization SHALL maintain the ability to run in Docker containers

## Edge Cases and Considerations

1. **Configuration Migration**: Existing users may have encrypted configurations via KMS. The modernization must not break decryption functionality.

2. **External API Changes**: When migrating HTTP to HTTPS, some APIs may have different response formats or authentication requirements on the secure endpoint.

3. **Lambda Cold Starts**: Changes to dependency management or initialization code should not significantly impact Lambda cold start times.

4. **Concurrent Plugin Execution**: Global state refactoring must consider that plugins may be executed concurrently in response to multiple Slack requests.

5. **Tumblr Plugin**: The `randomtumblr.go` plugin may need special attention as the Tumblr API may have changed since the code was written.

## Success Criteria

The modernization will be considered complete when:

1. All automated tests pass with the updated code
2. The application successfully deploys to AWS Lambda
3. The Docker container builds and runs successfully
4. All deprecated Go APIs have been replaced
5. GitHub Actions CI/CD is functional and passing
6. Manual testing confirms all plugins work as expected
7. No new warnings or errors are introduced in the build process

## Out of Scope

The following items are explicitly out of scope for this modernization effort:

- Adding new features or plugins
- Changing the Slack API integration approach
- Migrating from AWS Lambda to other serverless platforms
- Comprehensive test coverage improvements (beyond ensuring existing tests pass)
- Performance optimization beyond fixing obvious issues
- Migration from gopkg.in/yaml.v2 to v3 (unless required by other changes)
