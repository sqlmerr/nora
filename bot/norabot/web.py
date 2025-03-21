from uuid import UUID

from fastapi import FastAPI, APIRouter
from pydantic import BaseModel
from aiogram import Bot, Dispatcher

router = APIRouter()

class ConnectData(BaseModel):
    user_id: UUID
    telegram_id: int

@router.post("/connect")
async def connect_telegram(data: ConnectData, bot: Bot):
    await bot.send_message(data.telegram_id, "Are you sure to connect")

def create_app(bot: Bot, dp: Dispatcher) -> FastAPI:
    app = FastAPI(docs_url=None, openapi_url=None)
    app.include_router(router)
    app.dependency_overrides[Bot] = lambda: bot
    app.dependency_overrides[Dispatcher] = lambda: dp
    
    return app