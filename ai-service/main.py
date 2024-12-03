from fastapi import FastApi
from dotenv import load_dotenv
import os



app = FastApi()

load_dotenv()

ollama_negrok_url = os.getenv("OLLAMA_NEGROK_URL")
ollama_api_key = os.getenv("OLLAMA_API_KEY")


@app.post("/sumarize_chat")
def sumarize_chat():
    pass


@app.post("/detect_emotions")
def detect_emotions():
    pass
