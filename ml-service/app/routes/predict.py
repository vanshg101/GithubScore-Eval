"""
Prediction route — accepts a developer feature vector and returns
an ML-predicted impact score via a pre-trained Ridge Regression model.

The model is loaded once (singleton) from disk on first request.
"""

from fastapi import APIRouter, HTTPException
from app.schema.schemas import DeveloperFeatures, PredictionResponse
from app.model.loader import ModelLoader

router = APIRouter()


@router.post("/predict", response_model=PredictionResponse)
def predict_impact_score(features: DeveloperFeatures):
    """Predict Developer Impact Score from a 14-feature vector.

    The feature vector mirrors the Go backend's DeveloperMetrics struct.
    See DeveloperFeatures schema for field details and validation rules.
    """
    try:
        model = ModelLoader.get_instance()
    except FileNotFoundError as e:
        raise HTTPException(status_code=503, detail=str(e))

    feature_list = features.to_feature_list()
    impact_score = model.predict(feature_list)

    return PredictionResponse(
        impact_score=round(impact_score, 2),
        confidence=None,
    )
