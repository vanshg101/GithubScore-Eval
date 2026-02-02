from fastapi import FastAPI
from app.routes import health, predict

app = FastAPI(
    title="GitHub Score ML Service",
    description="ML microservice for predicting Developer Impact Score",
    version="1.0.0",
)

app.include_router(health.router, tags=["Health"])
app.include_router(predict.router, tags=["Prediction"])
