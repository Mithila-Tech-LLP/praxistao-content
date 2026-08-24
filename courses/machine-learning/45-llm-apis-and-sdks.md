# 45 | LLM APIs and SDKs

## Table of Contents
1. [Before You Start](#before-you-start)
2. [The LLM API Landscape](#the-llm-api-landscape)
3. [Anthropic API Deep Dive](#anthropic-api-deep-dive)
4. [OpenAI API and Compatibility Layer](#openai-api-and-compatibility-layer)
5. [Open-Source and Local Models](#open-source-and-local-models)
6. [LangChain and LlamaIndex](#langchain-and-llamaindex)
7. [Cost Optimization Strategies](#cost-optimization-strategies)
8. [Rate Limiting and Batching](#rate-limiting-and-batching)
9. [Mini Projects](#mini-projects)
10. [Exercises](#exercises)

---

## Before You Start

**Prerequisites:**
- Basic Python and HTTP/API concepts
- You've used at least one LLM API in earlier chapters

**What you'll build:** A unified LLM client that can switch between Anthropic, OpenAI, and local models with a consistent interface.

**Why this matters:** The LLM API market is fast-moving. Good code doesn't get tied to one provider.

---

## The LLM API Landscape

```
MAJOR LLM PROVIDERS:

Commercial:
  Anthropic ─── Claude Opus, Sonnet, Haiku
  OpenAI ─────── GPT-4o, GPT-4o-mini, o1
  Google ─────── Gemini 1.5 Pro, Gemini Flash
  Mistral ─────── Mistral Large, Mistral 7B
  Cohere ──────── Command R+

Open-Source (run locally):
  Meta ─────────── Llama 3.1 (8B, 70B, 405B)
  Mistral ─────────Mistral 7B (open weights)
  Google ──────────Gemma 2 (2B, 9B, 27B)
  Microsoft ───────Phi-3 Mini (3.8B)

APIs for open-source:
  Together AI ───── Run any open model via API
  Groq ─────────── Ultra-fast inference (Llama, Mixtral)
  Replicate ──────── Any model via API
  Ollama ──────────── Run locally, OpenAI-compatible API
```

### Choosing a Provider

| Factor | Anthropic Claude | OpenAI GPT-4o | Llama (local) |
|--------|-----------------|---------------|----------------|
| Quality | Excellent | Excellent | Good (70B) |
| Speed | Fast | Fast | Depends on hardware |
| Cost | $3-15/1M tokens | $2.5-10/1M tokens | $0 (hardware cost) |
| Privacy | Data sent to cloud | Data sent to cloud | All local |
| Context window | 200K | 128K | Varies |
| Fine-tuning | No public API | Yes | Yes (fully) |

---

## Anthropic API Deep Dive

### Basic Message API

```python
import anthropic

client = anthropic.Anthropic(
    api_key="your_key"  # or set ANTHROPIC_API_KEY env var
)

# Simple message
response = client.messages.create(
    model="claude-opus-4-7",
    max_tokens=1024,
    messages=[
        {"role": "user", "content": "Explain transformers in one paragraph."}
    ]
)

print(response.content[0].text)
print(f"Input tokens: {response.usage.input_tokens}")
print(f"Output tokens: {response.usage.output_tokens}")
```

### System Prompts and Multi-turn

```python
# System prompt + multi-turn conversation
messages = [
    {"role": "user", "content": "My name is Aditya"},
    {"role": "assistant", "content": "Hello Aditya! How can I help you?"},
    {"role": "user", "content": "What's my name?"},
]

response = client.messages.create(
    model="claude-opus-4-7",
    max_tokens=100,
    system="You are a helpful assistant who remembers user details.",
    messages=messages
)
# Response: "Your name is Aditya, as you mentioned earlier."
```

### Streaming

```python
# Stream tokens as they arrive (much better UX)
with client.messages.stream(
    model="claude-opus-4-7",
    max_tokens=1000,
    messages=[{"role": "user", "content": "Write a poem about clouds."}]
) as stream:
    for text in stream.text_stream:
        print(text, end="", flush=True)
    print()  # Newline at end

# Get final message with usage info
final = stream.get_final_message()
print(f"\nTokens: {final.usage.input_tokens} in, {final.usage.output_tokens} out")
```

### Vision (Images)

```python
import base64
from pathlib import Path

def encode_image(image_path: str) -> str:
    return base64.b64encode(Path(image_path).read_bytes()).decode()

# Send an image for analysis
response = client.messages.create(
    model="claude-opus-4-7",
    max_tokens=500,
    messages=[{
        "role": "user",
        "content": [
            {
                "type": "image",
                "source": {
                    "type": "base64",
                    "media_type": "image/jpeg",
                    "data": encode_image("photo.jpg"),
                },
            },
            {"type": "text", "text": "What's in this image?"}
        ],
    }]
)
print(response.content[0].text)

# Or use URL directly
response = client.messages.create(
    model="claude-opus-4-7",
    max_tokens=200,
    messages=[{
        "role": "user",
        "content": [
            {
                "type": "image",
                "source": {"type": "url", "url": "https://example.com/photo.jpg"},
            },
            {"type": "text", "text": "Describe this image."}
        ]
    }]
)
```

### Batch API (for large-scale processing)

```python
# Process thousands of requests asynchronously (50% cheaper!)
batch = client.messages.batches.create(
    requests=[
        {
            "custom_id": f"request_{i}",
            "params": {
                "model": "claude-haiku-4-5",
                "max_tokens": 100,
                "messages": [{"role": "user", "content": f"Classify: {text}"}]
            }
        }
        for i, text in enumerate(["Great product!", "Terrible service", "It's okay"])
    ]
)

print(f"Batch ID: {batch.id}, Status: {batch.processing_status}")

# Check status later
import time
while True:
    batch = client.messages.batches.retrieve(batch.id)
    if batch.processing_status == "ended":
        break
    time.sleep(60)

# Retrieve results
for result in client.messages.batches.results(batch.id):
    print(f"{result.custom_id}: {result.result.message.content[0].text}")
```

---

## OpenAI API and Compatibility Layer

### OpenAI SDK

```python
from openai import OpenAI

client = OpenAI(api_key="your_key")

response = client.chat.completions.create(
    model="gpt-4o",
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Tell me a joke."}
    ],
    max_tokens=200,
    temperature=0.7,
)

print(response.choices[0].message.content)
print(f"Tokens: {response.usage.total_tokens}")
```

### OpenAI-Compatible APIs

Many providers support the OpenAI format, so you can swap with just a URL change:

```python
from openai import OpenAI

# Groq (ultra-fast, OpenAI-compatible)
client = OpenAI(
    api_key="your_groq_key",
    base_url="https://api.groq.com/openai/v1"
)

response = client.chat.completions.create(
    model="llama-3.1-70b-versatile",  # Use Llama via Groq
    messages=[{"role": "user", "content": "Hello!"}],
    max_tokens=100,
)

# Together AI
client = OpenAI(
    api_key="your_together_key",
    base_url="https://api.together.xyz/v1"
)

# Local Ollama (zero cost!)
client = OpenAI(
    api_key="ollama",  # Any string works
    base_url="http://localhost:11434/v1"
)
response = client.chat.completions.create(
    model="llama3.2",
    messages=[{"role": "user", "content": "Explain neural networks."}]
)
```

---

## Open-Source and Local Models

### Ollama Setup

```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Pull a model
ollama pull llama3.2          # 2GB, fast
ollama pull llama3.1:8b       # 5GB, better quality
ollama pull mistral:7b        # 4GB, strong coder
ollama pull nomic-embed-text  # For embeddings

# Run interactively
ollama run llama3.2

# Start as API server (OpenAI-compatible)
ollama serve  # Runs on localhost:11434
```

### Using Ollama with Python

```python
import httpx
import json

def ollama_generate(prompt: str, model: str = "llama3.2") -> str:
    """Generate with Ollama REST API directly."""
    response = httpx.post(
        "http://localhost:11434/api/generate",
        json={
            "model": model,
            "prompt": prompt,
            "stream": False,
        },
        timeout=120,
    )
    return response.json()["response"]


def ollama_chat(messages: list, model: str = "llama3.2") -> str:
    """Chat with Ollama."""
    response = httpx.post(
        "http://localhost:11434/api/chat",
        json={
            "model": model,
            "messages": messages,
            "stream": False,
        },
        timeout=120,
    )
    return response.json()["message"]["content"]


# Or use OpenAI SDK for consistency
from openai import OpenAI

ollama = OpenAI(api_key="ollama", base_url="http://localhost:11434/v1")

response = ollama.chat.completions.create(
    model="llama3.2",
    messages=[{"role": "user", "content": "What is 2+2?"}]
)
print(response.choices[0].message.content)
```

### Ollama Embeddings

```python
import httpx
import numpy as np

def embed_text(text: str, model: str = "nomic-embed-text") -> list:
    """Generate embeddings with Ollama."""
    response = httpx.post(
        "http://localhost:11434/api/embeddings",
        json={"model": model, "prompt": text},
        timeout=30,
    )
    return response.json()["embedding"]

# Example: compute similarity
v1 = np.array(embed_text("The cat sat on the mat"))
v2 = np.array(embed_text("A feline rested on a rug"))
v3 = np.array(embed_text("Python programming tutorial"))

def cosine_sim(a, b):
    return np.dot(a, b) / (np.linalg.norm(a) * np.linalg.norm(b))

print(f"cat/feline: {cosine_sim(v1, v2):.3f}")    # High similarity
print(f"cat/python: {cosine_sim(v1, v3):.3f}")    # Low similarity
```

---

## LangChain and LlamaIndex

These frameworks provide higher-level abstractions over LLM APIs. Useful for building complex pipelines but can add complexity.

### LangChain Basics

```python
# pip install langchain langchain-anthropic

from langchain_anthropic import ChatAnthropic
from langchain_core.messages import HumanMessage, SystemMessage
from langchain_core.prompts import ChatPromptTemplate

# Basic chat
llm = ChatAnthropic(model="claude-opus-4-7")

response = llm.invoke([
    SystemMessage(content="You are a helpful assistant."),
    HumanMessage(content="What is the capital of France?")
])
print(response.content)

# Using prompt templates
prompt = ChatPromptTemplate.from_messages([
    ("system", "You are a {role}."),
    ("human", "{question}")
])

chain = prompt | llm

result = chain.invoke({
    "role": "Python expert",
    "question": "What is a generator?"
})
print(result.content)
```

### LlamaIndex for RAG

```python
# pip install llama-index llama-index-llms-anthropic llama-index-embeddings-huggingface

from llama_index.core import VectorStoreIndex, SimpleDirectoryReader
from llama_index.llms.anthropic import Anthropic
from llama_index.embeddings.huggingface import HuggingFaceEmbedding
from llama_index.core import Settings

# Configure
Settings.llm = Anthropic(model="claude-opus-4-7")
Settings.embed_model = HuggingFaceEmbedding(model_name="all-MiniLM-L6-v2")

# Index documents
documents = SimpleDirectoryReader("./docs").load_data()
index = VectorStoreIndex.from_documents(documents)

# Query
query_engine = index.as_query_engine()
response = query_engine.query("What are the main themes in these documents?")
print(response)
```

### When to Use Frameworks vs. Raw SDKs

```
Use RAW SDK when:
  ✓ Simple use cases (single model, basic prompting)
  ✓ You need full control over the API
  ✓ Performance matters (fewer abstraction layers)
  ✓ Building production systems where debugging is critical

Use LangChain/LlamaIndex when:
  ✓ Rapid prototyping of complex pipelines
  ✓ Switching between multiple LLM providers
  ✓ You need pre-built integrations (100+ vector DBs, retrievers, etc.)
  ✓ Team already knows the framework
```

---

## Cost Optimization Strategies

### Strategy 1: Model Selection

```python
TASK_MODEL_MAP = {
    # Simple tasks → cheapest model
    "classification": "claude-haiku-4-5",
    "extraction": "claude-haiku-4-5",
    "summarization": "claude-haiku-4-5",
    
    # Medium complexity → balanced model
    "explanation": "claude-sonnet-4-6",
    "code_generation": "claude-sonnet-4-6",
    "analysis": "claude-sonnet-4-6",
    
    # Complex tasks → best model
    "complex_reasoning": "claude-opus-4-7",
    "research": "claude-opus-4-7",
    "novel_tasks": "claude-opus-4-7",
}

def smart_model_select(task_type: str) -> str:
    return TASK_MODEL_MAP.get(task_type, "claude-sonnet-4-6")
```

### Strategy 2: Prompt Caching

```python
# Cache large static context (see Chapter 37 for full details)
# First call: pay full price
# Subsequent calls within 5 min: ~90% cheaper

cached_response = client.messages.create(
    model="claude-opus-4-7",
    max_tokens=512,
    system=[
        {"type": "text", "text": "You are a helpful assistant."},
        {
            "type": "text",
            "text": large_document,  # Could be 50K+ tokens
            "cache_control": {"type": "ephemeral"}
        }
    ],
    messages=[{"role": "user", "content": query}]
)
```

### Strategy 3: Max Tokens Management

```python
def estimate_needed_tokens(task_type: str) -> int:
    """Don't over-provision max_tokens."""
    limits = {
        "yes_no_classification": 5,
        "label_extraction": 20,
        "short_summary": 150,
        "paragraph_response": 300,
        "detailed_analysis": 1000,
        "long_form_content": 4096,
    }
    return limits.get(task_type, 500)
```

### Strategy 4: Batching

```python
# DON'T: 1000 individual API calls
for text in texts:
    result = client.messages.create(...)  # 1000 calls, full cost

# DO: Use Batch API for non-time-sensitive work (50% cheaper!)
batch = client.messages.batches.create(
    requests=[{"custom_id": str(i), "params": {...}} for i, text in enumerate(texts)]
)
# Wait for batch to complete, retrieve all results at once
```

### Strategy 5: Output Caching

```python
import hashlib
import json
from pathlib import Path
from functools import lru_cache

CACHE_DIR = Path("./llm_cache")
CACHE_DIR.mkdir(exist_ok=True)

def cached_llm_call(prompt: str, model: str, **kwargs) -> str:
    """Cache LLM responses to avoid redundant API calls."""
    cache_key = hashlib.md5(f"{model}:{prompt}:{json.dumps(kwargs, sort_keys=True)}".encode()).hexdigest()
    cache_file = CACHE_DIR / f"{cache_key}.json"
    
    if cache_file.exists():
        cached = json.loads(cache_file.read_text())
        return cached["response"]
    
    # Make the actual API call
    response = client.messages.create(
        model=model,
        messages=[{"role": "user", "content": prompt}],
        **kwargs
    )
    result = response.content[0].text
    
    # Cache the result
    cache_file.write_text(json.dumps({
        "prompt": prompt[:100],
        "response": result,
        "model": model,
    }))
    
    return result
```

---

## Rate Limiting and Batching

```python
import time
import asyncio
from typing import List, Any, Callable
from collections import deque

class RateLimiter:
    """Token bucket rate limiter for API calls."""
    
    def __init__(self, requests_per_minute: int = 50, tokens_per_minute: int = 40_000):
        self.rpm = requests_per_minute
        self.tpm = tokens_per_minute
        self.request_times = deque()
        self.token_count = 0
    
    def wait_if_needed(self, estimated_tokens: int = 1000):
        """Block if we're approaching rate limits."""
        now = time.time()
        
        # Remove requests older than 1 minute
        while self.request_times and now - self.request_times[0] > 60:
            self.request_times.popleft()
        
        # Check request limit
        if len(self.request_times) >= self.rpm:
            sleep_time = 60 - (now - self.request_times[0]) + 1
            print(f"Rate limit: sleeping {sleep_time:.1f}s")
            time.sleep(sleep_time)
        
        self.request_times.append(now)
    
    def process_batch(
        self,
        items: List[Any],
        process_fn: Callable,
        delay_between: float = 0.1,
    ) -> List[Any]:
        """Process a list of items with rate limiting."""
        results = []
        for i, item in enumerate(items):
            self.wait_if_needed()
            result = process_fn(item)
            results.append(result)
            
            if i < len(items) - 1:
                time.sleep(delay_between)
        
        return results


# Usage
limiter = RateLimiter(requests_per_minute=50)

def classify_text(text: str) -> str:
    limiter.wait_if_needed()
    response = client.messages.create(
        model="claude-haiku-4-5",
        max_tokens=10,
        messages=[{"role": "user", "content": f"Classify as POSITIVE/NEGATIVE/NEUTRAL: {text}"}]
    )
    return response.content[0].text.strip()

texts = ["Great product!", "Terrible service", "It works"] * 100
results = limiter.process_batch(texts, classify_text, delay_between=0.05)
```

---

## Mini Projects

### Mini Project 1: Unified LLM Client (2 hours)

**Goal:** Build one client class that works with Claude, OpenAI, and Ollama interchangeably.

```python
# unified_llm.py

from abc import ABC, abstractmethod
from typing import Optional
import anthropic
from openai import OpenAI


class LLMProvider(ABC):
    @abstractmethod
    def generate(self, prompt: str, system: str = None, max_tokens: int = 500) -> str:
        pass
    
    @abstractmethod
    def count_tokens(self, text: str) -> int:
        pass


class ClaudeProvider(LLMProvider):
    def __init__(self, model: str = "claude-opus-4-7"):
        self.client = anthropic.Anthropic()
        self.model = model
    
    def generate(self, prompt, system=None, max_tokens=500) -> str:
        kwargs = {"model": self.model, "max_tokens": max_tokens,
                  "messages": [{"role": "user", "content": prompt}]}
        if system:
            kwargs["system"] = system
        resp = self.client.messages.create(**kwargs)
        return resp.content[0].text
    
    def count_tokens(self, text: str) -> int:
        return len(text) // 4  # Approximation


class OpenAIProvider(LLMProvider):
    def __init__(self, model: str = "gpt-4o-mini", base_url: str = None, api_key: str = None):
        self.client = OpenAI(api_key=api_key or "key", base_url=base_url)
        self.model = model
    
    def generate(self, prompt, system=None, max_tokens=500) -> str:
        messages = []
        if system:
            messages.append({"role": "system", "content": system})
        messages.append({"role": "user", "content": prompt})
        resp = self.client.chat.completions.create(
            model=self.model, messages=messages, max_tokens=max_tokens
        )
        return resp.choices[0].message.content
    
    def count_tokens(self, text: str) -> int:
        return len(text) // 4


class OllamaProvider(OpenAIProvider):
    def __init__(self, model: str = "llama3.2"):
        super().__init__(model=model, base_url="http://localhost:11434/v1", api_key="ollama")


class UnifiedLLM:
    """Abstracts over multiple LLM providers."""
    
    def __init__(self, provider: LLMProvider):
        self.provider = provider
    
    @classmethod
    def claude(cls, model="claude-opus-4-7"):
        return cls(ClaudeProvider(model))
    
    @classmethod
    def openai(cls, model="gpt-4o-mini"):
        return cls(OpenAIProvider(model))
    
    @classmethod
    def ollama(cls, model="llama3.2"):
        return cls(OllamaProvider(model))
    
    def ask(self, prompt: str, system: str = None) -> str:
        return self.provider.generate(prompt, system)


# Usage - swap providers with one line change!
llm = UnifiedLLM.claude()
# llm = UnifiedLLM.openai()
# llm = UnifiedLLM.ollama()

print(llm.ask("What is the capital of France?"))
```

### Mini Project 2: Cost Tracker (1 hour)

```python
# Track actual cost across API calls

PRICING = {
    "claude-opus-4-7":   {"input": 15.0, "output": 75.0},  # per 1M tokens
    "claude-sonnet-4-6": {"input": 3.0,  "output": 15.0},
    "claude-haiku-4-5":  {"input": 0.80, "output": 4.0},
    "gpt-4o":            {"input": 5.0,  "output": 15.0},
    "gpt-4o-mini":       {"input": 0.15, "output": 0.60},
}

class CostTracker:
    def __init__(self):
        self.calls = []
    
    def record(self, model: str, input_tokens: int, output_tokens: int):
        pricing = PRICING.get(model, {"input": 5.0, "output": 15.0})
        cost = (input_tokens * pricing["input"] + output_tokens * pricing["output"]) / 1_000_000
        self.calls.append({"model": model, "in": input_tokens, "out": output_tokens, "cost": cost})
    
    def report(self):
        total = sum(c["cost"] for c in self.calls)
        print(f"Total calls: {len(self.calls)}, Total cost: ${total:.6f}")
        
        by_model = {}
        for call in self.calls:
            by_model.setdefault(call["model"], 0)
            by_model[call["model"]] += call["cost"]
        
        for model, cost in sorted(by_model.items(), key=lambda x: -x[1]):
            print(f"  {model}: ${cost:.6f}")
```

---

## Exercises

1. **Provider comparison:** Run the same 10 prompts through Claude Haiku, GPT-4o-mini, and Ollama Llama3.2. Compare: quality, latency, cost.

2. **Retry logic:** Implement an exponential backoff retry wrapper for API calls (handle RateLimitError, ServerError). Test it by simulating errors.

3. **Async calls:** Convert the `process_batch` function to use `asyncio` and `httpx.AsyncClient` for non-blocking parallel API calls.

4. **Model routing:** Build a router that uses a fast/cheap model first, and only escalates to an expensive model if the first model says "I'm not sure" or the confidence is low.

5. **Cost prediction:** Before running a batch job, predict the estimated cost based on prompt size and model pricing. Warn user if it exceeds $1.

---

**[← Chapter 44: Autonomous Agent Project](44-project-autonomous-agent.md) | [Chapter 46: Observability and Evals →](46-observability-and-evals.md)**
