FROM python:3.12

ENV PYTHONUNBUFFERED 1

# Set the Current Working Directory inside the container
WORKDIR /app

# Install dependencies
COPY ai-service/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the source code
COPY ai-service .

# Expose port 8084 to the outside world
EXPOSE 8084

# Command to run the service
CMD ["python", "main.py"]