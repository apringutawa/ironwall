from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    APP_NAME: str = "IronWall API"
    VERSION: str = "1.0.0"
    DASHBOARD_PORT: int = 8080
    
    TELEGRAM_BOT_TOKEN: str = ""
    TELEGRAM_CHAT_ID: str = ""
    
    CORS_ORIGINS: list = ["http://localhost:3000", "http://localhost:8080"]
    
    class Config:
        env_file = ".env"

settings = Settings()
