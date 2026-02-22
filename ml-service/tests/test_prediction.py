"""Unit tests for the prediction logic and schema validation."""

import pytest
from app.schema.schemas import DeveloperFeatures, PredictionResponse


def _sample_features(**overrides) -> dict:
    """Return a valid feature dict with optional overrides."""
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


class TestDeveloperFeatures:
    def test_valid_features(self):
        f = DeveloperFeatures(**_sample_features())
        assert f.total_commits == 200
        assert f.language_count == 5

    def test_to_feature_list_length(self):
        f = DeveloperFeatures(**_sample_features())
        fl = f.to_feature_list()
        assert len(fl) == 14
        assert all(isinstance(v, float) for v in fl)

    def test_negative_commits_rejected(self):
        with pytest.raises(Exception):
            DeveloperFeatures(**_sample_features(total_commits=-1))

    def test_active_weeks_over_52_rejected(self):
        with pytest.raises(Exception):
            DeveloperFeatures(**_sample_features(active_weeks=53))

    def test_commit_trend_out_of_range_rejected(self):
        with pytest.raises(Exception):
            DeveloperFeatures(**_sample_features(commit_trend_score=1.5))


class TestPredictionResponse:
    def test_valid_response(self):
        r = PredictionResponse(impact_score=75.5, confidence=0.9)
        assert r.impact_score == 75.5

    def test_no_confidence(self):
        r = PredictionResponse(impact_score=50.0)
        assert r.confidence is None


class TestModelPrediction:
    """Tests that require a trained model — skip if not available."""

    @pytest.fixture(autouse=True)
    def _load_model(self):
        try:
            from app.model.loader import ModelLoader

            self.model = ModelLoader.get_instance()
        except FileNotFoundError:
            pytest.skip("Trained model not available")

    def test_prediction_in_range(self):
        features = DeveloperFeatures(**_sample_features())
        score = self.model.predict(features.to_feature_list())
        assert 0 <= score <= 100

    def test_zero_activity_low_score(self):
        features = DeveloperFeatures(
            **_sample_features(
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
            )
        )
        score = self.model.predict(features.to_feature_list())
        assert score < 35  # very low activity should score low

    def test_high_activity_high_score(self):
        features = DeveloperFeatures(
            **_sample_features(
                total_commits=500,
                total_prs=100,
                merged_prs=90,
                total_issues_opened=50,
                total_issues_closed=45,
                review_comments=100,
                active_weeks=48,
                repos_contributed=20,
                total_stars=200,
                total_forks=50,
                avg_pr_lines_changed=200,
                avg_issue_response_hours=6,
                commit_trend_score=1.0,
                language_count=10,
            )
        )
        score = self.model.predict(features.to_feature_list())
        assert score > 60  # high activity should score high
