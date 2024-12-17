from queries import summarize_chat_multishot_query

from fastapi import FastAPI, HTTPException, Body
import uvicorn
from pydantic import BaseModel, Field

import requests
from dotenv import load_dotenv
import os
from typing import List, Dict


app = FastAPI()

load_dotenv()


OLLAMA_NGROK_URL = os.getenv("OLLAMA_NGROK_URL")
OLLAMA_API_KEY = os.getenv("OLLAMA_API_KEY")
PORT = int(os.getenv("PORT"))


class SummarizeChatResponse(BaseModel):
    response: str = Field(..., example="The user asked for information about services.")


def process_chat_input(chat_conversation: List[Dict]) -> str:
    """
    Process chat input by concatenating messages from the conversation.

    Args:
        chat_conversation (List[Dict[str, str]]): A list of dictionaries containing chat messages.

    Returns:
        str: A concatenated string of chat messages.
    """
    processed_chat = ""
    for message in chat_conversation:
        processed_chat += f"{message['sentBy']}: {message['content']}\n"
    return processed_chat


@app.post("/sumarize_chat", response_model=SummarizeChatResponse)
def sumarize_chat(
    input: Dict[str, List[Dict[str, str]]] = Body(
        ...,
        example={
            "messages": [
                {"sentBy": "user1", "content": "Hello!"},
                {"sentBy": "user2", "content": "Hi there! How can I help you today?"},
                {
                    "sentBy": "user1",
                    "content": "I need some information about your services.",
                },
            ]
        },
    )
):
    """
    Summarize a chat conversation.

    Args:
        input (list[dict]): A list of dictionaries containing chat messages.

    Returns:
        dict: The summarized chat response from the API.

    Raises:
        HTTPException: If the API request fails.
    """
    messages = input["messages"]
    processed_chat = process_chat_input(messages)
    query = summarize_chat_multishot_query(processed_chat)
    response = requests.post(
        OLLAMA_NGROK_URL,
        json={"query": f"{query}"},
        headers={"Authorization": f"Bearer {OLLAMA_API_KEY}"},
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


if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=PORT, reload=True)
