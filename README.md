# GitHub Contribution Scoring & Developer Evaluation System

A full-stack system that evaluates GitHub developers by aggregating contribution metrics, computing weighted scores, and predicting impact using ML.

## Tech Stack

- **Backend:** Go (Gin)
- **ML Service:** Python (FastAPI, Scikit-Learn)
- **Frontend:** Next.js, TypeScript, Tailwind CSS
- **Database:** Google Cloud Firestore
- **Infrastructure:** Docker, GitHub Actions, GCP Cloud Run

## Project Structure

```
├── backend/          # Go REST API
├── ml-service/       # FastAPI prediction service
├── frontend/         # Next.js dashboard
└── docker-compose.yml
```

## Getting Started

### Prerequisites

- Go 1.22+
- Python 3.12+
- Node.js 20+
- Docker & Docker Compose

### Run Locally

```bash
# Backend
cd backend
cp .env.example .env
go run ./cmd/server/

# ML Service
cd ml-service
pip install -r requirements.txt
uvicorn app.main:app --port 8000

# Frontend
cd frontend
npm install
npm run dev
```

### Docker

```bash
docker-compose up --build
```

## Services

| Service    | Port | Endpoint     |
|------------|------|--------------|
| Backend    | 8080 | `/health`    |
| ML Service | 8000 | `/health`    |
| Frontend   | 3000 | `/`          |