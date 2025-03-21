import uuid
import httpx

from aiogram import Router
from aiogram.filters import CommandObject, CommandStart
from aiogram.types import Message, CallbackQuery
from aiogram.utils.keyboard import InlineKeyboardBuilder

from norabot.config import settings
from norabot.utils.calldata import ConnectData


router = Router()

@router.message(CommandStart(deep_link=True))
async def start_deep_link(message: Message, command: CommandObject) -> None:
    deep_link = command.args
    if deep_link.startswith("c") and len(deep_link) > 1:
        deep_link = deep_link.removeprefix("c")
        try:
            user_id = uuid.UUID(deep_link)
        except ValueError:
            return
        b = InlineKeyboardBuilder().button(text="Connect", callback_data=ConnectData(telegram_id=message.from_user.id, user_id=str(user_id)).pack())
        await message.answer(f"Are you sure to connect with user {user_id}?", reply_markup=b.as_markup())
    else:
        await start_cmd(message)
        

@router.callback_query(ConnectData.filter())
async def connect_callback(callback: CallbackQuery, callback_data: ConnectData, client: httpx.AsyncClient) -> None:
    if callback.from_user.id != callback_data.telegram_id:
        await callback.answer(cache_time=30)
        return
    
    
    r = await client.post("/auth/telegram/connect", json={"telegram_id": callback_data.telegram_id, "user_id": callback_data.user_id}, headers={"secret-token": settings.secret_token})
    json = r.json()
    if json.get("ok"):
        await callback.answer("Connected")
        await callback.message.edit_reply_markup(reply_markup=None)
    else:
        await callback.answer("Error")

@router.message(CommandStart())
async def start_cmd(message: Message) -> None:
    await message.answer("hi")
    
