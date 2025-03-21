from aiogram.filters.callback_data import CallbackData


class ConnectData(CallbackData, prefix="connect"):
    telegram_id: int
    user_id: str