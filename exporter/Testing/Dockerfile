# Use the official Python image from Docker Hub
FROM python:3.11

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install the MQTT client library
RUN pip install "paho-mqtt<2.0.0"

# Expose the port the app runs on
EXPOSE 1883

# Run the Python script when the container launches
CMD ["python", "main.py"]