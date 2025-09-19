# DocGen2 Deployment Guide

## Table of Contents
- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Manual Deployment](#manual-deployment)
- [Automated CI/CD Deployment](#automated-cicd-deployment)
- [Configuration](#configuration)
- [Testing](#testing)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)
- [Security Considerations](#security-considerations)

## Overview

This guide provides complete instructions for deploying DocGen2 to Google Cloud Run. DocGen2 is a containerized microservice that generates Microsoft Word documents from JSON plans using a component-based approach.

### Architecture

- **Container**: Multi-stage Docker build with distroless base image
- **Platform**: Google Cloud Run (serverless, auto-scaling)
- **Storage**: Artifact Registry for container images
- **CI/CD**: Cloud Build with GitHub integration
- **Monitoring**: Cloud Run metrics + health endpoints

## Prerequisites

### Required Tools
- [Google Cloud SDK (gcloud)](https://cloud.google.com/sdk/docs/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Git](https://git-scm.com/downloads)
- [curl](https://curl.se/) and [jq](https://stedolan.jq.io/) (for testing)

### GCP Requirements
- GCP Project with billing enabled
- Sufficient IAM permissions:
  - Cloud Run Admin
  - Artifact Registry Admin
  - Cloud Build Editor
  - Service Account User

### Verify Prerequisites
```bash
# Check tool installations
gcloud version
docker --version
git --version
curl --version
jq --version

# Authenticate with GCP
gcloud auth login
gcloud auth configure-docker
```

## Quick Start

For first-time deployment, follow these steps:

### 1. Clone and Setup
```bash
git clone <your-repo-url>
cd docgen2
```

### 2. GCP Project Setup
```bash
# Run the automated setup script
./scripts/setup-gcp.sh -p YOUR_PROJECT_ID

# Or with custom region
./scripts/setup-gcp.sh -p YOUR_PROJECT_ID -r us-west1
```

### 3. Deploy the Service
```bash
# Deploy with default settings
./scripts/deploy-manual.sh -p YOUR_PROJECT_ID --allow-unauthenticated

# Or with custom configuration
./scripts/deploy-manual.sh -p YOUR_PROJECT_ID \
  --region us-west1 \
  --memory 1Gi \
  --max-instances 20 \
  --allow-unauthenticated
```

### 4. Test the Deployment
```bash
# Get the service URL from the deployment output, then test
./scripts/test-deployment.sh -u https://docgen-service-xyz.a.run.app
```

## Manual Deployment

### Step 1: GCP Project Setup

The setup script configures your GCP project with all necessary resources:

```bash
./scripts/setup-gcp.sh -p YOUR_PROJECT_ID [OPTIONS]
```

**Options:**
- `-p, --project-id`: GCP Project ID (required)
- `-r, --region`: GCP Region (default: us-central1)
- `--repository`: Artifact Registry repository name (default: docgen-repo)
- `--service-name`: Cloud Run service name (default: docgen-service)

**What it does:**
- ✅ Enables required GCP APIs (Cloud Run, Cloud Build, Artifact Registry, etc.)
- ✅ Creates Artifact Registry repository for Docker images
- ✅ Configures Docker authentication
- ✅ Sets up IAM permissions for Cloud Build
- ✅ Creates environment-specific configuration file

### Step 2: Build and Deploy

The deployment script handles building, pushing, and deploying:

```bash
./scripts/deploy-manual.sh -p YOUR_PROJECT_ID [OPTIONS]
```

**Key Options:**
- `-p, --project-id`: GCP Project ID (required)
- `--allow-unauthenticated`: Allow public access (recommended for initial testing)
- `--memory`: Memory allocation (default: 512Mi)
- `--cpu`: CPU allocation (default: 1)
- `--max-instances`: Maximum auto-scaling instances (default: 10)
- `--skip-tests`: Skip running tests before deployment
- `--skip-build`: Use existing image (for rapid deployment)

**Example Commands:**
```bash
# Basic deployment
./scripts/deploy-manual.sh -p my-project --allow-unauthenticated

# Production deployment with higher resources
./scripts/deploy-manual.sh -p my-project \
  --memory 1Gi \
  --cpu 2 \
  --max-instances 50 \
  --min-instances 1

# Quick redeploy without rebuild
./scripts/deploy-manual.sh -p my-project --skip-build --skip-tests
```

### Step 3: Verify Deployment

Test all endpoints to ensure proper deployment:

```bash
# Comprehensive testing
./scripts/test-deployment.sh -u https://your-service-url.a.run.app

# Quick health check only
curl https://your-service-url.a.run.app/health
```

## Automated CI/CD Deployment

### Cloud Build Configuration

The project includes a complete Cloud Build configuration (`cloudbuild.yaml`) that:

1. **Runs Tests**: Executes `go test ./...` to ensure code quality
2. **Builds Image**: Creates optimized Docker image with multi-stage build
3. **Pushes to Registry**: Stores image in Artifact Registry with SHA and latest tags
4. **Deploys to Cloud Run**: Updates the service with new image

### Setting Up CI/CD

#### 1. Connect Source Repository

**Option A: GitHub (Recommended)**
```bash
gcloud builds triggers create github \
  --repo-name=docgen2 \
  --repo-owner=YOUR_USERNAME \
  --branch-pattern=^main$ \
  --build-config=cloudbuild.yaml \
  --substitutions=_REGION=us-central1,_REPOSITORY=docgen-repo
```

**Option B: Cloud Source Repositories**
```bash
# Mirror your repository to Cloud Source Repositories first
gcloud source repos create docgen2
git remote add google https://source.developers.google.com/p/YOUR_PROJECT/r/docgen2

# Create trigger
gcloud builds triggers create cloud-source-repositories \
  --repo=docgen2 \
  --branch-pattern=^main$ \
  --build-config=cloudbuild.yaml
```

#### 2. Configure Trigger Settings

In the Cloud Console, configure your trigger with:

**Basic Settings:**
- Event: Push to branch
- Branch: `^main$`
- Configuration: Cloud Build configuration file
- Location: `cloudbuild.yaml`

**Advanced Settings (Substitution Variables):**
```yaml
_REGION: us-central1
_REPOSITORY: docgen-repo
_SERVICE_NAME: docgen-service
_MEMORY: 512Mi
_CPU: 1
_MIN_INSTANCES: 0
_MAX_INSTANCES: 10
```

### CI/CD Best Practices

1. **Branch Protection**: Only deploy from `main` branch
2. **Testing**: All tests must pass before deployment
3. **Rollback**: Keep previous image tags for quick rollback
4. **Notifications**: Configure Slack/email notifications for build status
5. **Security**: Use least-privilege IAM roles

## Configuration

### Environment Variables

DocGen2 is configured entirely through environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DOCGEN_SHELL_PATH` | `./assets/shell/template_shell.docx` | Path to shell document template |
| `DOCGEN_COMPONENTS_DIR` | `./assets/components/` | Directory containing component XML files |
| `DOCGEN_SCHEMA_PATH` | `./assets/schemas/rules.cue` | Path to CUE validation schema |

### Cloud Run Configuration

**Resource Allocation:**
- **Memory**: Start with 512Mi, increase if needed for large documents
- **CPU**: 1 CPU sufficient for most workloads
- **Concurrency**: Default (80) works well for document generation

**Scaling:**
- **Min Instances**: 0 for cost optimization, 1+ for consistent latency
- **Max Instances**: Set based on expected load and quota limits
- **Timeout**: Default (300s) sufficient for document generation

**Security:**
- **Authentication**: Enable IAM authentication for production
- **HTTPS**: Automatically provided by Cloud Run
- **VPC**: Consider VPC connector for private resources

### Custom Configuration Example

```bash
gcloud run services update docgen-service \
  --region=us-central1 \
  --memory=1Gi \
  --cpu=2 \
  --min-instances=1 \
  --max-instances=20 \
  --set-env-vars="CUSTOM_VAR=value" \
  --no-allow-unauthenticated
```

## Testing

### Automated Testing

The test script validates all service endpoints:

```bash
./scripts/test-deployment.sh -u https://your-service-url.a.run.app
```

**Test Coverage:**
- ✅ Health check endpoint (`GET /health`)
- ✅ Components list endpoint (`GET /components`)
- ✅ Plan validation endpoint (`POST /validate-plan`)
- ✅ Document generation endpoint (`POST /generate`)
- ✅ Error handling (invalid requests)

### Manual Testing

**Health Check:**
```bash
curl https://your-service-url.a.run.app/health
```

**List Available Components:**
```bash
curl https://your-service-url.a.run.app/components
```

**Validate a Document Plan:**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d @assets/plans/test_plan_01.json \
  https://your-service-url.a.run.app/validate-plan
```

**Generate a Document:**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d @assets/plans/test_plan_01.json \
  --output generated.docx \
  https://your-service-url.a.run.app/generate
```

### Load Testing

For production readiness, consider load testing:

```bash
# Using Apache Bench (ab)
ab -n 100 -c 10 -T application/json -p assets/plans/test_plan_01.json \
  https://your-service-url.a.run.app/generate

# Using Artillery
artillery quick --count 10 --num 20 https://your-service-url.a.run.app/health
```

## Monitoring

### Cloud Run Metrics

Monitor these key metrics in Cloud Console:

**Performance:**
- Request count and latency
- Error rate and status codes
- Memory and CPU utilization
- Instance count (active/idle)

**Reliability:**
- Cold start frequency
- Request timeout rate
- Container startup time

### Application Logs

View logs for debugging and monitoring:

```bash
# Recent logs
gcloud logs read --service=docgen-service --limit=50

# Follow logs in real-time
gcloud logs tail --service=docgen-service

# Filter by severity
gcloud logs read --service=docgen-service --filter="severity>=ERROR"
```

### Health Monitoring

Set up uptime checks and alerting:

1. **Cloud Monitoring Uptime Checks**
2. **Alerting Policies** for high error rates or latency
3. **Dashboard** with key metrics

### Custom Metrics

The service exposes health status at `/health`:

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "components": ["DocumentTitle", "TestBlock", ...],
  "version": "1.0.0"
}
```

## Troubleshooting

### Common Issues

**1. Service Not Starting**
```bash
# Check deployment status
gcloud run services describe docgen-service --region=us-central1

# Check logs for startup errors
gcloud logs read --service=docgen-service --filter="severity>=ERROR" --limit=20
```

**2. Authentication Errors**
```bash
# Check IAM permissions
gcloud projects get-iam-policy YOUR_PROJECT_ID

# Test with authentication
curl -H "Authorization: Bearer $(gcloud auth print-access-token)" \
  https://your-service-url.a.run.app/health
```

**3. Build Failures**
```bash
# Check Cloud Build history
gcloud builds list --limit=10

# View specific build logs
gcloud builds log BUILD_ID
```

**4. Resource Limits**
```bash
# Check resource usage
gcloud run services describe docgen-service --region=us-central1 \
  --format="table(spec.template.spec.containers[0].resources)"

# Update resources if needed
gcloud run services update docgen-service --memory=1Gi --cpu=2
```

### Debug Commands

```bash
# Service information
gcloud run services describe docgen-service --region=us-central1

# Recent deployments
gcloud run revisions list --service=docgen-service --region=us-central1

# Traffic allocation
gcloud run services describe docgen-service --region=us-central1 \
  --format="table(status.traffic[0].revisionName,status.traffic[0].percent)"

# Container image details
gcloud run services describe docgen-service --region=us-central1 \
  --format="value(spec.template.spec.containers[0].image)"
```

### Log Analysis

**Common log patterns to monitor:**

```bash
# Failed requests
gcloud logs read --filter='resource.type="cloud_run_revision" AND httpRequest.status>=400'

# Slow requests (>5 seconds)
gcloud logs read --filter='resource.type="cloud_run_revision" AND httpRequest.latency>"5s"'

# Memory or CPU issues
gcloud logs read --filter='resource.type="cloud_run_revision" AND ("memory" OR "CPU")'
```

## Security Considerations

### Authentication and Authorization

**Production Deployment:**
```bash
# Deploy with authentication required
./scripts/deploy-manual.sh -p YOUR_PROJECT_ID  # (no --allow-unauthenticated)

# Test with authentication
curl -H "Authorization: Bearer $(gcloud auth print-access-token)" \
  https://your-service-url.a.run.app/health
```

**Service-to-Service Authentication:**
```bash
# Create service account for clients
gcloud iam service-accounts create docgen-client

# Grant Cloud Run Invoker role
gcloud run services add-iam-policy-binding docgen-service \
  --member="serviceAccount:docgen-client@YOUR_PROJECT.iam.gserviceaccount.com" \
  --role="roles/run.invoker"
```

### Network Security

**VPC Connector (Optional):**
```bash
# Create VPC connector for private access
gcloud compute networks vpc-access connectors create docgen-connector \
  --region=us-central1 \
  --subnet=your-subnet-name

# Deploy with VPC connector
gcloud run deploy docgen-service \
  --vpc-connector=docgen-connector \
  --vpc-egress=private-ranges-only
```

### Container Security

**Security Features:**
- ✅ **Distroless base image** (minimal attack surface)
- ✅ **Non-root user** (security best practice)
- ✅ **Read-only filesystem** (except /tmp)
- ✅ **No shell access** (distroless = no shell)

### Data Security

**Sensitive Data Handling:**
- Document plans may contain sensitive information
- Use HTTPS for all communication (automatic with Cloud Run)
- Consider implementing encryption for document storage
- Audit access logs regularly

### Security Monitoring

```bash
# Monitor failed authentication attempts
gcloud logs read --filter='resource.type="cloud_run_revision" AND httpRequest.status=401'

# Check for unusual traffic patterns
gcloud logs read --filter='resource.type="cloud_run_revision"' \
  --format='table(timestamp,httpRequest.remoteIp,httpRequest.requestMethod,httpRequest.requestUrl,httpRequest.status)'
```

---

## Quick Reference

### Essential Commands

```bash
# Deploy to production
./scripts/deploy-manual.sh -p YOUR_PROJECT_ID

# Test deployment
./scripts/test-deployment.sh -u https://your-service-url.a.run.app

# View logs
gcloud logs read --service=docgen-service --limit=20

# Update service
gcloud run services update docgen-service --memory=1Gi

# Get service URL
gcloud run services describe docgen-service --format="value(status.url)"
```

### Support

- **Documentation**: `/docs/` directory in the repository
- **Issues**: GitHub Issues for bug reports and feature requests
- **Logs**: Use `gcloud logs read` for troubleshooting
- **Monitoring**: Cloud Console > Cloud Run > docgen-service

---

This completes the comprehensive deployment guide for DocGen2 on Google Cloud Run. The service is now ready for production use with proper monitoring, security, and scalability configurations.