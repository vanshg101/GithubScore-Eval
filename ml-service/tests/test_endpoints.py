"""Integration tests for the /predict and /health endpoints."""

import pytest
from fastapi.testclient import TestClient
from app.main import app


client = TestClient(app)


def _sample_payload(**overrides) -> dict:
    defaults = {
        "total_commits": 200,
        "total_prs": 30,
        "merged_prs": 25,
        "total_issues_opened": 15,
        "total_issues_closed": 10,
        "review_comments": 40,
        "active_weeks": 30,
        "repos_contributed": 10,
        "total_stars": 50,
        "total_forks": 10,
        "avg_pr_lines_changed": 180.0,
        "avg_issue_response_hours": 24.0,
        "commit_trend_score": 0.5,
        "language_count": 5,
    }
    defaults.update(overrides)
    return defaults


class TestHealthEndpoint:
    def test_health_returns_200(self):
        resp = client.get("/health")
        assert resp.status_code == 200
        data = resp.json()
        assert data["status"] == "healthy"
        assert "model_loaded" in data


class TestPredictEndpoint:
    def test_predict_returns_score(self):
        resp = client.post("/predict", json=_sample_payload())
        if resp.status_code == 503:
            pytest.skip("Model not trained yet")
        assert resp.status_code == 200
        data = resp.json()
        assert "impact_score" in data
        assert 0 <= data["impact_score"] <= 100

    def test_predict_missing_field(self):
        payload = _sample_payload()
        del payload["total_commits"]
        resp = client.post("/predict", json=payload)
        assert resp.status_code == 422  # validation error

    def test_predict_negative_value(self):
        resp = client.post("/predict", json=_sample_payload(total_commits=-5))
        assert resp.status_code == 422

    def test_predict_high_activity(self):
        resp = client.post(
            "/predict",
            json=_sample_payload(
                total_commits=500,
                total_prs=100,
                merged_prs=90,
                active_weeks=48,
                repos_contributed=20,
                total_stars=200,
                commit_trend_score=1.0,
                language_count=10,
            ),
        )
        if resp.status_code == 503:
            pytest.skip("Model not trained yet")
        assert resp.status_code == 200
        assert resp.json()["impact_score"] > 50

    def test_predict_zero_activity(self):
        resp = client.post(
            "/predict",
            json=_sample_payload(
                total_commits=0,
                total_prs=0,
                merged_prs=0,
                total_issues_opened=0,
                total_issues_closed=0,
                review_comments=0,
                active_weeks=0,
                repos_contributed=0,
                total_stars=0,
                total_forks=0,
                avg_pr_lines_changed=0,
                avg_issue_response_hours=0,
                commit_trend_score=0.2,
                language_count=0,
            ),
        )
        if resp.status_code == 503:
            pytest.skip("Model not trained yet")
        assert resp.status_code == 200
        assert resp.json()["impact_score"] < 35
