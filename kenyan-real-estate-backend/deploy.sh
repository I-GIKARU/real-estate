#!/bin/bash

# Deployment script for Google Cloud Run
# Make sure to set up the required environment variables and authenticate with gcloud

set -e

# Configuration
PROJECT_ID=${PROJECT_ID:-"your-project-id"}
SERVICE_NAME="kenyan-real-estate-backend"
REGION=${REGION:-"us-central1"}
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME"

echo "🚀 Deploying Kenyan Real Estate Backend to Google Cloud Run"

# Check if gcloud is authenticated
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "."; then
    echo "❌ Please authenticate with gcloud first: gcloud auth login"
    exit 1
fi

# Set the project
gcloud config set project $PROJECT_ID

# Enable required APIs
echo "📡 Enabling required APIs..."
gcloud services enable containerregistry.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com

# Build and push the Docker image
echo "🏗️  Building Docker image..."
gcloud builds submit --tag $IMAGE_NAME

# Deploy to Cloud Run
echo "🚀 Deploying to Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --memory 512Mi \
  --cpu 1 \
  --concurrency 100 \
  --timeout 300s \
  --max-instances 10 \
  --min-instances 0 \
  --port 8080 \
  --set-env-vars="APP_ENV=production,SERVER_HOST=0.0.0.0,SERVER_PORT=8080"

# Get the service URL
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --platform managed --region $REGION --format 'value(status.url)')

echo "✅ Deployment complete!"
echo "🌐 Service URL: $SERVICE_URL"
echo "🔍 Health check: $SERVICE_URL/health"

# Test the health endpoint
echo "🏥 Testing health endpoint..."
if curl -f "$SERVICE_URL/health" > /dev/null 2>&1; then
    echo "✅ Health check passed!"
else
    echo "❌ Health check failed. Please check the logs:"
    echo "📝 View logs: gcloud logs read --project=$PROJECT_ID --resource.type=cloud_run_revision"
fi
