import asyncio
import os
import logging

from aiogram import Bot, Dispatcher
from httpx import AsyncClient
from norabot.config import settings
from norabot.handlers import init_router

logging.basicConfig(level=logging.INFO)
log = logging.getLogger(__name__)

async def main():
    bot = Bot(settings.bot_token)
    client = AsyncClient(base_url=settings.api_url)
    dp = Dispatcher(client=client)
    dp.startup.register(on_startup)
    dp.include_router(init_router())
    
    log.info("Starting bot")
    await dp.start_polling(bot)

async def on_startup():
    log.info("Bot started!")
    

if __name__ == "__main__":
    asyncio.run(main())