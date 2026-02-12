"""
Generate synthetic training data for the Developer Impact Score model.

The impact score label is derived from a weighted formula that mimics
how a senior engineering manager would evaluate developer effectiveness.
Random noise is added to simulate real-world variance.
"""

import csv
import os
import random
import math

FEATURE_COLUMNS = [
    "total_commits",
    "total_prs",
    "merged_prs",
    "total_issues_opened",
    "total_issues_closed",
    "review_comments",
    "active_weeks",
    "repos_contributed",
    "total_stars",
    "total_forks",
    "avg_pr_lines_changed",
    "avg_issue_response_hours",
    "commit_trend_score",
    "language_count",
]

LABEL_COLUMN = "impact_score"

NUM_SAMPLES = 2000
SEED = 42


def _clamp(v: float, lo: float = 0.0, hi: float = 100.0) -> float:
    return max(lo, min(hi, v))


def _generate_row(rng: random.Random) -> dict:
    """Generate one synthetic developer sample with a realistic impact label."""

    # --- features ---
    total_commits = rng.randint(0, 800)
    total_prs = rng.randint(0, 200)
    merged_prs = rng.randint(0, total_prs) if total_prs > 0 else 0
    total_issues_opened = rng.randint(0, 150)
    total_issues_closed = rng.randint(0, total_issues_opened) if total_issues_opened > 0 else 0
    review_comments = rng.randint(0, 200)
    active_weeks = rng.randint(0, 52)
    repos_contributed = rng.randint(0, 40)
    total_stars = int(rng.paretovariate(1.5) * 5)  # long-tail
    total_forks = int(total_stars * rng.uniform(0.05, 0.4)) if total_stars > 0 else 0
    avg_pr_lines = round(rng.uniform(0, 600), 1)
    avg_response_hours = round(rng.uniform(0, 168), 1)
    commit_trend = rng.choice([0.2, 0.5, 1.0])
    language_count = rng.randint(0, 15)

    # --- label: weighted impact formula ---
    commit_norm = min(total_commits / 500, 1.0)
    pr_merge_rate = merged_prs / total_prs if total_prs > 0 else 0
    issue_ratio = total_issues_closed / total_issues_opened if total_issues_opened > 0 else 0
    review_norm = min(review_comments / 100, 1.0)
    consistency = active_weeks / 52
    repo_norm = min(repos_contributed / 20, 1.0)
    star_norm = min(math.log1p(total_stars) / math.log1p(500), 1.0)
    fork_norm = min(math.log1p(total_forks) / math.log1p(100), 1.0)

    # PR size: ideal is ~200 lines
    pr_size_score = max(0, 1 - abs(avg_pr_lines - 200) / 200)

    # Response time: lower is better
    response_score = max(0, 1 - avg_response_hours / 168)

    lang_norm = min(language_count / 10, 1.0)

    impact = (
        commit_norm * 15
        + pr_merge_rate * 12
        + issue_ratio * 8
        + review_norm * 10
        + consistency * 12
        + repo_norm * 8
        + star_norm * 8
        + fork_norm * 5
        + pr_size_score * 7
        + response_score * 5
        + commit_trend * 5
        + lang_norm * 5
    )

    # Add noise (±5 points)
    noise = rng.gauss(0, 2.5)
    impact = _clamp(impact + noise)

    return {
        "total_commits": total_commits,
        "total_prs": total_prs,
        "merged_prs": merged_prs,
        "total_issues_opened": total_issues_opened,
        "total_issues_closed": total_issues_closed,
        "review_comments": review_comments,
        "active_weeks": active_weeks,
        "repos_contributed": repos_contributed,
        "total_stars": total_stars,
        "total_forks": total_forks,
        "avg_pr_lines_changed": avg_pr_lines,
        "avg_issue_response_hours": avg_response_hours,
        "commit_trend_score": commit_trend,
        "language_count": language_count,
        LABEL_COLUMN: round(impact, 2),
    }


def generate_dataset(output_path: str, n: int = NUM_SAMPLES, seed: int = SEED):
    rng = random.Random(seed)
    fieldnames = FEATURE_COLUMNS + [LABEL_COLUMN]

    os.makedirs(os.path.dirname(output_path), exist_ok=True)

    with open(output_path, "w", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        writer.writeheader()
        for _ in range(n):
            writer.writerow(_generate_row(rng))

    print(f"Generated {n} samples → {output_path}")


if __name__ == "__main__":
    data_dir = os.path.join(os.path.dirname(__file__), "..", "..", "data")
    generate_dataset(os.path.join(data_dir, "training_data.csv"))
