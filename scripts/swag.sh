#!/bin/bash

# Swagger documentation generation script
# Swagger 문서 생성 스크립트

set -e

# Colors for output
# 출력용 색상
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
# 기본값
MAIN_FILE="cmd/server/main.go"
OUTPUT_DIR="./docs"
PACKAGE_NAME="github.com/kyungseok-lee/fiber-gorm-starter"

# Function to print colored output
# 색상 출력 함수
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if swag is installed
# swag 설치 확인 함수
check_swag() {
    if ! command -v swag &> /dev/null; then
        print_error "swag command not found"
        echo ""
        echo "Please install swag with:"
        echo "  go install github.com/swaggo/swag/cmd/swag@latest"
        echo ""
        echo "Or use the Makefile:"
        echo "  make install-tools"
        exit 1
    fi
}

# Function to generate swagger docs
# swagger 문서 생성 함수
generate_docs() {
    print_info "Generating Swagger documentation..."
    
    # Create output directory if it doesn't exist
    # 출력 디렉터리가 없으면 생성
    mkdir -p "$OUTPUT_DIR"
    
    # Generate swagger documentation
    # swagger 문서 생성
    swag init \
        --generalInfo "$MAIN_FILE" \
        --dir ./ \
        --output "$OUTPUT_DIR" \
        --parseInternal \
        --parseDependency \
        --parseDepth 1
    
    if [ $? -eq 0 ]; then
        print_info "Swagger documentation generated successfully!"
        print_info "Files created:"
        ls -la "$OUTPUT_DIR"
        echo ""
        print_info "Access documentation at: http://localhost:8080/docs/index.html"
    else
        print_error "Failed to generate Swagger documentation"
        exit 1
    fi
}

# Function to validate generated docs
# 생성된 문서 검증 함수
validate_docs() {
    print_info "Validating generated documentation..."
    
    # Check if required files exist
    # 필수 파일 존재 확인
    required_files=("docs.go" "swagger.json" "swagger.yaml")
    for file in "${required_files[@]}"; do
        if [ ! -f "$OUTPUT_DIR/$file" ]; then
            print_error "Required file not found: $OUTPUT_DIR/$file"
            return 1
        fi
    done
    
    # Check if swagger.json is valid JSON
    # swagger.json이 유효한 JSON인지 확인
    if command -v jq &> /dev/null; then
        if ! jq . "$OUTPUT_DIR/swagger.json" > /dev/null 2>&1; then
            print_error "Generated swagger.json is not valid JSON"
            return 1
        fi
    else
        print_warning "jq not found, skipping JSON validation"
    fi
    
    print_info "Documentation validation passed!"
    return 0
}

# Function to show usage
# 사용법 표시 함수
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Generate Swagger documentation for the Fiber API"
    echo ""
    echo "Options:"
    echo "  -h, --help       Show this help message"
    echo "  -o, --output     Output directory (default: $OUTPUT_DIR)"
    echo "  -m, --main       Main file path (default: $MAIN_FILE)"
    echo "  -v, --validate   Validate generated documentation"
    echo "  --clean          Clean output directory before generation"
    echo ""
    echo "Examples:"
    echo "  $0                           # Generate with defaults"
    echo "  $0 --clean                   # Clean and generate"
    echo "  $0 -o ./api-docs             # Custom output directory"
    echo "  $0 --validate                # Generate and validate"
}

# Function to clean output directory
# 출력 디렉터리 정리 함수
clean_output() {
    print_info "Cleaning output directory: $OUTPUT_DIR"
    if [ -d "$OUTPUT_DIR" ]; then
        rm -rf "$OUTPUT_DIR"/*
        print_info "Output directory cleaned"
    fi
}

# Parse command line arguments
# 명령줄 인수 파싱
VALIDATE=false
CLEAN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -m|--main)
            MAIN_FILE="$2"
            shift 2
            ;;
        -v|--validate)
            VALIDATE=true
            shift
            ;;
        --clean)
            CLEAN=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Main execution
# 메인 실행
main() {
    print_info "Starting Swagger documentation generation..."
    print_info "Main file: $MAIN_FILE"
    print_info "Output directory: $OUTPUT_DIR"
    echo ""
    
    # Check if swag is installed
    # swag 설치 확인
    check_swag
    
    # Clean output directory if requested
    # 요청시 출력 디렉터리 정리
    if [ "$CLEAN" = true ]; then
        clean_output
    fi
    
    # Generate documentation
    # 문서 생성
    generate_docs
    
    # Validate documentation if requested
    # 요청시 문서 검증
    if [ "$VALIDATE" = true ]; then
        validate_docs
    fi
    
    echo ""
    print_info "Documentation generation completed!"
    print_info "You can now start the server and visit: http://localhost:8080/docs/index.html"
}

# Run main function
# 메인 함수 실행
main "$@"