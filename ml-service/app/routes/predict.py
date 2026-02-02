from fastapi import APIRouter

router = APIRouter()


@router.post("/predict")
def predict_impact_score():
    return {
        "message": "Prediction endpoint",
        "impact_score": None,
    }
