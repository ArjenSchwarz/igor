# Decision Log: Go Modernization

This document records key decisions made during the requirements gathering and implementation phases of the Igor Go modernization project.

## Decision 1: Feature Scope - Modernization vs. Feature Addition

**Date**: 2025-10-20

**Context**: The project could focus on adding new features OR modernizing existing code.

**Decision**: Focus exclusively on modernization of existing code without adding new features.

**Rationale**:
- The codebase has accumulated significant technical debt with deprecated APIs
- Multiple security issues (HTTP vs HTTPS) need addressing
- Build infrastructure (Wercker) is no longer maintained
- Modernization provides foundation for future feature work
- Reduces scope creep and allows focused delivery

**Alternatives Considered**:
- Combined modernization + new features: Rejected due to complexity and risk of scope creep
- Partial modernization: Rejected as it would leave inconsistent code patterns

---

## Decision 2: Minimum Go Version Target

**Date**: 2025-10-20

**Context**: Need to choose which Go version to target as minimum supported version.

**Decision**: Target Go 1.21 or later as minimum supported version.

**Rationale**:
- Go 1.21 is stable and widely adopted (released August 2023)
- Provides modern standard library APIs we want to use
- Still recent enough that most environments support it
- AWS Lambda Go runtime supports Go 1.21+
- Not cutting edge (avoiding Go 1.23+) to ensure broad compatibility

**Alternatives Considered**:
- Go 1.16: Too old, though it introduced os.ReadFile()
- Go 1.23+: Too new, may not be available in all deployment environments

---

## Decision 3: CI/CD Platform Selection

**Date**: 2025-10-20

**Context**: Wercker CI/CD platform is deprecated. Need to select replacement.

**Decision**: Migrate to GitHub Actions.

**Rationale**:
- Project is hosted on GitHub, making Actions the natural choice
- No additional setup or external accounts required
- Excellent integration with GitHub releases, issues, PRs
- Large ecosystem of pre-built actions
- Free for open source projects
- Industry standard with strong community support

**Alternatives Considered**:
- GitLab CI: Would require mirroring repository or moving to GitLab
- CircleCI: Requires external account, additional complexity
- Travis CI: Has had reliability issues and pricing changes

---

## Decision 4: Global State Refactoring Approach

**Date**: 2025-10-20

**Context**: Config package uses global variables which hurt testability and thread safety.

**Decision**: Refactor to use dependency injection with Config struct/interface while maintaining backwards compatibility with plugin interfaces.

**Rationale**:
- Improves testability significantly
- Enables parallel test execution
- Thread-safe by design
- Plugin interface can remain stable if done carefully
- Aligns with Go best practices

**Constraints**:
- Must maintain backwards compatibility with existing plugins
- Cannot break existing configuration file format
- Should not require plugins to be rewritten

**Alternatives Considered**:
- Leave global state: Rejected due to testing and concurrency issues
- Complete plugin interface redesign: Too invasive, breaks existing plugins

---

## Decision 5: HTTP Client Timeout Strategy

**Date**: 2025-10-20

**Context**: Current HTTP clients have no timeouts, which can cause indefinite hangs.

**Decision**: Add 30-second default timeout for all HTTP clients with context support where feasible.

**Rationale**:
- 30 seconds is reasonable for external API calls
- Prevents indefinite hangs on slow/dead services
- Context support enables request cancellation
- Improves reliability in production

**Implementation Notes**:
- Each plugin's HTTP client should have explicit timeout
- Consider per-plugin timeout customization if needed
- Use context.WithTimeout for better control

**Alternatives Considered**:
- Global HTTP client: Rejected as different plugins may need different timeouts
- No timeout: Rejected due to reliability concerns
- Shorter timeout (10s): May be too aggressive for some APIs

---

## Decision 6: goquery API Migration Strategy

**Date**: 2025-10-20

**Context**: `goquery.NewDocument(url)` is deprecated.

**Decision**: Migrate to `goquery.NewDocumentFromReader()` with explicit HTTP client control.

**Rationale**:
- Gives us control over HTTP client (enables timeout configuration)
- Aligns with modern goquery API
- Allows us to handle HTTP errors properly
- More explicit and testable

**Implementation**:
```go
// Instead of:
doc, err := goquery.NewDocument(url)

// Use:
resp, err := httpClient.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()
doc, err := goquery.NewDocumentFromReader(resp.Body)
```

---

## Decision 7: TODO Comment Resolution Strategy

**Date**: 2025-10-20

**Context**: Multiple TODO comments exist in the codebase.

**Decision**: For each TODO, either fix immediately, create GitHub issue, or remove if obsolete.

**Rationale**:
- Untracked TODOs become invisible technical debt
- GitHub issues provide visibility and can be prioritized
- Some TODOs may be obsolete after other changes
- Forces conscious decision about each item

**Process**:
1. Evaluate if TODO is still relevant
2. If fixable quickly (< 30 min), fix during modernization
3. If larger effort, create GitHub issue and reference in comment
4. If obsolete, remove

---

## Decision 8: Backwards Compatibility Requirements

**Date**: 2025-10-20

**Context**: Need to define what level of backwards compatibility to maintain.

**Decision**: Maintain strict backwards compatibility for:
- Configuration file formats (YAML/JSON)
- Plugin interface
- Environment variables and Lambda integration
- Docker container usage

Allow breaking changes for:
- Build scripts (as long as documented)
- Development tooling
- Internal unexported APIs

**Rationale**:
- Users have existing deployments with configurations
- Plugin authors may have external plugins
- Lambda deployments need to work without code changes
- Internal refactoring shouldn't impact users
- Build process changes can be documented in migration guide

---

## Decision 9: Go Modules Module Path

**Date**: 2025-10-20

**Context**: Need to initialize go.mod with appropriate module path.

**Decision**: Use `github.com/arjenschwarz/igor` as the module path.

**Rationale**:
- Matches the actual GitHub repository location
- Follows Go module naming conventions
- Enables proper import paths for external users
- Allows for go get installation

---

## Decision 10: Testing Strategy for Modernization

**Date**: 2025-10-20

**Context**: Need to ensure modernization doesn't break functionality.

**Decision**:
1. All existing tests must pass before and after changes
2. Manual testing checklist for each plugin
3. Lambda deployment test in AWS
4. Docker container build and run test

**Rationale**:
- Existing tests provide baseline validation
- Some functionality may not be covered by automated tests
- Real-world deployment scenarios must be validated
- Plugin functionality is critical to user experience

**Success Criteria**:
- ✓ All `go test ./...` passes
- ✓ Application deploys to Lambda successfully
- ✓ Docker container builds without errors
- ✓ Manual verification of at least one command from each plugin
- ✓ No new build warnings or errors

---

## Questions for Clarification

None at this time. Requirements are ready for review.

---

## Future Decisions

The following decisions may need to be made during implementation:

1. **Specific dependency versions**: May need to evaluate and choose specific versions when running `go mod init`
2. **GitHub Actions workflow details**: Exact Go versions to test against, release automation details
3. **Makefile inclusion**: Decide if Makefile adds value or if Go commands are sufficient
4. **Error handling patterns**: May discover additional error handling issues during implementation
5. **Channel buffer sizing**: For concurrent goroutines, may need to determine appropriate buffer sizes

These will be documented as they arise during implementation.
