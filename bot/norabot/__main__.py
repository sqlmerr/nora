import asyncio
import os
import logging

from aiogram import Bot, Dispatcher
from norabot.config import settings

logging.basicConfig(level=logging.INFO)
log = logging.getLogger(__name__)

async def main():    
    bot = Bot(settings.bot_token)
    dp = Dispatcher()
    dp.startup.register(on_startup)
    
    log.info("Starting bot")
    await dp.start_polling(bot)

async def on_startup():
    log.info("Bot started!")
    

if __name__ == "__main__":
    asyncio.run(main())