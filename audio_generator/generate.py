import os
from datetime import datetime
from openai import OpenAI
open_key = os.environ.get('OPENAI_API_KEY')
client = OpenAI(api_key=open_key)
answer = """Hello, I am Global Predictor AI. Curious to know what insights I have for you today?
"""
string_to_replace = ['-',':',' ','.']
now = str(datetime.now())
for character in string_to_replace: 
    now = now.replace(character,'_') 

audio = client.audio.speech.create(
                model="tts-1",
                voice="shimmer",
                input=answer,
                )

with open(f"../backend/audio-idle/{now}.mp3", 'wb') as f:
    f.write(audio.content)
