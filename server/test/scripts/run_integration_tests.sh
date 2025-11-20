#!/bin/bash

# Integration Test Runner for My Collection Server
# Runs comprehensive integration tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Default values
VERBOSE=true
STRESS_TESTS=false
COVERAGE=false
NO_CACHE=false
PACKAGE="./test/integration"
TEST_PATTERN=""
SUITE=""

print_usage() {
    echo "Usage: $0 [OPTIONS] [SUITE]"
    echo ""
    echo "Test Suites:"
    echo "  fssync           Run filesystem synchronization tests"
    echo "  autotags         Run AutoTags integration tests"
    echo "  server           Run server API integration tests"
    echo "  tags             Run tags API integration tests"
    echo "  stress           Run stress and performance tests"
    echo "  all              Run all integration tests (default)"
    echo ""
    echo "Options:"
    echo "  -h, --help       Show this help message"
    echo "  -v, --verbose    Enable verbose output"
    echo "  -s, --stress     Include stress tests (may take 10+ minutes)"
    echo "  -c, --coverage   Generate test coverage report"
    echo "  -p, --pattern    Run specific test pattern (e.g., 'TestBasic*')"
    echo "  --short          Run tests in short mode (skip long-running tests)"
    echo "  --no-cache       Disable test caching (run all tests fresh)"
    echo ""
    echo "Examples:"
    echo "  $0                           # Run all integration tests"
    echo "  $0 fssync                    # Run only filesystem sync tests"
    echo "  $0 autotags -v               # Run AutoTags tests with verbose output"
    echo "  $0 server                    # Run only server API tests"
    echo "  $0 tags                      # Run only tags API tests"
    echo "  $0 stress                    # Run only stress tests"
    echo "  $0 -v -c                     # Run all tests with verbose output and coverage"
    echo "  $0 -p 'TestBasic*'           # Run only basic tests"
    echo "  $0 --short                   # Run quick tests only"
}

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_suite() {
    echo -e "${PURPLE}[SUITE]${NC} $1"
}

check_dependencies() {
    log_info "Checking dependencies..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check if we're in the right directory
    if [ ! -f "go.mod" ]; then
        log_error "go.mod not found. Please run this script from the server directory."
        exit 1
    fi
    
    # Check if integration test directory exists
    if [ ! -d "test/integration" ]; then
        log_error "Integration test directory not found. Please ensure test/integration exists."
        exit 1
    fi
    
    log_success "Dependencies check passed"
}

build_test_command() {
    local test_file=$1
    local run_pattern=$2
    
    local cmd="go test"
    
    if [ "$VERBOSE" = true ]; then
        cmd="$cmd -v"
    fi
    
    if [ "$NO_CACHE" = true ]; then
        cmd="$cmd -count=1"
    fi
    
    if [ "$COVERAGE" = true ]; then
        local coverage_file="coverage_${test_file}.out"
        cmd="$cmd -coverprofile=$coverage_file -coverpkg=./..."
    fi
    
    if [ "$SHORT_MODE" = true ]; then
        cmd="$cmd -short"
    fi
    
    if [ -n "$run_pattern" ]; then
        cmd="$cmd -run '$run_pattern'"
    fi
    
    cmd="$cmd $PACKAGE -timeout 30m"
    
    echo "$cmd"
}

run_fssync_tests() {
    log_suite "Running Filesystem Synchronization Tests"
    
    local pattern=""
    if [ -n "$TEST_PATTERN" ]; then
        pattern="$TEST_PATTERN"
    else
        pattern="TestBasic.*|TestFile.*|TestDirectory.*|TestComplex.*|TestStale.*|TestMixed.*|TestRealistic.*|TestSync.*"
    fi
    
    local cmd=$(build_test_command "fssync" "$pattern")
    eval $cmd
    
    if [ $? -eq 0 ]; then
        log_success "Filesystem synchronization tests passed"
        return 0
    else
        log_error "Filesystem synchronization tests failed"
        return 1
    fi
}

run_autotags_tests() {
    log_suite "Running AutoTags Integration Tests"
    
    local pattern=""
    if [ -n "$TEST_PATTERN" ]; then
        pattern="$TEST_PATTERN"
    else
        pattern="TestAutoTags.*"
    fi
    
    local cmd=$(build_test_command "autotags" "$pattern")
    eval $cmd
    
    if [ $? -eq 0 ]; then
        log_success "AutoTags integration tests passed"
        return 0
    else
        log_error "AutoTags integration tests failed"
        return 1
    fi
}

run_server_tests() {
    log_suite "Running Server API Integration Tests"
    
    local pattern=""
    if [ -n "$TEST_PATTERN" ]; then
        pattern="$TEST_PATTERN"
    else
        pattern="TestServer.*"
    fi
    
    local cmd=$(build_test_command "server" "$pattern")
    eval $cmd
    
    if [ $? -eq 0 ]; then
        log_success "Server API integration tests passed"
        return 0
    else
        log_error "Server API integration tests failed"
        return 1
    fi
}

run_tags_tests() {
    log_suite "Running Tags API Integration Tests"
    
    local pattern=""
    if [ -n "$TEST_PATTERN" ]; then
        pattern="$TEST_PATTERN"
    else
        pattern="TestTagsServer.*"
    fi
    
    local cmd=$(build_test_command "tags" "$pattern")
    eval $cmd
    
    if [ $? -eq 0 ]; then
        log_success "Tags API integration tests passed"
        return 0
    else
        log_error "Tags API integration tests failed"
        return 1
    fi
}

run_stress_tests() {
    log_suite "Running Stress and Performance Tests"
    log_warning "Stress tests may take 10+ minutes and consume significant disk space temporarily"
    
    local pattern=""
    if [ -n "$TEST_PATTERN" ]; then
        pattern="$TEST_PATTERN"
    else
        pattern="TestMassive.*|TestDeep.*|TestRapid.*|TestSpecial.*|TestConcurrent.*|TestInconsistent.*|TestPerformance.*|TestEdge.*|TestCircular.*|TestVeryLong.*|TestEmpty.*|TestFileReplacement.*|TestCase.*"
    fi
    
    local cmd=$(build_test_command "stress" "$pattern")
    eval $cmd
    
    if [ $? -eq 0 ]; then
        log_success "Stress tests passed"
        return 0
    else
        log_error "Stress tests failed"
        return 1
    fi
}

run_custom_pattern() {
    local pattern=$1
    log_suite "Running tests matching pattern: $pattern"
    
    local cmd=$(build_test_command "custom" "$pattern")
    eval $cmd
    
    if [ $? -eq 0 ]; then
        log_success "Custom pattern tests passed"
        return 0
    else
        log_error "Custom pattern tests failed"
        return 1
    fi
}

run_all_tests() {
    log_suite "Running All Integration Tests"
    
    local cmd="go test"
    if [ "$VERBOSE" = true ]; then
        cmd="$cmd -v"
    fi
    if [ "$NO_CACHE" = true ]; then
        cmd="$cmd -count=1"
    fi
    if [ "$COVERAGE" = true ]; then
        cmd="$cmd -coverprofile=coverage_integration_all.out -coverpkg=./..."
    fi
    if [ "$SHORT_MODE" = true ]; then
        cmd="$cmd -short"
    fi
    if [ "$STRESS_TESTS" = false ] && [ "$SHORT_MODE" != true ]; then
        cmd="$cmd -short"  # Skip stress tests by default unless explicitly enabled
    fi
    if [ -n "$TEST_PATTERN" ]; then
        cmd="$cmd -run '$TEST_PATTERN'"
    fi
    
    cmd="$cmd $PACKAGE -timeout 45m"
    eval $cmd
    
    if [ $? -eq 0 ]; then
        log_success "All integration tests passed"
        return 0
    else
        log_error "Some integration tests failed"
        return 1
    fi
}

generate_coverage_report() {
    if [ "$COVERAGE" = true ]; then
        log_info "Generating coverage report..."
        
        # Find all coverage files
        coverage_files=$(find . -name "coverage_*.out" 2>/dev/null | grep -v merged)
        
        if [ -n "$coverage_files" ]; then
            # Merge coverage files if multiple exist
            if [ $(echo "$coverage_files" | wc -l) -gt 1 ]; then
                log_info "Merging multiple coverage files..."
                echo "mode: set" > coverage_integration_merged.out
                for file in $coverage_files; do
                    tail -n +2 "$file" >> coverage_integration_merged.out
                done
                coverage_file="coverage_integration_merged.out"
            else
                coverage_file="$coverage_files"
            fi
            
            # Generate HTML report
            go tool cover -html="$coverage_file" -o coverage_integration.html
            
            # Show coverage percentage
            if command -v go &> /dev/null; then
                coverage_percent=$(go tool cover -func="$coverage_file" | grep total | awk '{print $3}')
                log_success "Coverage report generated: coverage_integration.html"
                log_info "Total coverage: $coverage_percent"
            fi
            
            # Cleanup individual coverage files (keep the main one for reference)
            # Remove only the per-suite coverage files, keep coverage_integration_all.out
            find . -name "coverage_*.out" ! -name "coverage_integration_all.out" ! -name "coverage_integration_merged.out" -delete 2>/dev/null || true
        else
            log_warning "No coverage files found"
        fi
    fi
}

cleanup() {
    log_info "Cleaning up temporary files..."
    # Remove any temporary test databases or directories that might be left over
    find /tmp -name "integration-test-*" -type d -exec rm -rf {} + 2>/dev/null || true
    find . -name "*.db-shm" -o -name "*.db-wal" -exec rm -f {} + 2>/dev/null || true
}

show_test_summary() {
    echo ""
    log_info "Integration Test Summary:"
    echo "  📁 Location: server/test/integration/"
    echo "  🔧 Framework: Real database and filesystem operations"
    echo "  🏷️  AutoTags: Comprehensive testing of automatic tag behavior"
    echo "  🌐 Server API: Real HTTP server with end-to-end testing"
    echo "  🚀 Performance: Stress tests with 1000+ files and deep hierarchies"
    echo "  🔍 Edge Cases: Special characters, long filenames, circular operations"
    echo "  📊 Coverage: Integration test coverage across multiple packages"
    echo ""
    if [ "$COVERAGE" = true ]; then
        echo "  📈 Coverage Report: coverage_integration.html"
        echo ""
    fi
    echo "  💡 Next Steps:"
    echo "    - Run 'server/test/scripts/run_integration_tests.sh fssync' for fs tests only"
    echo "    - Run 'server/test/scripts/run_integration_tests.sh autotags' for AutoTags tests"
    echo "    - Run 'server/test/scripts/run_integration_tests.sh server' for server API tests"
    echo "    - Run 'server/test/scripts/run_integration_tests.sh tags' for tags API tests"
    echo "    - Run 'server/test/scripts/run_integration_tests.sh stress' for stress testing"
    echo "    - Add new integration tests to server/test/integration/"
}

# Trap to ensure cleanup on exit
trap cleanup EXIT

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            print_usage
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -s|--stress)
            STRESS_TESTS=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -p|--pattern)
            TEST_PATTERN="$2"
            shift 2
            ;;
        --short)
            SHORT_MODE=true
            shift
            ;;
        --no-cache)
            NO_CACHE=true
            shift
            ;;
        fssync|autotags|server|tags|stress|all)
            SUITE="$1"
            shift
            ;;
        *)
            log_error "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Set default suite if none specified
if [ -z "$SUITE" ]; then
    SUITE="all"
fi

# Main execution
main() {
    echo "🧪 My Collection Integration Test Runner"
    echo "========================================"
    log_info "Suite: $SUITE"
    log_info "Package: $PACKAGE"
    if [ "$VERBOSE" = true ]; then
        log_info "Verbose output enabled"
    fi
    if [ "$COVERAGE" = true ]; then
        log_info "Coverage reporting enabled"
    fi
    if [ "$STRESS_TESTS" = true ]; then
        log_info "Stress tests enabled"
    fi
    if [ "$SHORT_MODE" = true ]; then
        log_info "Short mode enabled (skipping long-running tests)"
    fi
    if [ "$NO_CACHE" = true ]; then
        log_info "Test caching disabled (all tests will run fresh)"
    fi
    if [ -n "$TEST_PATTERN" ]; then
        log_info "Test pattern: $TEST_PATTERN"
    fi
    echo ""
    
    check_dependencies
    
    # Set test environment
    export GO_TEST_TIMEOUT="45m"
    export GOMAXPROCS=$(nproc)
    
    # Run tests based on suite
    case "$SUITE" in
        "fssync")
            run_fssync_tests
            ;;
        "autotags")
            run_autotags_tests
            ;;
        "server")
            run_server_tests
            ;;
        "tags")
            run_tags_tests
            ;;
        "stress")
            run_stress_tests
            ;;
        "all")
            run_all_tests
            ;;
        *)
            if [ -n "$TEST_PATTERN" ]; then
                run_custom_pattern "$TEST_PATTERN"
            else
                log_error "Unknown suite: $SUITE"
                print_usage
                exit 1
            fi
            ;;
    esac
    
    test_result=$?
    
    generate_coverage_report
    
    if [ $test_result -eq 0 ]; then
        echo ""
        log_success "🎉 All requested tests completed successfully!"
        show_test_summary
    else
        echo ""
        log_error "❌ Some tests failed. Check the output above for details."
        echo ""
        log_info "Debugging tips:"
        echo "  - Run with -v flag for verbose output"
        echo "  - Run specific test with -p 'TestName'"
        echo "  - Check /tmp for any leftover test files"
        echo "  - Ensure sufficient disk space for stress tests"
        exit 1
    fi
}

main 