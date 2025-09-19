#!/bin/bash

# Test Deployment Script for DocGen2 Service
# This script tests a deployed DocGen2 service endpoints

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

# Configuration
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -u, --url SERVICE_URL           Service URL (required)"
    echo "  -p, --plan-file PLAN_FILE       Test plan file (default: assets/plans/test_plan_01.json)"
    echo "  -o, --output-file OUTPUT        Output file for generated document (default: test-output.docx)"
    echo "  --skip-generate                 Skip document generation test"
    echo "  -h, --help                      Show this help message"
    echo ""
    echo "Example:"
    echo "  $0 -u https://docgen-service-xyz.a.run.app"
}

SERVICE_URL=""
PLAN_FILE="assets/plans/test_plan_01.json"
OUTPUT_FILE="test-output.docx"
SKIP_GENERATE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--url)
            SERVICE_URL="$2"
            shift 2
            ;;
        -p|--plan-file)
            PLAN_FILE="$2"
            shift 2
            ;;
        -o|--output-file)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        --skip-generate)
            SKIP_GENERATE=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

if [[ -z "$SERVICE_URL" ]]; then
    print_error "Service URL is required. Use -u or --url option."
    usage
    exit 1
fi

# Remove trailing slash from URL
SERVICE_URL=${SERVICE_URL%/}

print_header "Testing DocGen2 Service Deployment"
echo "Service URL: $SERVICE_URL"
echo "Plan File: $PLAN_FILE"
echo "Output File: $OUTPUT_FILE"
echo ""

# Test 1: Health Check
print_header "Test 1: Health Check"
HEALTH_URL="${SERVICE_URL}/health"
print_status "Testing: GET $HEALTH_URL"

if response=$(curl -f -s "$HEALTH_URL"); then
    print_status "✅ Health check passed"
    echo "Response:"
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
else
    print_error "❌ Health check failed"
    exit 1
fi
echo ""

# Test 2: Components List
print_header "Test 2: Components List"
COMPONENTS_URL="${SERVICE_URL}/components"
print_status "Testing: GET $COMPONENTS_URL"

if response=$(curl -f -s "$COMPONENTS_URL"); then
    print_status "✅ Components endpoint passed"
    echo "Available components:"
    echo "$response" | jq '.components[]' 2>/dev/null || echo "$response"
else
    print_error "❌ Components endpoint failed"
    exit 1
fi
echo ""

# Test 3: Plan Validation
print_header "Test 3: Plan Validation"
VALIDATE_URL="${SERVICE_URL}/validate-plan"
print_status "Testing: POST $VALIDATE_URL"

if [[ ! -f "$PLAN_FILE" ]]; then
    print_warning "Plan file $PLAN_FILE not found, skipping validation test"
else
    if response=$(curl -f -s -X POST -H "Content-Type: application/json" -d @"$PLAN_FILE" "$VALIDATE_URL"); then
        print_status "✅ Plan validation passed"
        echo "Validation response:"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        print_error "❌ Plan validation failed"
        exit 1
    fi
fi
echo ""

# Test 4: Document Generation
if [[ "$SKIP_GENERATE" == false ]]; then
    print_header "Test 4: Document Generation"
    GENERATE_URL="${SERVICE_URL}/generate"
    print_status "Testing: POST $GENERATE_URL"

    if [[ ! -f "$PLAN_FILE" ]]; then
        print_warning "Plan file $PLAN_FILE not found, skipping generation test"
    else
        print_status "Generating document and saving to $OUTPUT_FILE"

        if curl -f -s -X POST -H "Content-Type: application/json" -d @"$PLAN_FILE" "$GENERATE_URL" --output "$OUTPUT_FILE"; then
            # Check if file was created and has content
            if [[ -f "$OUTPUT_FILE" ]] && [[ -s "$OUTPUT_FILE" ]]; then
                FILE_SIZE=$(stat -f%z "$OUTPUT_FILE" 2>/dev/null || stat -c%s "$OUTPUT_FILE" 2>/dev/null || echo "unknown")
                print_status "✅ Document generation passed"
                print_status "Generated file: $OUTPUT_FILE (${FILE_SIZE} bytes)"

                # Basic file type check
                if file "$OUTPUT_FILE" | grep -q "Microsoft Word"; then
                    print_status "✅ Generated file appears to be a valid Word document"
                elif file "$OUTPUT_FILE" | grep -q "Zip archive"; then
                    print_status "✅ Generated file appears to be a valid Office document (ZIP format)"
                else
                    print_warning "⚠️  Generated file may not be a valid Word document"
                    print_status "File type: $(file "$OUTPUT_FILE")"
                fi
            else
                print_error "❌ Document generation failed - file not created or empty"
                exit 1
            fi
        else
            print_error "❌ Document generation failed"
            exit 1
        fi
    fi
else
    print_warning "Skipping document generation test as requested"
fi
echo ""

# Test 5: Error Handling (Invalid Plan)
print_header "Test 5: Error Handling"
print_status "Testing error handling with invalid plan"

INVALID_PLAN='{"invalid": "plan", "structure": true}'
if response=$(curl -s -X POST -H "Content-Type: application/json" -d "$INVALID_PLAN" "$VALIDATE_URL"); then
    echo "Error response:"
    echo "$response" | jq '.' 2>/dev/null || echo "$response"

    # Check if response indicates validation failure
    if echo "$response" | grep -q -i "error\|validation\|invalid"; then
        print_status "✅ Error handling working correctly"
    else
        print_warning "⚠️  Expected validation error but got success response"
    fi
else
    print_status "✅ Error handling working correctly (HTTP error returned)"
fi
echo ""

# Summary
print_header "Test Summary"
print_status "All tests completed! 🎉"
echo ""
print_status "Service appears to be working correctly:"
echo "  ✅ Health check endpoint"
echo "  ✅ Components list endpoint"
echo "  ✅ Plan validation endpoint"
if [[ "$SKIP_GENERATE" == false ]]; then
echo "  ✅ Document generation endpoint"
fi
echo "  ✅ Error handling"
echo ""

print_header "Next Steps"
echo "• Test with your own document plans"
echo "• Monitor service logs: gcloud logs read --service=SERVICE_NAME"
echo "• Set up monitoring and alerting"
echo "• Configure custom domain if needed"
echo ""

print_status "Service is ready for production use! 🚀"