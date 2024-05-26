FROM python:3.8

WORKDIR /app

# Клонируем репозиторий с Python-проектом
RUN git clone https://github.com/blackHATred/kinoskop-profanity-filter .

# Устанавливаем зависимости
RUN pip3 install --no-cache-dir -r requirements.txt
# RUN python3 -m spacy download ru_core_news_md
RUN python3 -m spacy download en

# Запускаем main.py при старте контейнера
CMD ["python3", "-u", "./main.py", "--addr=0.0.0.0:8050"]
