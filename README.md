# GitHub Contribution Scoring & Developer Evaluation System

A full-stack system that evaluates GitHub developers by aggregating contribution metrics, computing weighted scores across 12 indicators, and predicting impact using a Ridge Regression ML model. Features include developer comparison, organization evaluation, leaderboard rankings, and interactive data visualizations.

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)
![Python](https://img.shields.io/badge/Python-3.13-3776AB?logo=python)
![Next.js](https://img.shields.io/badge/Next.js-16-000000?logo=next.js)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)
![License](https://img.shields.io/badge/License-MIT-green)

---

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Next.js 16    │────▶│   Go Backend    │────▶│  FastAPI ML     │
│   Frontend      │     │   (Gin)         │     │  Service        │
│   Port 3000     │     │   Port 8080     │     │  Port 8000      │
└─────────────────┘     └────────┬────────┘     └─────────────────┘
                                 │
                    ┌────────────┼────────────┐
                    ▼            ▼            ▼
              ┌──────────┐ ┌──────────┐ ┌──────────┐
              │Firestore │ │GitHub API│ │Cron Jobs │
              │ Database │ │  v3/v4   │ │Scheduler │
              └──────────┘ └──────────┘ └──────────┘
```

## Tech Stack

| Layer          | Technology                                        |
|----------------|---------------------------------------------------|
| **Backend**    | Go 1.24, Gin, JWT auth, GitHub OAuth              |
| **ML Service** | Python 3.13, FastAPI, Scikit-Learn (Ridge Regression) |
| **Frontend**   | Next.js 16, React 19, TypeScript, Tailwind CSS v4, Recharts |
| **Database**   | Google Cloud Firestore                            |
| **Infra**      | Docker, Docker Compose, GitHub Actions CI/CD      |
| **Deployment** | GCP Cloud Run, Artifact Registry                  |

## Project Structure

```
GithubScoreEval/
├── backend/                    # Go REST API
│   ├── cmd/server/main.go      # Entry point
│   └── internal/
│       ├── auth/               # GitHub OAuth + JWT
│       ├── config/             # Environment config
│       ├── cron/               # Scheduled jobs (refresh, ranking)
│       ├── firestore/          # Firestore client
│       ├── github/             # GitHub API client (commits, PRs, issues, etc.)
│       ├── handler/            # HTTP handlers (auth, developer, score, ranking)
│       ├── middleware/         # Auth, CORS, logger
│       ├── mlclient/           # ML service HTTP client
│       ├── model/              # Data models (Developer, Score, Ranking, User)
│       ├── repository/         # Firestore repositories
│       ├── router/             # Route definitions
│       ├── scoring/            # 12-indicator scoring engine
│       └── service/            # Business logic
├── ml-service/                 # FastAPI ML prediction service
│   ├── app/
│   │   ├── main.py             # FastAPI app
│   │   ├── model/              # Model loader + training script
│   │   ├── routes/             # /health, /predict endpoints
│   │   ├── schema/             # Pydantic request/response models
│   │   └── utils/              # Training data generation
│   ├── trained_models/         # Serialized model + scaler (.pkl)
│   └── data/                   # Training dataset
├── frontend/                   # Next.js dashboard
│   └── src/
│       ├── app/                # Pages (dashboard, search, compare, leaderboard, org)
│       ├── components/         # UI components + 7 chart components (Recharts)
│       ├── context/            # Auth + Theme providers
│       └── lib/                # Typed API client
├── .github/workflows/          # CI (lint+test) + CD (Cloud Run deploy)
├── scripts/                    # GCP setup script
├── docker-compose.yml          # Full stack orchestration
└── README.md
```

---

## Getting Started

### Prerequisites

- **Go** 1.24+
- **Python** 3.13+
- **Node.js** 20+
- **Docker** & Docker Compose (for containerized mode)
- **GitHub OAuth App** - create one at [github.com/settings/developers](https://github.com/settings/developers)
- **GCP Project** with Firestore enabled + service account JSON

### 1. Clone & Configure

```bash
git clone https://github.com/vanshg101/GithubScoreEval.git
cd GithubScoreEval
```

Create `backend/.env` from the example:

```bash
cp backend/.env.example backend/.env
```

Fill in the required values:

```env
PORT=8080
ENVIRONMENT=development

# GitHub OAuth App credentials
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
GITHUB_REDIRECT_URL=http://localhost:8080/auth/github/callback

# JWT signing secret (random 32+ char string)
JWT_SECRET=your_secret_key_here

# Google Cloud / Firestore
GCP_PROJECT_ID=your-gcp-project-id
FIRESTORE_CREDENTIALS=./path-to-service-account.json

# ML service URL
ML_SERVICE_URL=http://localhost:8000

# Frontend URL (for OAuth redirect)
FRONTEND_URL=http://localhost:3000
```

Create `frontend/.env.local`:

```bash
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > frontend/.env.local
```

### 2a. Run with Docker (Recommended)

```bash
docker compose up --build -d

# Check status
docker compose ps

# View logs
docker compose logs -f

# Stop
docker compose down
```

### 2b. Run Locally (Individual Services)

**Backend:**
```bash
cd backend
go run ./cmd/server/
```

**ML Service:**
```bash
cd ml-service
python -m venv .venv
source .venv/bin/activate        # Windows: .venv\Scripts\activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8000
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

### 3. Access the App

| Service      | URL                          | Health Check              |
|--------------|------------------------------|---------------------------|
| Frontend     | http://localhost:3000         | -                         |
| Backend API  | http://localhost:8080         | http://localhost:8080/health |
| ML Service   | http://localhost:8000         | http://localhost:8000/health |

---

## Scoring Algorithm

The system evaluates developers across **12 weighted indicators** (total weight = 1.00):

| Indicator                 | Weight | Max  | Description                                |
|---------------------------|--------|------|--------------------------------------------|
| Total Commits             | 0.15   | 500  | Raw commit count across repos              |
| PR Merge Rate             | 0.12   | 1.0  | `merged_prs / total_prs`                   |
| Contribution Consistency  | 0.12   | 1.0  | `active_weeks / 52`                        |
| Code Review Participation | 0.10   | 100  | Number of review comments                  |
| Issue Resolution Ratio    | 0.08   | 1.0  | `closed_issues / opened_issues`            |
| Repo Diversity            | 0.08   | 20   | Number of repos contributed to             |
| Stars Earned              | 0.08   | 500  | Total stars across repositories            |
| Avg PR Size               | 0.07   | 1.0  | Closeness to ideal 200 LOC changes         |
| Fork Impact               | 0.05   | 100  | Total forks across repositories            |
| Issue Response Time       | 0.05   | 1.0  | Inverse of avg response hours (max 168h)   |
| Language Diversity        | 0.05   | 10   | Number of programming languages used       |
| Commit Trend              | 0.05   | 1.0  | Increasing=1.0, Stable=0.5, Decreasing=0.2 |

Each indicator is normalized to `[0, 1]`, multiplied by its weight, and summed to produce a **weighted score (0–100)**. An **ML Impact Score** from the prediction service supplements the evaluation.

## ML Model

- **Algorithm:** Ridge Regression (α=1.0) with StandardScaler
- **Features:** 14 developer metrics matching the Go backend's `DeveloperMetrics` model
- **Output:** Impact score (0–100) with confidence measure
- **Training:** 80/20 split, evaluated with R², MAE, RMSE

Retrain the model:
```bash
cd ml-service
python -m app.model.train
```

---

## API Reference

### Authentication

| Method | Path                        | Auth | Description              |
|--------|-----------------------------|------|--------------------------|
| GET    | `/auth/github/login`        | No   | Redirect to GitHub OAuth |
| GET    | `/auth/github/callback`     | No   | OAuth callback → JWT     |
| POST   | `/auth/logout`              | No   | Clear auth cookie        |
| GET    | `/auth/me`                  | Yes  | Get current user profile |

### Developers

| Method | Path                                   | Auth | Description                    |
|--------|----------------------------------------|------|--------------------------------|
| GET    | `/api/developers`                      | Yes  | List all fetched developers    |
| GET    | `/api/developers/:username`            | Yes  | Get developer by username      |
| POST   | `/api/developers/:username/fetch`      | Yes  | Fetch developer from GitHub    |

### Scoring

| Method | Path                                   | Auth | Description                    |
|--------|----------------------------------------|------|--------------------------------|
| POST   | `/api/developers/:username/score`      | Yes  | Compute score for developer    |
| GET    | `/api/developers/:username/score`      | Yes  | Get stored score               |

### Rankings & Comparison

| Method | Path                        | Auth | Description                          |
|--------|-----------------------------|------|--------------------------------------|
| GET    | `/api/rankings`             | Yes  | Paginated leaderboard (`?page=1&page_size=10`) |
| POST   | `/api/compare`              | Yes  | Compare 2–10 developers             |
| GET    | `/api/orgs/:org/evaluate`   | Yes  | Evaluate all members of an org       |

### ML Service

| Method | Path       | Description                           |
|--------|------------|---------------------------------------|
| GET    | `/health`  | Service health + model loaded status  |
| POST   | `/predict` | Predict impact score from 14 features |

---

## Frontend Pages

| Route                    | Description                                      |
|--------------------------|--------------------------------------------------|
| `/`                      | Landing page                                     |
| `/login`                 | GitHub OAuth login                               |
| `/dashboard`             | Main dashboard with search + quick links         |
| `/search`                | Developer search                                 |
| `/developer/[username]`  | Full profile: stats, radar chart, commit trends, PR donut, language pie, heatmap, ML score |
| `/leaderboard`           | Sortable ranking table with pagination           |
| `/compare`               | Compare 2–10 developers with overlay radar chart |
| `/org`                   | Organization member evaluation                   |

### Charts (Recharts)

- **Score Radar Chart** - 12-indicator spider/radar chart
- **Commit Trend Chart** - Monthly commit line chart
- **PR Merge Rate Chart** - Donut chart (merged vs unmerged)
- **Language Distribution** - Pie chart of programming languages
- **Contribution Heatmap** - 52-week CSS grid heatmap
- **ML Impact Score** - Color-coded progress bar card
- **Comparison Overlay** - Multi-user radar chart overlay

---

## CI/CD

### CI Pipeline (on PR / push to main)
1. **Go Backend** - `go vet` → `go build` → `go test`
2. **Python ML** - `ruff` lint → `mypy` check → `pytest`
3. **Next.js Frontend** - `npm run lint` → `npm run build`
4. **Docker** - Build all 3 images with Buildx caching

### CD Pipeline (on push to main)
1. Authenticate to GCP via Workload Identity Federation
2. Build & push images to Artifact Registry
3. Deploy ML Service → Backend → Frontend to Cloud Run
4. Smoke test all `/health` endpoints

---

## Docker

Multi-stage builds produce optimized images:

| Image      | Base                           | Final Size |
|------------|--------------------------------|------------|
| Backend    | `golang:1.25-alpine` → `alpine:3.20` | ~56 MB |
| ML Service | `python:3.13-slim` (2-stage)   | ~694 MB    |
| Frontend   | `node:20-alpine` (3-stage, standalone) | ~294 MB |

All containers run as non-root users with health checks.

---

## Environment Variables

| Variable                  | Service  | Description                            |
|---------------------------|----------|----------------------------------------|
| `PORT`                    | Backend  | Server port (default: 8080)            |
| `GITHUB_CLIENT_ID`        | Backend  | GitHub OAuth App client ID             |
| `GITHUB_CLIENT_SECRET`    | Backend  | GitHub OAuth App client secret         |
| `GITHUB_REDIRECT_URL`     | Backend  | OAuth callback URL                     |
| `JWT_SECRET`              | Backend  | JWT signing secret                     |
| `GCP_PROJECT_ID`          | Backend  | GCP project for Firestore              |
| `FIRESTORE_CREDENTIALS`   | Backend  | Path to service account JSON           |
| `ML_SERVICE_URL`          | Backend  | ML service URL                         |
| `FRONTEND_URL`            | Backend  | Frontend URL for CORS + OAuth redirect |
| `NEXT_PUBLIC_API_URL`     | Frontend | Backend API URL                        |

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.