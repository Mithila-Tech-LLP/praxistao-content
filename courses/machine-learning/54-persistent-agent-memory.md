# Chapter 54: Persistent Memory — Agents That Never Forget

*Chapter 41 introduced the four types of agent memory. This chapter builds the real thing: a production-ready memory system backed by a database that persists across sessions, can be searched semantically, and gets smarter over time.*

---

## Table of Contents

1. The Problem — Why In-Context Memory Is Not Enough
2. What Persistent Memory Looks Like
3. Architecture — A Three-Layer Memory System
4. Building the Memory Database (SQLite)
5. Semantic Search Over Memories (Embeddings)
6. Memory Retrieval — Finding What Matters
7. Memory Aging and Importance
8. Memory Writing — What to Remember and When
9. The Complete Memory Manager
10. Integrating Memory into Your Agent
11. Production Memory with PostgreSQL
12. Mini Project: Personal AI Assistant with Long-Term Memory
13. Summary
14. Exercises

---

## 1. The Problem — Why In-Context Memory Is Not Enough

Let's be concrete about the problem with stateless AI.

```
Day 1:
You: "I am building a sentiment classifier for Twitter data. 
      My deadline is Friday. I am using BERT."

Agent: [gives helpful advice about BERT for Twitter]

Day 2 (new session):
You: "I am still working on my project. Should I use LoRA?"

Agent: "What project are you referring to? What is your use case?"
```

Every conversation starts from zero. The agent has no idea:
- Who you are
- What projects you are working on
- What decisions you have already made
- What problems you have already encountered and solved
- What your preferences are

This is not a minor inconvenience. For an AI assistant to be genuinely useful across time, it needs to remember.

**The scale problem with in-context memory:**

Even modern long-context models have limits:
- GPT-4: 128K tokens = ~300 pages of text
- Claude 3.5: 200K tokens = ~450 pages of text

After a few months of conversations, you have far more data than fits in any context window. And injecting your entire conversation history into every prompt would be expensive and slow.

The solution: store memories in a database, retrieve only the relevant ones for each conversation.

---

## 2. What Persistent Memory Looks Like

Here is what a memory-enabled agent looks like from the outside:

```
Session 1 (March 1):
You:   "I prefer PyTorch over TensorFlow. Also, I am learning ML — 
        explain things from first principles."
Agent: "Got it! I'll explain everything from scratch and use PyTorch 
        in all code examples."

Session 2 (March 15):
You:   "What optimizer should I use for my image classifier?"
Agent: "For image classification in PyTorch [automatically uses 
        PyTorch — remembered your preference], I recommend AdamW 
        with a cosine annealing schedule. Since you mentioned you 
        are learning ML [remembered you are a learner], let me 
        explain why: ..."

Session 3 (April 2):
You:   "I got stuck on the batch normalization tutorial."
Agent: "Earlier you were working on an image classifier [connects 
        current conversation to prior context]. Batch normalization 
        is especially important in your architecture because... 
        You also mentioned you prefer first-principles explanations, 
        so let's start from why normalization helps training..."
```

The agent is learning about you over time and using that knowledge to be more useful.

---

## 3. Architecture — A Three-Layer Memory System

```
┌─────────────────────────────────────────────────────────────┐
│                      AGENT CONTEXT                           │
│   Current conversation + retrieved memories injected here    │
└───────────────────────────┬─────────────────────────────────┘
                            │  retrieve relevant + inject
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    MEMORY RETRIEVAL                          │
│   Given current conversation, find most relevant memories    │
│   Methods: semantic search (embeddings) + recency + importance│
└───────────────────────────┬─────────────────────────────────┘
                            │  search
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    MEMORY DATABASE                           │
│                                                              │
│   memories table:                                            │
│   id | content | category | embedding | importance | age     │
│                                                              │
│   Types: fact, preference, project, decision, skill, event   │
└─────────────────────────────────────────────────────────────┘
```

**Layer 1: Memory Database** — stores all memories persistently.

**Layer 2: Memory Retrieval** — given a query (the current conversation), finds the most relevant memories.

**Layer 3: Agent Context** — injects relevant memories into the prompt, making the agent behave as if it remembers.

---

## 4. Building the Memory Database (SQLite)

We will start with SQLite — a database stored in a single file. Perfect for a personal AI assistant.

```python
# memory/database.py
import sqlite3
import json
import time
from pathlib import Path
from dataclasses import dataclass, field
from typing import Optional

@dataclass
class Memory:
    content: str
    category: str          # "fact", "preference", "project", "skill", "event"
    importance: float = 1.0   # 0.0 to 1.0
    id: Optional[int] = None
    created_at: Optional[float] = None
    last_accessed: Optional[float] = None
    access_count: int = 0
    embedding: Optional[list[float]] = None

class MemoryDatabase:
    def __init__(self, db_path: str = "~/.ai_memory/memories.db"):
        self.db_path = Path(db_path).expanduser()
        self.db_path.parent.mkdir(parents=True, exist_ok=True)
        self._setup()

    def _setup(self):
        """Create the database schema if it does not exist."""
        with self._conn() as conn:
            conn.execute("""
                CREATE TABLE IF NOT EXISTS memories (
                    id            INTEGER PRIMARY KEY AUTOINCREMENT,
                    content       TEXT NOT NULL,
                    category      TEXT NOT NULL DEFAULT 'fact',
                    importance    REAL NOT NULL DEFAULT 1.0,
                    created_at    REAL NOT NULL,
                    last_accessed REAL,
                    access_count  INTEGER NOT NULL DEFAULT 0,
                    embedding     TEXT    -- JSON-encoded list of floats
                )
            """)
            conn.execute("""
                CREATE INDEX IF NOT EXISTS idx_category ON memories(category)
            """)
            conn.execute("""
                CREATE INDEX IF NOT EXISTS idx_importance ON memories(importance DESC)
            """)

    def _conn(self):
        return sqlite3.connect(str(self.db_path))

    def add(self, memory: Memory) -> int:
        """Store a new memory. Returns the assigned ID."""
        now = time.time()
        embedding_json = json.dumps(memory.embedding) if memory.embedding else None

        with self._conn() as conn:
            cursor = conn.execute("""
                INSERT INTO memories 
                    (content, category, importance, created_at, embedding)
                VALUES (?, ?, ?, ?, ?)
            """, (
                memory.content,
                memory.category,
                memory.importance,
                now,
                embedding_json,
            ))
            return cursor.lastrowid

    def get_all(self, category: Optional[str] = None) -> list[Memory]:
        """Retrieve all memories, optionally filtered by category."""
        with self._conn() as conn:
            if category:
                rows = conn.execute(
                    "SELECT * FROM memories WHERE category = ? ORDER BY importance DESC",
                    (category,)
                ).fetchall()
            else:
                rows = conn.execute(
                    "SELECT * FROM memories ORDER BY importance DESC"
                ).fetchall()
        return [self._row_to_memory(r) for r in rows]

    def get_recent(self, limit: int = 20) -> list[Memory]:
        """Retrieve the most recently created memories."""
        with self._conn() as conn:
            rows = conn.execute(
                "SELECT * FROM memories ORDER BY created_at DESC LIMIT ?",
                (limit,)
            ).fetchall()
        return [self._row_to_memory(r) for r in rows]

    def record_access(self, memory_id: int):
        """Record that a memory was accessed (for tracking usage)."""
        with self._conn() as conn:
            conn.execute("""
                UPDATE memories 
                SET last_accessed = ?, access_count = access_count + 1
                WHERE id = ?
            """, (time.time(), memory_id))

    def update_importance(self, memory_id: int, importance: float):
        """Update a memory's importance score."""
        with self._conn() as conn:
            conn.execute(
                "UPDATE memories SET importance = ? WHERE id = ?",
                (importance, memory_id)
            )

    def delete(self, memory_id: int):
        """Remove a memory."""
        with self._conn() as conn:
            conn.execute("DELETE FROM memories WHERE id = ?", (memory_id,))

    def update_embedding(self, memory_id: int, embedding: list[float]):
        """Store the embedding vector for a memory."""
        with self._conn() as conn:
            conn.execute(
                "UPDATE memories SET embedding = ? WHERE id = ?",
                (json.dumps(embedding), memory_id)
            )

    def count(self) -> int:
        with self._conn() as conn:
            return conn.execute("SELECT COUNT(*) FROM memories").fetchone()[0]

    def _row_to_memory(self, row) -> Memory:
        id_, content, category, importance, created_at, last_accessed, access_count, emb = row
        return Memory(
            id=id_,
            content=content,
            category=category,
            importance=importance,
            created_at=created_at,
            last_accessed=last_accessed,
            access_count=access_count,
            embedding=json.loads(emb) if emb else None,
        )
```

---

## 5. Semantic Search Over Memories

Text search (searching for exact words) is not good enough. We need **semantic search**: finding memories that are conceptually related to the query, even if they use different words.

For example: the query "what optimizer do you prefer?" should retrieve the memory "I like AdamW for its stability" — even though "AdamW" and "optimizer" are not in the same sentence.

We achieve this with embeddings (Chapter 33): convert each memory to a vector, convert the query to a vector, find memories whose vector is closest to the query vector.

```python
# memory/embeddings.py
import numpy as np
from sentence_transformers import SentenceTransformer

# Same free, local, no-API-key embedding model Chapter 41 used — a real,
# working default rather than a stub you must fill in before anything else
# in this chapter works. Swap in Voyage AI or OpenAI embeddings later if you
# want higher-quality vectors in production.
_embedder = SentenceTransformer("all-MiniLM-L6-v2")

def embed_text(text: str) -> list[float]:
    """Convert text to an embedding vector."""
    return _embedder.encode(text).tolist()

def cosine_similarity(a: list[float], b: list[float]) -> float:
    """How similar are two vectors? Returns 0 to 1 (1 = identical)."""
    a_arr = np.array(a)
    b_arr = np.array(b)
    
    dot = np.dot(a_arr, b_arr)
    norm_a = np.linalg.norm(a_arr)
    norm_b = np.linalg.norm(b_arr)
    
    if norm_a == 0 or norm_b == 0:
        return 0.0
    
    return float(dot / (norm_a * norm_b))

def find_similar(
    query_embedding: list[float],
    memories_with_embeddings: list[tuple],  # list of (memory, embedding)
    top_k: int = 5,
    min_similarity: float = 0.5,
) -> list[tuple]:
    """Find the top-k most similar memories to the query."""
    scored = []
    for memory, embedding in memories_with_embeddings:
        if embedding is None:
            continue
        sim = cosine_similarity(query_embedding, embedding)
        if sim >= min_similarity:
            scored.append((memory, sim))
    
    # Sort by similarity, highest first
    scored.sort(key=lambda x: x[1], reverse=True)
    return scored[:top_k]
```

---

## 6. Memory Retrieval — Combining Signals

Good memory retrieval is not just "find the most semantically similar memories." It combines three signals:

1. **Semantic relevance**: how closely related is the memory to the current query?
2. **Recency**: memories from the last week matter more than memories from a year ago
3. **Importance**: memories explicitly marked as important (user preferences, key facts) should surface more often

```python
# memory/retrieval.py
import time
import math
from .database import Memory, MemoryDatabase
from .embeddings import embed_text, find_similar

def retrieval_score(
    memory: Memory,
    semantic_similarity: float,
    recency_weight: float = 0.2,
    importance_weight: float = 0.3,
    semantic_weight: float = 0.5,
) -> float:
    """
    Combine three signals into a single retrieval score.
    
    The weights should sum to 1.0.
    """
    # Recency score: 1.0 for very recent, decays exponentially
    if memory.created_at:
        age_days = (time.time() - memory.created_at) / 86400
        recency = math.exp(-age_days / 30)   # half-life of ~30 days
    else:
        recency = 0.5

    score = (
        semantic_weight   * semantic_similarity +
        recency_weight    * recency +
        importance_weight * memory.importance
    )
    return score

class MemoryRetriever:
    def __init__(self, db: MemoryDatabase):
        self.db = db

    def retrieve(
        self,
        query: str,
        top_k: int = 5,
        categories: list[str] | None = None,
    ) -> list[Memory]:
        """
        Given a query string, return the most relevant memories.
        """
        # Get all memories with embeddings
        all_memories = self.db.get_all(category=None)
        memories_with_embeddings = [
            (m, m.embedding) for m in all_memories
            if m.embedding is not None
            and (categories is None or m.category in categories)
        ]

        if not memories_with_embeddings:
            # Fall back to recency if no embeddings
            return self.db.get_recent(top_k)

        # Embed the query
        query_embedding = embed_text(query)

        # Find similar memories
        similar = find_similar(query_embedding, memories_with_embeddings, top_k=top_k * 2)

        # Re-rank by combined score
        scored = [
            (memory, retrieval_score(memory, similarity))
            for memory, similarity in similar
        ]
        scored.sort(key=lambda x: x[1], reverse=True)

        # Record accesses
        top_memories = [m for m, _ in scored[:top_k]]
        for memory in top_memories:
            if memory.id:
                self.db.record_access(memory.id)

        return top_memories
```

---

## 7. Memory Writing — What to Remember and When

The hardest part of memory systems is deciding what to remember. If you remember everything, you end up with too much noise. If you remember too little, the system is not useful.

```python
# memory/writer.py
from anthropic import Anthropic
from .database import Memory, MemoryDatabase
from .embeddings import embed_text

EXTRACTION_PROMPT = """
You are a memory extraction system for an AI assistant.
Given a conversation, extract important facts that should be remembered 
for future conversations. Focus on:
- User preferences (tools, styles, approaches they like)
- Facts about the user (name, role, projects they are working on)
- Decisions made (which approach to use, which tool to pick)
- Important constraints (deadlines, limitations, requirements)
- Skills and knowledge level

For each fact, output:
MEMORY: [the fact to remember, in one clear sentence]
CATEGORY: [fact | preference | project | skill | event]
IMPORTANCE: [0.0 to 1.0, where 1.0 is very important]

Only extract things worth remembering long-term.
Do NOT extract: greetings, general questions, temporary context.
"""

client = Anthropic()

def extract_memories(conversation: list[dict]) -> list[dict]:
    """
    Given a conversation, extract memories worth saving.
    Returns list of dicts with keys: content, category, importance
    """
    formatted = "\n".join(
        f"{msg['role'].upper()}: {msg['content']}"
        for msg in conversation
    )

    response = client.messages.create(
        model="claude-sonnet-4-6",
        max_tokens=1000,
        system=EXTRACTION_PROMPT,
        messages=[{"role": "user", "content": f"Extract memories from:\n\n{formatted}"}]
    )

    text = response.content[0].text
    memories = []
    
    for block in text.strip().split("\n\n"):
        lines = block.strip().split("\n")
        if len(lines) < 3:
            continue
        
        memory = {}
        for line in lines:
            if line.startswith("MEMORY:"):
                memory["content"] = line.replace("MEMORY:", "").strip()
            elif line.startswith("CATEGORY:"):
                memory["category"] = line.replace("CATEGORY:", "").strip()
            elif line.startswith("IMPORTANCE:"):
                try:
                    memory["importance"] = float(line.replace("IMPORTANCE:", "").strip())
                except ValueError:
                    memory["importance"] = 0.7
        
        if "content" in memory and "category" in memory:
            memories.append(memory)
    
    return memories

def save_memories(memories: list[dict], db: MemoryDatabase) -> int:
    """
    Save extracted memories to the database, with embeddings.
    Returns number of memories saved.
    """
    saved = 0
    for mem_dict in memories:
        memory = Memory(
            content=mem_dict["content"],
            category=mem_dict.get("category", "fact"),
            importance=mem_dict.get("importance", 0.7),
        )
        memory_id = db.add(memory)
        
        # Compute and store embedding
        try:
            embedding = embed_text(mem_dict["content"])
            db.update_embedding(memory_id, embedding)
        except Exception:
            pass   # embedding failure is not fatal
        
        saved += 1
    
    return saved
```

---

## 8. The Complete Memory Manager

```python
# memory/manager.py
from .database import MemoryDatabase, Memory
from .retrieval import MemoryRetriever
from .writer import extract_memories, save_memories

class MemoryManager:
    """
    High-level interface for all memory operations.
    This is what your agent uses.
    """
    
    def __init__(self, db_path: str = "~/.ai_memory/memories.db"):
        self.db = MemoryDatabase(db_path)
        self.retriever = MemoryRetriever(self.db)

    def remember_explicitly(
        self,
        content: str,
        category: str = "fact",
        importance: float = 0.8,
    ) -> int:
        """
        Explicitly save something to memory.
        Use this when the user says "remember that..."
        """
        memory = Memory(content=content, category=category, importance=importance)
        memory_id = self.db.add(memory)
        
        try:
            from .embeddings import embed_text
            embedding = embed_text(content)
            self.db.update_embedding(memory_id, embedding)
        except Exception:
            pass
        
        return memory_id

    def learn_from_conversation(self, conversation: list[dict]) -> int:
        """
        Extract and save anything worth remembering from a conversation.
        Call this at the end of each conversation.
        Returns number of new memories saved.
        """
        memories = extract_memories(conversation)
        return save_memories(memories, self.db)

    def recall(self, query: str, top_k: int = 5) -> list[str]:
        """
        Retrieve memories relevant to a query.
        Returns a list of memory content strings.
        """
        memories = self.retriever.retrieve(query, top_k=top_k)
        return [m.content for m in memories]

    def inject_into_prompt(self, query: str, top_k: int = 5) -> str:
        """
        Build a memory context string to inject into a system prompt.
        """
        relevant = self.recall(query, top_k=top_k)
        if not relevant:
            return ""
        
        memories_str = "\n".join(f"- {m}" for m in relevant)
        return f"\n\nRelevant things I know about you:\n{memories_str}\n"

    def stats(self) -> dict:
        return {
            "total_memories": self.db.count(),
            "by_category": {
                cat: len(self.db.get_all(cat))
                for cat in ["fact", "preference", "project", "skill", "event"]
            }
        }
```

---

## 9. Integrating Memory into Your Agent

Here is a complete memory-enabled agent:

```python
# agent_with_memory.py
from anthropic import Anthropic
from memory.manager import MemoryManager

client = Anthropic()

class MemoryAgent:
    def __init__(self, name: str = "AI Assistant"):
        self.name = name
        self.memory = MemoryManager()
        self.conversation = []   # current session

    def chat(self, user_message: str) -> str:
        # 1. Get relevant memories for this query
        memory_context = self.memory.inject_into_prompt(user_message)
        
        # 2. Build system prompt with memory context
        system_prompt = f"""You are {self.name}, a helpful AI assistant.
You remember things about the user from past conversations.
Be concise. If you remember something relevant to the question, use it naturally.
{memory_context}"""

        # 3. Add user message to current conversation
        self.conversation.append({"role": "user", "content": user_message})

        # 4. Get response from the LLM
        response = client.messages.create(
            model="claude-sonnet-4-6",
            max_tokens=1000,
            system=system_prompt,
            messages=self.conversation,
        )
        
        assistant_message = response.content[0].text
        self.conversation.append({"role": "assistant", "content": assistant_message})

        # 5. Check if user wants to save something explicitly
        if "remember" in user_message.lower() and "that" in user_message.lower():
            # Heuristic: if user says "remember that X", save it
            self.memory.remember_explicitly(
                user_message.replace("remember that", "").strip(),
                category="fact",
                importance=0.9
            )

        return assistant_message

    def end_session(self):
        """Call at the end of each conversation to save memories."""
        if len(self.conversation) > 2:  # only save if meaningful conversation
            saved = self.memory.learn_from_conversation(self.conversation)
            print(f"[Memory] Saved {saved} new memories from this session.")
        self.conversation = []

# Usage
agent = MemoryAgent()

print(agent.chat("I am building a sentiment classifier. I prefer using HuggingFace."))
print(agent.chat("My name is Alex and I am a senior ML engineer."))
print(agent.chat("What model architecture should I use?"))

agent.end_session()  # saves memories from this session

# New session — agent will remember Alex and HuggingFace preference
print("\n--- New Session ---")
print(agent.chat("Can you help me with my project?"))
# Agent knows: Alex, senior ML engineer, sentiment classifier, HuggingFace
```

---

## 10. Production Memory with PostgreSQL

For production systems serving multiple users, switch to PostgreSQL. Install with:

```bash
pip install psycopg2-binary pgvector
```

The `pgvector` extension adds native vector similarity search to PostgreSQL:

```sql
-- Enable pgvector
CREATE EXTENSION IF NOT EXISTS vector;

-- Create memories table with vector column
CREATE TABLE memories (
    id            BIGSERIAL PRIMARY KEY,
    user_id       TEXT NOT NULL,
    content       TEXT NOT NULL,
    category      TEXT NOT NULL DEFAULT 'fact',
    importance    FLOAT NOT NULL DEFAULT 1.0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_accessed TIMESTAMPTZ,
    access_count  INT NOT NULL DEFAULT 0,
    embedding     vector(1536)   -- dimension matches your embedding model
);

-- HNSW index for fast approximate nearest neighbor search
CREATE INDEX ON memories USING hnsw (embedding vector_cosine_ops);
```

```python
# production_memory.py
import psycopg2
from psycopg2.extras import execute_values
from pgvector.psycopg2 import register_vector
import numpy as np

class ProductionMemoryDB:
    def __init__(self, connection_string: str):
        self.conn = psycopg2.connect(connection_string)
        # Without this, psycopg2 sends Python lists in its default array
        # format, not pgvector's expected "[1,2,3]" wire format — inserts
        # and similarity queries below will fail without it.
        register_vector(self.conn)

    def add(self, user_id: str, content: str, embedding: list[float],
            category: str = "fact", importance: float = 0.8) -> int:
        with self.conn.cursor() as cur:
            cur.execute("""
                INSERT INTO memories (user_id, content, embedding, category, importance)
                VALUES (%s, %s, %s, %s, %s)
                RETURNING id
            """, (user_id, content, embedding, category, importance))
            self.conn.commit()
            return cur.fetchone()[0]

    def search(self, user_id: str, query_embedding: list[float], 
               top_k: int = 5) -> list[dict]:
        """Semantic search using pgvector's cosine similarity."""
        with self.conn.cursor() as cur:
            cur.execute("""
                SELECT id, content, category, importance,
                       1 - (embedding <=> %s::vector) AS similarity
                FROM memories
                WHERE user_id = %s
                ORDER BY embedding <=> %s::vector
                LIMIT %s
            """, (query_embedding, user_id, query_embedding, top_k))
            
            rows = cur.fetchall()
            return [
                {
                    "id": row[0],
                    "content": row[1],
                    "category": row[2],
                    "importance": row[3],
                    "similarity": row[4],
                }
                for row in rows
            ]
```

---

## 11. Mini Project: Personal AI Assistant with Long-Term Memory

Build a CLI AI assistant that remembers you across sessions.

```python
# assistant.py
"""
Run with: python assistant.py
Type messages and press Enter. Type 'quit' to exit.
Type 'stats' to see memory statistics.
Type 'remember: <fact>' to explicitly save something.
"""
import sys
from agent_with_memory import MemoryAgent

def main():
    agent = MemoryAgent("Aria")
    print("Aria: Hi! I am Aria. I remember things between our conversations.")
    print("     Type 'stats' to see my memory, 'quit' to exit.\n")

    while True:
        try:
            user_input = input("You: ").strip()
        except (EOFError, KeyboardInterrupt):
            break

        if not user_input:
            continue

        if user_input.lower() == "quit":
            agent.end_session()
            print("Aria: Goodbye! I will remember our conversation.")
            break

        if user_input.lower() == "stats":
            stats = agent.memory.stats()
            print(f"Aria: I have {stats['total_memories']} memories about you:")
            for cat, count in stats["by_category"].items():
                if count > 0:
                    print(f"       {cat}: {count}")
            continue

        response = agent.chat(user_input)
        print(f"Aria: {response}\n")

if __name__ == "__main__":
    main()
```

---

## Summary

- Persistent memory transforms a stateless chatbot into a personal AI that improves over time.
- A three-layer architecture: database → retrieval → context injection.
- Memory types: facts, preferences, projects, skills, events.
- Retrieval combines semantic similarity (embeddings), recency, and importance.
- Memory extraction uses the LLM itself to decide what is worth saving.
- SQLite for personal/small systems; PostgreSQL + pgvector for production.

---

## Exercises

**Easy:**

1. Set up the `MemoryDatabase` class and add 5 test memories manually. Retrieve them all. Delete one. Verify it is gone.

2. Write a function `format_memories_for_prompt(memories)` that takes a list of `Memory` objects and returns a nicely formatted string ready to inject into a system prompt.

**Medium:**

3. Implement a memory importance decay system: every time you retrieve memories, reduce the importance of memories that were NOT retrieved by 1%. This means frequently irrelevant memories gradually become less important.

4. Add a "memory consolidation" function: find pairs of memories that are very similar (cosine similarity > 0.9), and replace them with a single merged memory. This reduces memory bloat.

5. Add a simple user interface: when the agent extracts memories from a conversation, show them to the user and ask "Should I remember this? (y/n)" before saving.

**Hard:**

6. Build the complete personal assistant from Section 11 with real embeddings (using OpenAI's `text-embedding-3-small`). Have a 10-turn conversation, end the session, start a new session, and verify that the agent remembers key facts from the first session.

7. Extend the system to support multiple users: each user has their own isolated memory space. Add a simple authentication system (just a username is enough for this exercise). Test that user A's memories do not appear for user B.

---

**[← Chapter 41: Memory Systems for Agents](41-memory-systems-for-agents.md) | [Chapter 42: MCP Servers →](42-mcp-servers.md)**
