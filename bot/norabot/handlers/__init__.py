from aiogram import Router

from . import start

def init_router() -> Router:
    router = Router()
    router.include_routers(start.router)
    return router