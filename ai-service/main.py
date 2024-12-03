from fastapi import FastAPI, HTTPException, Body
from pydantic import BaseModel, Field
from dotenv import load_dotenv
from queries import summarize_chat_multishot_query
import requests
import os

app = FastAPI()

load_dotenv()

ollama_negrok_url = os.getenv("OLLAMA_NEGROK_URL")
ollama_api_key = os.getenv("OLLAMA_API_KEY")

class SummarizeChatResponse(BaseModel):
    response: str = Field(..., example="The user asked for information about services.")

def process_chat_input(chat_conversation: list[dict]) -> str:
    """
    Process chat input by concatenating messages from the conversation.

    Args:
        chat_conversation (list[dict]): A list of dictionaries containing chat messages.

    Returns:
        str: A concatenated string of chat messages.
    """
    processed_chat = ""
    for message in chat_conversation:
        processed_chat += f"{message['sender']}: {message['message']}\n"
    return processed_chat

@app.post("/sumarize_chat", response_model=SummarizeChatResponse)
def sumarize_chat(input: list[dict] = Body(
    ...,
    example=[
        {"sender": "user1", "message": "Hello!"},
        {"sender": "user2", "message": "Hi there! How can I help you today?"},
        {"sender": "user1", "message": "I need some information about your services."}
    ]
)):
    """
    Summarize a chat conversation.

    Args:
        input (list[dict]): A list of dictionaries containing chat messages.

    Returns:
        dict: The summarized chat response from the API.

    Raises:
        HTTPException: If the API request fails.
    """
    processed_chat = process_chat_input(input)
    query = summarize_chat_multishot_query(processed_chat)
    response = requests.post(
        ollama_negrok_url,
        json={"query": f"{query}"},
        headers={"Authorization": f"Bearer {ollama_api_key}"},
    )
    if response.status_code == 200:
        return response.json()
    else:
        raise HTTPException(status_code=response.status_code, detail=response.text)

@app.post("/detect_emotions")
def detect_emotions():
    """
    Detect emotions in a chat conversation.

    This endpoint is not yet implemented.
    """
    pass