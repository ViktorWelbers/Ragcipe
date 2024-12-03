# REWE Recipe Scraper and Vector Search 🍽️🔍

## Overview

This Golang project provides a comprehensive solution for scraping, embedding, and querying recipes from REWE.com using advanced vector search technologies.

## Features

- **Web Scraping**: Automated scraping of approximately 20,000 recipes from REWE.com
- **AI Embedding**: Utilizes Ollama for generating recipe embeddings
- **Vector Database**: Stores and indexes recipes using Qdrant vector database
- **CLI Integration**: Seamless command-line interface for recipe querying and management

## Tech Stack

- **Language**: Go (Golang)
- **Scraping**: Custom web scraping implementation
- **Embedding**: Ollama 
- **Vector Database**: Qdrant
- **CLI**: Custom CLI tool with recipe search functionality

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/rewe-recipe-scraper.git

# Navigate to project directory
cd rewe-recipe-scraper

# Install dependencies
go mod download

# Build the project
go build
```

## Usage

### Scraping Recipes
```bash
# Scrape recipes from REWE
./recipe-scraper scrape
```

### Querying Recipes
```bash
# Search for recipes (returns 3-5 matching recipes)
./recipe-scraper search "vegetarian pasta"
```

## Project Architecture

1. **Scraper**: Extracts recipe data from REWE.com
2. **Embedder**: Generates vector representations using Ollama
3. **Indexer**: Stores embeddings in Qdrant vector database
4. **CLI**: Provides interface for searching and managing recipes

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgements

- [Ollama](https://ollama.ai/) for embedding technology
- [Qdrant](https://qdrant.tech/) for vector search capabilities
- [REWE](https://www.rewe.de/) for recipe data source

---
