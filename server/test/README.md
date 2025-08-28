# Integration Tests for My Collection Server

This directory contains comprehensive integration tests that verify end-to-end functionality across multiple packages using real databases and filesystems.

## 🏗️ Structure

```
server/test/
├── integration/           # Integration test packages
│   ├── fssync_test.go     # Filesystem synchronization tests
│   ├── autotags_test.go   # AutoTags behavior tests
│   └── stress_test.go     # Stress and performance tests
├── testutils/             # Shared testing utilities
│   └── framework.go       # IntegrationTestFramework
├── scripts/               # Test runners and utilities
│   └── run_integration_tests.sh
└── README.md              # This file
```

## 🚀 Quick Start

### Run All Tests
```bash
# From the server directory
cd server
./test/scripts/run_integration_tests.sh
```

### Run Specific Test Suites
```bash
# Filesystem synchronization tests only
./test/scripts/run_integration_tests.sh fssync

# AutoTags tests only
./test/scripts/run_integration_tests.sh autotags

# Stress tests only (may take 10+ minutes)
./test/scripts/run_integration_tests.sh stress
```

### Run with Options
```bash
# Verbose output with coverage
./test/scripts/run_integration_tests.sh -v -c

# Quick tests only (skip long-running tests)
./test/scripts/run_integration_tests.sh --short

# Specific test pattern
./test/scripts/run_integration_tests.sh -p "TestBasic*"

# Include stress tests (disabled by default)
./test/scripts/run_integration_tests.sh -s
```

## 📋 Test Suites

### 1. Filesystem Synchronization (`fssync_test.go`)
Tests the core filesystem synchronization functionality:
- **Basic Operations**: File/directory creation, deletion, moves
- **Complex Hierarchies**: Deep nesting, bulk operations, reorganization
- **Stale Handling**: Recovery from inconsistent filesystem/DB states
- **Performance**: Sync time measurement and regression detection

**Key Tests:**
- `TestBasicFileOperations` - Basic file add/remove/move
- `TestComplexHierarchyOperations` - Complex nested directory operations
- `TestStaleFileHandling` - Handling of orphaned database entries
- `TestRealisticLibraryOperations` - Real media library scenarios

### 2. AutoTags Integration (`autotags_test.go`)
Comprehensive testing of AutoTags behavior during filesystem operations:
- **Creation**: AutoTags automatically created for directories
- **Updates**: AutoTags updated when files/directories move
- **Cleanup**: AutoTags removed when directories are deleted
- **Consistency**: AutoTags remain consistent during complex operations
- **Interaction**: AutoTags coexist with manual tags

**Key Tests:**
- `TestAutoTagsBasicBehavior` - Basic AutoTag creation and assignment
- `TestAutoTagsFileMovement` - AutoTag updates during file moves
- `TestAutoTagsDirectoryRename` - AutoTag updates during directory renames
- `TestAutoTagsConsistencyAfterComplexOperations` - Complex reorganization scenarios

### 3. Stress & Performance (`stress_test.go`)
High-load testing and edge case scenarios:
- **Scale Testing**: 1000+ files across 100+ directories
- **Deep Hierarchies**: 20+ level directory nesting
- **Rapid Changes**: Rapid file creation/modification/deletion cycles
- **Edge Cases**: Special characters, long filenames, circular operations
- **Performance**: Regression detection and benchmarking

**Key Tests:**
- `TestMassiveFileOperations` - Large-scale operations (1000+ files)
- `TestDeepHierarchy` - Very deep directory nesting
- `TestSpecialCharacters` - Files with special characters
- `TestPerformanceRegression` - Performance monitoring

## 🛠️ Testing Framework

### IntegrationTestFramework
The core testing framework provides:

```go
// Create isolated test environment
framework := testutils.NewIntegrationTestFramework(t)
defer framework.Cleanup()

// Filesystem operations
framework.CreateFile("path/to/file.mp4", "content")
framework.CreateDir("path/to/directory")
framework.MoveFile("src/file.mp4", "dst/file.mp4")
framework.DeleteFile("path/to/file.mp4")

// Synchronization
framework.Sync()

// Verification
framework.AssertItemExists("directory", "filename.mp4")
framework.AssertDirectoryExists("path")
framework.AssertAutoTagsExist("directory", []string{"expected", "autotags"})
```

### Key Features
- **Isolation**: Each test gets fresh temporary directory and database
- **Real Operations**: Uses actual SQLite database and filesystem
- **AutoTags Testing**: Comprehensive AutoTags verification
- **Performance Monitoring**: Built-in timing and regression detection
- **Automatic Cleanup**: No leftover files or databases

## 📊 Coverage and Reporting

### Generate Coverage Reports
```bash
./test/scripts/run_integration_tests.sh -c
# Generates: coverage_integration.html
```

### Coverage Goals
- **Core Sync Logic**: >90%
- **AutoTags Functionality**: >95%
- **Error Handling**: >80%
- **Edge Cases**: >70%

## 🔧 Adding New Tests

### 1. Adding Tests to Existing Suites

```go
// In fssync_test.go, autotags_test.go, or stress_test.go
func TestYourNewFeature(t *testing.T) {
    framework := testutils.NewIntegrationTestFramework(t)
    defer framework.Cleanup()

    // Your test logic here
    framework.CreateFile("test.mp4", "content")
    framework.Sync()
    framework.AssertItemExists("", "test.mp4")
}
```

### 2. Creating New Test Suites

1. **Create new test file**: `server/test/integration/newsuite_test.go`
2. **Add test pattern**: Update `run_integration_tests.sh` to include your patterns
3. **Document the suite**: Update this README

### 3. Extending the Framework

Add new helper methods to `testutils/framework.go`:

```go
// Example: Add method for complex operations
func (f *IntegrationTestFramework) CreateComplexLibrary(config LibraryConfig) {
    // Implementation
}
```

## 🚨 Troubleshooting

### Common Issues

**"directory not found in db"**
- Database not properly initialized
- Solution: `rm -rf /tmp/integration-test-*`

**"permission denied"**
- Filesystem permission issues
- Solution: Check temp directory permissions

**Test timeouts**
- Performance issues or infinite loops
- Solution: Run with `-v` to see detailed execution

**"coverage_*.out not found"**
- Coverage files not generated
- Solution: Ensure tests complete successfully before coverage generation

### Debugging Tips

1. **Verbose Mode**: Use `-v` flag to see detailed test execution
2. **Isolation**: Each test runs independently - no shared state
3. **Cleanup**: Automatic cleanup on test completion or failure
4. **Logging**: Framework provides detailed logging for debugging

### Performance Considerations

- **Stress Tests**: May require several GB disk space temporarily
- **Timeouts**: Default 45-minute timeout for full test suite
- **Parallel Execution**: Tests run sequentially for database isolation
- **Memory Usage**: Monitor during large-scale tests

## 🎯 Best Practices

### Test Design
1. **Arrange-Act-Assert**: Set up state, perform operations, verify results
2. **Isolation**: Don't depend on other tests
3. **Cleanup**: Use `defer framework.Cleanup()`
4. **Realistic Scenarios**: Test real-world usage patterns

### Performance
1. **Use `testing.Short()`**: Skip long tests in quick mode
2. **Monitor Timing**: Use framework timing helpers
3. **Realistic Data**: Use representative file sizes and counts
4. **Memory Awareness**: Monitor memory usage in stress tests

### Documentation
1. **Test Names**: Use descriptive names (`TestAutoTagsFileMovement`)
2. **Comments**: Explain complex test scenarios
3. **Examples**: Provide usage examples for new framework features

## 🔄 CI/CD Integration

### Quick Tests (Pull Requests)
```bash
./test/scripts/run_integration_tests.sh --short -v
```

### Full Tests (Main Branch)
```bash
./test/scripts/run_integration_tests.sh -s -c
```

### Timeout Protection
```bash
timeout 45m ./test/scripts/run_integration_tests.sh stress
```

## 📈 Metrics and Monitoring

The integration tests provide several metrics:

### Performance Metrics
- **Sync Duration**: Time to synchronize filesystem changes
- **Operation Throughput**: Files processed per second
- **Memory Usage**: Peak memory consumption during operations
- **Database Performance**: Query execution times

### Quality Metrics
- **Test Coverage**: Code coverage across integration scenarios
- **Success Rate**: Percentage of passing tests over time
- **Regression Detection**: Performance degradation alerts

## 🚀 Future Enhancements

Planned improvements:
1. **Concurrent Testing**: Multi-threaded filesystem operations
2. **Network Storage**: Testing with remote filesystems
3. **Large File Testing**: Multi-GB file handling
4. **Real-time Monitoring**: Filesystem watcher integration
5. **Cross-platform Testing**: Windows/macOS compatibility tests

## 📞 Support

For questions or issues:
1. Check test output for specific error messages
2. Run with verbose flag for detailed execution logs
3. Verify filesystem permissions and disk space
4. Review recent changes to sync logic or database schema
5. Check the main project documentation

---

**Note**: These integration tests complement unit tests by verifying end-to-end functionality. They use real databases and filesystems to catch issues that mocks might miss. 