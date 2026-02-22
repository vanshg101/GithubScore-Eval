"""
Train a Ridge Regression model to predict Developer Impact Score.

Usage:
    python -m app.model.train

Outputs:
    - trained_models/impact_model.pkl
    - Evaluation metrics (R², MAE, RMSE)
"""

import os
import joblib
import numpy as np
import pandas as pd
from sklearn.linear_model import Ridge
from sklearn.model_selection import train_test_split
from sklearn.metrics import r2_score, mean_absolute_error, mean_squared_error
from sklearn.preprocessing import StandardScaler


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

DATA_PATH = os.path.join(
    os.path.dirname(__file__), "..", "..", "data", "training_data.csv"
)
MODEL_DIR = os.path.join(os.path.dirname(__file__), "..", "..", "trained_models")
MODEL_PATH = os.path.join(MODEL_DIR, "impact_model.pkl")
SCALER_PATH = os.path.join(MODEL_DIR, "scaler.pkl")


def train():
    print("Loading training data...")
    df = pd.read_csv(DATA_PATH)
    print(f"  Samples: {len(df)}")

    X = df[FEATURE_COLUMNS].values
    y = df[LABEL_COLUMN].values

    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, random_state=42
    )

    # Scale features
    scaler = StandardScaler()
    X_train_scaled = scaler.fit_transform(X_train)
    X_test_scaled = scaler.transform(X_test)

    # Train Ridge Regression
    model = Ridge(alpha=1.0)
    model.fit(X_train_scaled, y_train)

    # Evaluate
    y_pred = model.predict(X_test_scaled)
    r2 = r2_score(y_test, y_pred)
    mae = mean_absolute_error(y_test, y_pred)
    rmse = np.sqrt(mean_squared_error(y_test, y_pred))

    print("\n--- Model Evaluation ---")
    print(f"  R² Score:  {r2:.4f}")
    print(f"  MAE:       {mae:.4f}")
    print(f"  RMSE:      {rmse:.4f}")

    # Feature importances (coefficients)
    print("\n--- Feature Coefficients ---")
    for name, coef in sorted(
        zip(FEATURE_COLUMNS, model.coef_), key=lambda x: abs(x[1]), reverse=True
    ):
        print(f"  {name:30s} {coef:+.4f}")

    # Save model + scaler
    os.makedirs(MODEL_DIR, exist_ok=True)
    joblib.dump(model, MODEL_PATH)
    joblib.dump(scaler, SCALER_PATH)
    print(f"\nModel saved → {MODEL_PATH}")
    print(f"Scaler saved → {SCALER_PATH}")

    return r2, mae, rmse


if __name__ == "__main__":
    train()
