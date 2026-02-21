#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────
# GCP Setup Script for GitHub Score Eval
# Run this once to configure GCP project, Artifact Registry,
# Secret Manager, and Workload Identity Federation for GitHub Actions.
# ─────────────────────────────────────────────────────────────
set -euo pipefail

# ── Configuration ────────────────────────────────────────────
PROJECT_ID="${GCP_PROJECT_ID:?Set GCP_PROJECT_ID env var}"
REGION="${GCP_REGION:-us-central1}"
REPO_NAME="ghscore"
SA_NAME="ghscore-deployer"
GITHUB_REPO="${GITHUB_REPO:?Set GITHUB_REPO env var (owner/repo)}"
WIF_POOL="github-actions-pool"
WIF_PROVIDER="github-actions-provider"

echo "=== GCP Setup for github-score-eval ==="
echo "Project:  $PROJECT_ID"
echo "Region:   $REGION"
echo "Repo:     $GITHUB_REPO"
echo ""

# ── 1. Enable required APIs ─────────────────────────────────
echo "→ Enabling APIs..."
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com \
  cloudbuild.googleapis.com \
  iam.googleapis.com \
  iamcredentials.googleapis.com \
  --project="$PROJECT_ID"

# ── 2. Create Artifact Registry repository ───────────────────
echo "→ Creating Artifact Registry repo..."
gcloud artifacts repositories create "$REPO_NAME" \
  --repository-format=docker \
  --location="$REGION" \
  --description="GitHub Score Eval container images" \
  --project="$PROJECT_ID" 2>/dev/null || echo "  (already exists)"

# ── 3. Create Service Account ────────────────────────────────
echo "→ Creating service account..."
gcloud iam service-accounts create "$SA_NAME" \
  --display-name="GitHub Score Eval Deployer" \
  --project="$PROJECT_ID" 2>/dev/null || echo "  (already exists)"

SA_EMAIL="${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

# ── 4. Grant roles to Service Account ────────────────────────
echo "→ Granting IAM roles..."
for ROLE in \
  roles/run.admin \
  roles/artifactregistry.writer \
  roles/secretmanager.secretAccessor \
  roles/iam.serviceAccountUser; do
  gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SA_EMAIL" \
    --role="$ROLE" \
    --quiet > /dev/null
  echo "  Granted $ROLE"
done

# ── 5. Workload Identity Federation (keyless auth for GitHub Actions) ──
echo "→ Setting up Workload Identity Federation..."

# Create pool
gcloud iam workload-identity-pools create "$WIF_POOL" \
  --location="global" \
  --display-name="GitHub Actions Pool" \
  --project="$PROJECT_ID" 2>/dev/null || echo "  Pool already exists"

# Create provider
gcloud iam workload-identity-pools providers create-oidc "$WIF_PROVIDER" \
  --location="global" \
  --workload-identity-pool="$WIF_POOL" \
  --display-name="GitHub Actions Provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --project="$PROJECT_ID" 2>/dev/null || echo "  Provider already exists"

# Get the full provider name
WIF_PROVIDER_FULL=$(gcloud iam workload-identity-pools providers describe "$WIF_PROVIDER" \
  --location="global" \
  --workload-identity-pool="$WIF_POOL" \
  --project="$PROJECT_ID" \
  --format="value(name)")

# Allow GitHub repo to impersonate the SA
gcloud iam service-accounts add-iam-policy-binding "$SA_EMAIL" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/${WIF_PROVIDER_FULL}/attribute.repository/${GITHUB_REPO}" \
  --project="$PROJECT_ID" \
  --quiet > /dev/null

# ── 6. Create secrets in Secret Manager ──────────────────────
echo "→ Creating Secret Manager secrets..."

declare -a SECRETS=(
  "ghscore-firestore-key"
  "ghscore-jwt-secret"
  "ghscore-gh-oauth-client-id"
  "ghscore-gh-oauth-client-secret"
)

for SECRET in "${SECRETS[@]}"; do
  gcloud secrets create "$SECRET" \
    --replication-policy="automatic" \
    --project="$PROJECT_ID" 2>/dev/null || echo "  $SECRET already exists"
done

# ── 7. Print GitHub Secrets to configure ─────────────────────
echo ""
echo "════════════════════════════════════════════════════════"
echo " DONE! Now set these GitHub repository secrets:"
echo "════════════════════════════════════════════════════════"
echo ""
echo "  GCP_PROJECT_ID       = $PROJECT_ID"
echo "  GCP_WIF_PROVIDER     = $WIF_PROVIDER_FULL"
echo "  GCP_SERVICE_ACCOUNT  = $SA_EMAIL"
echo "  GH_OAUTH_CLIENT_ID   = <your GitHub OAuth App client ID>"
echo "  GH_OAUTH_CLIENT_SECRET = <your GitHub OAuth App client secret>"
echo "  JWT_SECRET           = <random 32+ char string>"
echo "  ML_SERVICE_URL       = https://ghscore-ml-xxxxx-uc.a.run.app"
echo "  BACKEND_URL          = https://ghscore-backend-xxxxx-uc.a.run.app"
echo "  FRONTEND_URL         = https://ghscore-frontend-xxxxx-uc.a.run.app"
echo ""
echo "Then add the Firestore service account key:"
echo "  gcloud secrets versions add ghscore-firestore-key \\"
echo "    --data-file=backend/github-score-eval-eb16d2dd3723.json \\"
echo "    --project=$PROJECT_ID"
echo ""
echo "  gcloud secrets versions add ghscore-jwt-secret \\"
echo "    --data-file=- <<< '<your-jwt-secret>' \\"
echo "    --project=$PROJECT_ID"
echo ""
