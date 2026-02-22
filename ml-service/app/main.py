"""
FastAPI application entry point for the GitHub Score ML microservice.

This service exposes a /predict endpoint that accepts developer metrics
and returns an ML-predicted impact score using a Ridge Regression model.
Auto-generated Swagger docs are available at /docs.
"""

from fastapi import FastAPI
from app.routes import health, predict

app = FastAPI(
    title="GitHub Score ML Service",
    description="ML microservice for predicting Developer Impact Score",
    version="1.0.0",
)

app.include_router(health.router, tags=["Health"])
app.include_router(predict.router, tags=["Prediction"])
