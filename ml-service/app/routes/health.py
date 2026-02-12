from fastapi import APIRouter
from app.model.loader import ModelLoader

router = APIRouter()


@router.get("/health")
def health_check():
    model_loaded = False
    try:
        model_loaded = ModelLoader.get_instance().is_loaded
    except Exception:
        pass

    return {
        "status": "healthy",
        "service": "github-score-ml-service",
        "model_loaded": model_loaded,
    }
