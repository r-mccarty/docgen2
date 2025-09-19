#!/bin/bash

# Manual Deployment Script for DocGen2 Service
# This script manually builds and deploys DocGen2 to Google Cloud Run

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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
    echo -e "${BLUE}[DEPLOY]${NC} $1"
}

# Configuration defaults
DEFAULT_REGION="us-central1"
DEFAULT_REPOSITORY="docgen-repo"
DEFAULT_SERVICE_NAME="docgen-service"
DEFAULT_MEMORY="512Mi"
DEFAULT_CPU="1"
DEFAULT_MIN_INSTANCES="0"
DEFAULT_MAX_INSTANCES="10"

# Parse command line arguments
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -p, --project-id PROJECT_ID     GCP Project ID (required)"
    echo "  -r, --region REGION             GCP Region (default: us-central1)"
    echo "  --repository REPO               Artifact Registry repository name (default: docgen-repo)"
    echo "  --service-name NAME             Cloud Run service name (default: docgen-service)"
    echo "  --memory MEMORY                 Cloud Run memory allocation (default: 512Mi)"
    echo "  --cpu CPU                       Cloud Run CPU allocation (default: 1)"
    echo "  --min-instances MIN             Cloud Run minimum instances (default: 0)"
    echo "  --max-instances MAX             Cloud Run maximum instances (default: 10)"
    echo "  --tag TAG                       Docker image tag (default: current timestamp)"
    echo "  --skip-build                    Skip Docker build (use existing image)"
    echo "  --skip-tests                    Skip running tests before deployment"
    echo "  --allow-unauthenticated         Allow unauthenticated access to the service"
    echo "  -h, --help                      Show this help message"
    echo ""
    echo "Example:"
    echo "  $0 -p my-gcp-project --allow-unauthenticated"
}

PROJECT_ID=""
REGION="$DEFAULT_REGION"
REPOSITORY="$DEFAULT_REPOSITORY"
SERVICE_NAME="$DEFAULT_SERVICE_NAME"
MEMORY="$DEFAULT_MEMORY"
CPU="$DEFAULT_CPU"
MIN_INSTANCES="$DEFAULT_MIN_INSTANCES"
MAX_INSTANCES="$DEFAULT_MAX_INSTANCES"
TAG=""
SKIP_BUILD=false
SKIP_TESTS=false
ALLOW_UNAUTHENTICATED=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--project-id)
            PROJECT_ID="$2"
            shift 2
            ;;
        -r|--region)
            REGION="$2"
            shift 2
            ;;
        --repository)
            REPOSITORY="$2"
            shift 2
            ;;
        --service-name)
            SERVICE_NAME="$2"
            shift 2
            ;;
        --memory)
            MEMORY="$2"
            shift 2
            ;;
        --cpu)
            CPU="$2"
            shift 2
            ;;
        --min-instances)
            MIN_INSTANCES="$2"
            shift 2
            ;;
        --max-instances)
            MAX_INSTANCES="$2"
            shift 2
            ;;
        --tag)
            TAG="$2"
            shift 2
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        --allow-unauthenticated)
            ALLOW_UNAUTHENTICATED=true
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

# Validate required parameters
if [[ -z "$PROJECT_ID" ]]; then
    print_error "Project ID is required. Use -p or --project-id option."
    usage
    exit 1
fi

# Set default tag if not provided
if [[ -z "$TAG" ]]; then
    TAG=$(date +%Y%m%d-%H%M%S)
fi

# Validate we're in the right directory
if [[ ! -f "go.mod" ]] || [[ ! -f "Dockerfile" ]]; then
    print_error "This script must be run from the DocGen2 project root directory"
    exit 1
fi

print_header "Deploying DocGen2 Service to Google Cloud Run"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION"
echo "Repository: $REPOSITORY"
echo "Service Name: $SERVICE_NAME"
echo "Tag: $TAG"
echo "Memory: $MEMORY"
echo "CPU: $CPU"
echo "Min Instances: $MIN_INSTANCES"
echo "Max Instances: $MAX_INSTANCES"
echo ""

# Check if gcloud is installed and authenticated
if ! command -v gcloud &> /dev/null; then
    print_error "gcloud CLI is not installed. Please install it from https://cloud.google.com/sdk/"
    exit 1
fi

# Set the project
print_status "Setting GCP project to $PROJECT_ID"
gcloud config set project "$PROJECT_ID"

# Run tests unless skipped
if [[ "$SKIP_TESTS" == false ]]; then
    print_header "Running tests"
    print_status "Executing Go tests..."
    go test ./... -v
    print_status "All tests passed ✅"
else
    print_warning "Skipping tests as requested"
fi

# Build and push Docker image unless skipped
IMAGE_URL="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/docgen-service:${TAG}"
LATEST_URL="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/docgen-service:latest"

if [[ "$SKIP_BUILD" == false ]]; then
    print_header "Building Docker image"
    print_status "Building image: $IMAGE_URL"

    docker build \
        -t "$IMAGE_URL" \
        -t "$LATEST_URL" \
        .

    print_status "Docker image built successfully ✅"

    print_header "Pushing image to Artifact Registry"
    print_status "Pushing $IMAGE_URL"
    docker push "$IMAGE_URL"

    print_status "Pushing $LATEST_URL"
    docker push "$LATEST_URL"

    print_status "Image pushed successfully ✅"
else
    print_warning "Skipping Docker build as requested"
    print_status "Using existing image: $IMAGE_URL"
fi

# Deploy to Cloud Run
print_header "Deploying to Cloud Run"

# Build gcloud run deploy command
DEPLOY_CMD=(
    "gcloud" "run" "deploy" "$SERVICE_NAME"
    "--image=$IMAGE_URL"
    "--platform=managed"
    "--region=$REGION"
    "--port=8080"
    "--memory=$MEMORY"
    "--cpu=$CPU"
    "--min-instances=$MIN_INSTANCES"
    "--max-instances=$MAX_INSTANCES"
    "--set-env-vars=DOCGEN_SHELL_PATH=./assets/shell/template_shell.docx,DOCGEN_COMPONENTS_DIR=./assets/components/,DOCGEN_SCHEMA_PATH=./assets/schemas/rules.cue"
)

# Add authentication setting
if [[ "$ALLOW_UNAUTHENTICATED" == true ]]; then
    DEPLOY_CMD+=("--allow-unauthenticated")
    print_status "Service will allow unauthenticated access"
else
    DEPLOY_CMD+=("--no-allow-unauthenticated")
    print_warning "Service will require authentication"
fi

print_status "Executing deployment command..."
"${DEPLOY_CMD[@]}"

# Get service URL
SERVICE_URL=$(gcloud run services describe "$SERVICE_NAME" --region="$REGION" --format="value(status.url)")

print_header "Deployment completed successfully! 🎉"
echo ""
echo -e "${GREEN}✅ Service deployed successfully${NC}"
echo -e "${GREEN}✅ Service URL: $SERVICE_URL${NC}"
echo ""

# Test the deployment
print_header "Testing deployment"
print_status "Testing health endpoint..."

HEALTH_URL="${SERVICE_URL}/health"
print_status "GET $HEALTH_URL"

if curl -f -s "$HEALTH_URL" > /dev/null; then
    print_status "Health check passed ✅"
    echo ""
    print_status "Fetching health details:"
    curl -s "$HEALTH_URL" | jq '.' || curl -s "$HEALTH_URL"
else
    print_warning "Health check failed - service may still be starting up"
    print_status "You can check the logs with: gcloud logs read --service=$SERVICE_NAME"
fi

echo ""
print_header "Deployment Summary"
echo "Service Name: $SERVICE_NAME"
echo "Service URL: $SERVICE_URL"
echo "Image: $IMAGE_URL"
echo "Region: $REGION"
echo ""

print_header "Available Endpoints"
echo "Health Check: $SERVICE_URL/health"
echo "Components List: $SERVICE_URL/components"
echo "Validate Plan: $SERVICE_URL/validate-plan (POST)"
echo "Generate Document: $SERVICE_URL/generate (POST)"
echo ""

print_header "Useful Commands"
echo "View logs: gcloud logs read --service=$SERVICE_NAME --region=$REGION"
echo "Service info: gcloud run services describe $SERVICE_NAME --region=$REGION"
echo "Update service: gcloud run services update $SERVICE_NAME --region=$REGION"
echo "Delete service: gcloud run services delete $SERVICE_NAME --region=$REGION"

print_status "Deployment completed! 🚀"