"""
Model loader — loads the trained Ridge model and scaler at startup.
"""

import os
import joblib
import numpy as np
from typing import Optional

MODEL_PATH = os.path.join(os.path.dirname(__file__), "..", "..", "trained_models", "impact_model.pkl")
SCALER_PATH = os.path.join(os.path.dirname(__file__), "..", "..", "trained_models", "scaler.pkl")


class ModelLoader:
    """Singleton-style model loader for the impact prediction model."""

    _instance: Optional["ModelLoader"] = None
    _model = None
    _scaler = None
    _loaded = False

    @classmethod
    def get_instance(cls) -> "ModelLoader":
        if cls._instance is None:
            cls._instance = cls()
            cls._instance._load()
        return cls._instance

    def _load(self):
        if not os.path.exists(MODEL_PATH):
            raise FileNotFoundError(
                f"Model not found at {MODEL_PATH}. Run training first: "
                "python -m app.model.train"
            )
        if not os.path.exists(SCALER_PATH):
            raise FileNotFoundError(
                f"Scaler not found at {SCALER_PATH}. Run training first: "
                "python -m app.model.train"
            )

        self._model = joblib.load(MODEL_PATH)
        self._scaler = joblib.load(SCALER_PATH)
        self._loaded = True

    def predict(self, features: list[float]) -> float:
        """Predict impact score from a feature vector."""
        if not self._loaded:
            raise RuntimeError("Model not loaded")

        X = np.array(features).reshape(1, -1)
        X_scaled = self._scaler.transform(X)
        prediction = self._model.predict(X_scaled)[0]

        # Clamp to [0, 100]
        return float(max(0.0, min(100.0, prediction)))

    @property
    def is_loaded(self) -> bool:
        return self._loaded
