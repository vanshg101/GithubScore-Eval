from fastapi import APIRouter, HTTPException
from app.schema.schemas import DeveloperFeatures, PredictionResponse
from app.model.loader import ModelLoader

router = APIRouter()


@router.post("/predict", response_model=PredictionResponse)
def predict_impact_score(features: DeveloperFeatures):
    """Predict Developer Impact Score from feature vector."""
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
