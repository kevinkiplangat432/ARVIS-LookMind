# Contributing to AI Control Layer SDK

Thank you for your interest in contributing to the AI Control Layer SDK!

## Development Setup

1. Clone the repository
2. Install dependencies: `pip install -e ".[dev]"`
3. Copy `.env.example` to `.env` and configure
4. Start services: `docker-compose up -d`
5. Run migrations: `alembic upgrade head`

## Code Standards

- Follow PEP 8 style guidelines
- Use type hints for all functions
- Write docstrings for public APIs
- Maintain test coverage above 80%
- Use async/await for I/O operations

## Testing

```bash
pytest tests/
pytest --cov=src tests/
```

## Pull Request Process

1. Create a feature branch from `main`
2. Make your changes with clear commit messages
3. Add tests for new functionality
4. Update documentation as needed
5. Ensure all tests pass
6. Submit PR with detailed description

## Code Review

All submissions require review. We use GitHub pull requests for this purpose.

## Questions?

Open an issue or contact: kiplangatkevin335@gmail.com
