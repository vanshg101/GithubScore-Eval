from pydantic import BaseModel, Field
from typing import Optional


class DeveloperFeatures(BaseModel):
    """Feature vector for ML prediction — maps to Go DeveloperMetrics."""

    total_commits: int = Field(..., ge=0, description="Total commits in last 12 months")
    total_prs: int = Field(..., ge=0, description="Total pull requests")
    merged_prs: int = Field(..., ge=0, description="Merged pull requests")
    total_issues_opened: int = Field(..., ge=0, description="Issues opened")
    total_issues_closed: int = Field(..., ge=0, description="Issues closed")
    review_comments: int = Field(..., ge=0, description="Code review comments")
    active_weeks: int = Field(..., ge=0, le=52, description="Active weeks out of 52")
    repos_contributed: int = Field(..., ge=0, description="Repos contributed to")
    total_stars: int = Field(..., ge=0, description="Stars across owned repos")
    total_forks: int = Field(..., ge=0, description="Forks across owned repos")
    avg_pr_lines_changed: float = Field(
        ..., ge=0, description="Avg lines changed per PR"
    )
    avg_issue_response_hours: float = Field(
        ..., ge=0, description="Avg hours to first response"
    )
    commit_trend_score: float = Field(
        ...,
        ge=0,
        le=1,
        description="Commit trend: 1.0=increasing, 0.5=stable, 0.2=decreasing",
    )
    language_count: int = Field(..., ge=0, description="Number of languages used")

    def to_feature_list(self) -> list[float]:
        """Convert to ordered float list for model input."""
        return [
            float(self.total_commits),
            float(self.total_prs),
            float(self.merged_prs),
            float(self.total_issues_opened),
            float(self.total_issues_closed),
            float(self.review_comments),
            float(self.active_weeks),
            float(self.repos_contributed),
            float(self.total_stars),
            float(self.total_forks),
            self.avg_pr_lines_changed,
            self.avg_issue_response_hours,
            self.commit_trend_score,
            float(self.language_count),
        ]


class PredictionResponse(BaseModel):
    """Response from /predict endpoint."""

    impact_score: float = Field(..., description="Predicted impact score (0-100)")
    confidence: Optional[float] = Field(None, description="Model confidence metric")


class ErrorResponse(BaseModel):
    """Standard error response."""

    error: str
    detail: Optional[str] = None
