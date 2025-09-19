#!/bin/bash

# GCP Project Setup Script for DocGen2 Service
# This script sets up the necessary GCP resources for deploying DocGen2 to Cloud Run

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
    echo -e "${BLUE}[SETUP]${NC} $1"
}

# Configuration
DEFAULT_PROJECT_ID=""
DEFAULT_REGION="us-central1"
DEFAULT_REPOSITORY="docgen-repo"
DEFAULT_SERVICE_NAME="docgen-service"

# Parse command line arguments
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -p, --project-id PROJECT_ID    GCP Project ID (required)"
    echo "  -r, --region REGION           GCP Region (default: us-central1)"
    echo "  --repository REPO             Artifact Registry repository name (default: docgen-repo)"
    echo "  --service-name NAME           Cloud Run service name (default: docgen-service)"
    echo "  -h, --help                    Show this help message"
    echo ""
    echo "Example:"
    echo "  $0 -p my-gcp-project -r us-west1"
}

PROJECT_ID=""
REGION="$DEFAULT_REGION"
REPOSITORY="$DEFAULT_REPOSITORY"
SERVICE_NAME="$DEFAULT_SERVICE_NAME"

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

print_header "Setting up GCP project for DocGen2 deployment"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION"
echo "Repository: $REPOSITORY"
echo "Service Name: $SERVICE_NAME"
echo ""

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    print_error "gcloud CLI is not installed. Please install it from https://cloud.google.com/sdk/"
    exit 1
fi

# Set the project
print_status "Setting GCP project to $PROJECT_ID"
gcloud config set project "$PROJECT_ID"

# Enable required APIs
print_header "Enabling required GCP APIs"
APIS=(
    "run.googleapis.com"
    "cloudbuild.googleapis.com"
    "artifactregistry.googleapis.com"
    "secretmanager.googleapis.com"
    "iam.googleapis.com"
    "cloudresourcemanager.googleapis.com"
)

for api in "${APIS[@]}"; do
    print_status "Enabling $api"
    gcloud services enable "$api"
done

# Create Artifact Registry repository
print_header "Creating Artifact Registry repository"
if gcloud artifacts repositories describe "$REPOSITORY" --location="$REGION" &>/dev/null; then
    print_warning "Repository $REPOSITORY already exists in $REGION"
else
    print_status "Creating Docker repository: $REPOSITORY"
    gcloud artifacts repositories create "$REPOSITORY" \
        --repository-format=docker \
        --location="$REGION" \
        --description="Docker repository for DocGen2 service"
fi

# Configure Docker authentication
print_header "Configuring Docker authentication"
print_status "Configuring Docker to authenticate with Artifact Registry"
gcloud auth configure-docker "${REGION}-docker.pkg.dev"

# Create Cloud Build trigger (optional, requires repository connection)
print_header "Cloud Build Setup"
print_warning "Cloud Build trigger setup requires connecting your source repository (GitHub, etc.)"
print_status "You can create a trigger manually in the Cloud Console or using:"
echo "  gcloud builds triggers create github \\"
echo "    --repo-name=your-repo-name \\"
echo "    --repo-owner=your-github-username \\"
echo "    --branch-pattern=^main$ \\"
echo "    --build-config=cloudbuild.yaml"

# Set up IAM permissions for Cloud Build
print_header "Setting up IAM permissions"
PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format="value(projectNumber)")
CLOUDBUILD_SA="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"

print_status "Granting Cloud Run Admin role to Cloud Build service account"
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$CLOUDBUILD_SA" \
    --role="roles/run.admin"

print_status "Granting IAM Service Account User role to Cloud Build service account"
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$CLOUDBUILD_SA" \
    --role="roles/iam.serviceAccountUser"

# Create environment-specific substitution file
print_header "Creating environment configuration"
cat > "cloudbuild-substitutions.yaml" << EOF
# Cloud Build substitution variables for environment: $PROJECT_ID
substitutions:
  _REGION: '$REGION'
  _REPOSITORY: '$REPOSITORY'
  _SERVICE_NAME: '$SERVICE_NAME'
  _MEMORY: '512Mi'
  _CPU: '1'
  _MIN_INSTANCES: '0'
  _MAX_INSTANCES: '10'
EOF

print_status "Created cloudbuild-substitutions.yaml with environment-specific configuration"

# Output summary
print_header "Setup Summary"
echo -e "${GREEN}✅ GCP Project configured: $PROJECT_ID${NC}"
echo -e "${GREEN}✅ Required APIs enabled${NC}"
echo -e "${GREEN}✅ Artifact Registry repository created: $REPOSITORY${NC}"
echo -e "${GREEN}✅ Docker authentication configured${NC}"
echo -e "${GREEN}✅ Cloud Build IAM permissions configured${NC}"
echo -e "${GREEN}✅ Environment configuration created${NC}"
echo ""

print_header "Next Steps"
echo "1. Connect your source repository to Cloud Build triggers"
echo "2. Run manual deployment: ./scripts/deploy-manual.sh -p $PROJECT_ID"
echo "3. Test the deployment: ./scripts/test-deployment.sh"
echo ""

print_status "Setup completed successfully!"

# Output important URLs and commands
print_header "Important Information"
echo "Artifact Registry URL: https://console.cloud.google.com/artifacts/docker/$PROJECT_ID/$REGION/$REPOSITORY"
echo "Cloud Build Console: https://console.cloud.google.com/cloud-build/dashboard?project=$PROJECT_ID"
echo "Cloud Run Console: https://console.cloud.google.com/run?project=$PROJECT_ID"
echo ""
echo "Image URL format: $REGION-docker.pkg.dev/$PROJECT_ID/$REPOSITORY/docgen-service:TAG"