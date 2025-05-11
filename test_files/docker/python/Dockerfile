FROM python:3.9-slim

WORKDIR /app

RUN pip install flask

COPY . /app

EXPOSE 5000

CMD ["python", "server.py"]