from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    bot_token: str
    api_url: str
    secret_token: str
    
    model_config = SettingsConfigDict(env_file=".env", extra="ignore")
    

settings = Settings()
