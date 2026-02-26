#!/bin/bash

echo "Setting up AI Control Layer SDK..."

# Create virtual environment
python3 -m venv venv
source venv/bin/activate

# Install dependencies
pip install --upgrade pip
pip install -e ".[dev]"

# Copy environment file
if [ ! -f .env ]; then
    cp .env.example .env
    echo "Created .env file - please configure it"
fi

# Start Docker services
docker-compose up -d

# Wait for services
echo "Waiting for services to be ready..."
sleep 5

# Run migrations
alembic upgrade head

echo "Setup complete! Run 'source venv/bin/activate' to activate the environment"
