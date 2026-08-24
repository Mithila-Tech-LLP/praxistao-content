# 38 | Project: Build a RAG Chatbot

## Table of Contents
1. [Before You Start](#before-you-start)
2. [Project Overview](#project-overview)
3. [Phase 1: Document Ingestion Pipeline](#phase-1-document-ingestion-pipeline)
4. [Phase 2: Retrieval and Generation](#phase-2-retrieval-and-generation)
5. [Phase 3: Conversation Memory](#phase-3-conversation-memory)
6. [Phase 4: Gradio Web UI](#phase-4-gradio-web-ui)
7. [Phase 5: Evaluation and Improvement](#phase-5-evaluation-and-improvement)
8. [Mini Extensions](#mini-extensions)
9. [Exercises](#exercises)

---

## Before You Start

**What you need:**
- Chapters 33 (Embeddings), 34 (Vector DBs), 35 (RAG), 37 (Context Engineering), 45 (LLM APIs and SDKs)
- An API key: Claude (Anthropic), OpenAI, or a local Ollama model
- Python 3.10+

**What you'll build:** A fully functional RAG chatbot that:
- Ingests PDF/text/markdown documents
- Retrieves relevant context for each question
- Maintains conversation history
- Has a clean web UI
- Can evaluate its own answer quality

**Estimated time:** 4-6 hours for the full build, 1-2 hours for the core version.

```
FINAL PRODUCT:
┌─────────────────────────────────────────────────────┐
│  🤖 RAG Chatbot                                     │
│  ─────────────────────────────────────────────────  │
│  📁 Upload Documents: [Browse...]  [Index]          │
│  ─────────────────────────────────────────────────  │
│  💬 Chat                                            │
│                                                     │
│  You: What does the manual say about installation? │
│                                                     │
│  Bot: According to the installation guide (p.3):   │
│       1. Download the installer...                  │
│       [Source: manual.pdf, page 3]                  │
│                                                     │
│  [Type message...]              [Send]              │
└─────────────────────────────────────────────────────┘
```

---

## Project Overview

### What We're Building

A "chat with your documents" application. Users upload documents, then ask questions in natural language. The chatbot retrieves relevant passages and generates grounded answers.

### File Structure

```
rag_chatbot/
├── src/
│   ├── __init__.py
│   ├── ingestion.py      # Document loading + chunking + indexing
│   ├── retrieval.py      # Embedding + search
│   ├── generation.py     # LLM integration + prompting
│   ├── memory.py         # Conversation history
│   └── evaluation.py     # Quality metrics
├── data/
│   └── documents/        # Put your PDFs/docs here
├── app.py                # Gradio web UI
├── config.py             # Settings
├── requirements.txt
└── README.md
```

### Install Dependencies

```bash
pip install anthropic chromadb sentence-transformers PyMuPDF \
            gradio python-dotenv tiktoken numpy
```

---

## Phase 1: Document Ingestion Pipeline

### config.py

```python
# config.py
import os
from dataclasses import dataclass

@dataclass
class Config:
    # Models
    embedding_model: str = "all-MiniLM-L6-v2"
    llm_model: str = "claude-opus-4-7"
    
    # Chunking
    chunk_size: int = 500        # tokens per chunk
    chunk_overlap: int = 50      # overlap between chunks
    
    # Retrieval
    top_k: int = 5               # number of chunks to retrieve
    rerank_top_k: int = 3        # chunks after reranking
    
    # Memory
    max_history_messages: int = 10
    
    # Paths
    chroma_db_path: str = "./chroma_db"
    collection_name: str = "documents"

config = Config()
```

### src/ingestion.py

```python
# src/ingestion.py
import os
import re
import hashlib
from typing import List, Dict, Any
from pathlib import Path

import fitz  # PyMuPDF — same PDF library Chapter 35 uses (PyPDF2 is deprecated)
import chromadb
from sentence_transformers import SentenceTransformer

from config import config


def load_document(filepath: str) -> str:
    """Load text from PDF, markdown, or plain text file."""
    path = Path(filepath)
    
    if path.suffix.lower() == ".pdf":
        return _load_pdf(filepath)
    elif path.suffix.lower() in [".md", ".markdown"]:
        return path.read_text(encoding="utf-8")
    elif path.suffix.lower() == ".txt":
        return path.read_text(encoding="utf-8")
    else:
        raise ValueError(f"Unsupported file type: {path.suffix}")


def _load_pdf(filepath: str) -> str:
    """Extract text from PDF."""
    text_parts = []
    doc = fitz.open(filepath)
    for page_num, page in enumerate(doc):
        text = page.get_text()
        if text.strip():
            text_parts.append(f"[Page {page_num + 1}]\n{text}")
    doc.close()
    return "\n\n".join(text_parts)


def chunk_text(text: str, source: str) -> List[Dict[str, Any]]:
    """
    Split text into overlapping chunks with metadata.

    _load_pdf() inserts "[Page N]" markers into the text stream — this strips
    those back out of the word list before chunking (so they don't pollute the
    chunk text itself) while recording which page each remaining word came
    from, so every chunk's metadata can cite an actual page number.
    """
    raw_words = text.split()

    words: List[str] = []
    word_pages: List[Any] = []
    current_page = None
    i = 0
    while i < len(raw_words):
        if raw_words[i] == "[Page" and i + 1 < len(raw_words) and raw_words[i + 1].endswith("]"):
            current_page = raw_words[i + 1].rstrip("]")
            i += 2
            continue
        words.append(raw_words[i])
        word_pages.append(current_page)
        i += 1

    chunks = []
    i = 0
    chunk_idx = 0
    while i < len(words):
        # Take chunk_size words
        chunk_words = words[i:i + config.chunk_size]
        chunk_pages = word_pages[i:i + config.chunk_size]
        chunk_text = " ".join(chunk_words)
        
        # Skip very short chunks
        if len(chunk_words) < 20:
            break
        
        chunks.append({
            "text": chunk_text,
            "metadata": {
                "source": source,
                "chunk_index": chunk_idx,
                "word_count": len(chunk_words),
                # Page this chunk starts on — None for non-PDF sources, or
                # PDFs with no extractable page markers.
                "page": chunk_pages[0] if chunk_pages and chunk_pages[0] else None,
            }
        })
        
        chunk_idx += 1
        i += config.chunk_size - config.chunk_overlap
    
    return chunks


def get_document_id(filepath: str) -> str:
    """Generate stable ID from file content hash."""
    content = Path(filepath).read_bytes()
    return hashlib.md5(content).hexdigest()[:16]


class DocumentIndexer:
    """Handles indexing documents into ChromaDB."""
    
    def __init__(self):
        self.client = chromadb.PersistentClient(path=config.chroma_db_path)
        self.collection = self.client.get_or_create_collection(
            name=config.collection_name,
            metadata={"hnsw:space": "cosine"}
        )
        self.embedder = SentenceTransformer(config.embedding_model)
    
    def index_document(self, filepath: str) -> Dict[str, Any]:
        """Load, chunk, embed, and store a document."""
        filename = Path(filepath).name
        doc_id = get_document_id(filepath)
        
        # Check if already indexed
        existing = self.collection.get(where={"doc_id": doc_id})
        if existing["ids"]:
            return {
                "status": "already_indexed",
                "filename": filename,
                "chunks": len(existing["ids"])
            }
        
        # Load and chunk
        print(f"Loading {filename}...")
        text = load_document(filepath)
        chunks = chunk_text(text, source=filename)
        
        if not chunks:
            return {"status": "error", "message": "No text extracted"}
        
        # Embed all chunks
        print(f"Embedding {len(chunks)} chunks...")
        texts = [c["text"] for c in chunks]
        embeddings = self.embedder.encode(texts, show_progress_bar=True).tolist()
        
        # Store in ChromaDB
        ids = [f"{doc_id}_{i}" for i in range(len(chunks))]
        metadatas = [
            {**c["metadata"], "doc_id": doc_id, "filename": filename}
            for c in chunks
        ]
        
        self.collection.add(
            ids=ids,
            embeddings=embeddings,
            documents=texts,
            metadatas=metadatas,
        )
        
        print(f"✓ Indexed {filename}: {len(chunks)} chunks")
        return {
            "status": "indexed",
            "filename": filename,
            "chunks": len(chunks),
            "doc_id": doc_id
        }
    
    def get_stats(self) -> Dict:
        """Return collection statistics."""
        count = self.collection.count()
        return {"total_chunks": count}
    
    def delete_document(self, doc_id: str):
        """Remove all chunks for a document."""
        self.collection.delete(where={"doc_id": doc_id})
```

---

## Phase 2: Retrieval and Generation

### src/retrieval.py

```python
# src/retrieval.py
from typing import List, Dict, Any, Optional
import chromadb
from sentence_transformers import SentenceTransformer, CrossEncoder
from config import config


class Retriever:
    """Retrieves relevant document chunks for a query."""
    
    def __init__(self, use_reranker: bool = True):
        self.client = chromadb.PersistentClient(path=config.chroma_db_path)
        self.collection = self.client.get_or_create_collection(config.collection_name)
        self.embedder = SentenceTransformer(config.embedding_model)
        
        self.reranker = None
        if use_reranker:
            try:
                self.reranker = CrossEncoder("cross-encoder/ms-marco-MiniLM-L-6-v2")
            except Exception:
                print("Reranker not available, using embedding similarity only")
    
    def retrieve(self, query: str, top_k: Optional[int] = None) -> List[Dict[str, Any]]:
        """Retrieve top-k most relevant chunks for a query."""
        k = top_k or config.top_k
        
        if self.collection.count() == 0:
            return []
        
        # Embed the query
        query_embedding = self.embedder.encode(query).tolist()
        
        # Search ChromaDB
        results = self.collection.query(
            query_embeddings=[query_embedding],
            n_results=min(k * 2, self.collection.count()),  # Fetch extra for reranking
            include=["documents", "metadatas", "distances"]
        )
        
        chunks = []
        for doc, meta, dist in zip(
            results["documents"][0],
            results["metadatas"][0],
            results["distances"][0]
        ):
            chunks.append({
                "text": doc,
                "metadata": meta,
                "score": 1 - dist  # Convert distance to similarity
            })
        
        # Rerank if available
        if self.reranker and len(chunks) > config.rerank_top_k:
            chunks = self._rerank(query, chunks)
        
        return chunks[:k]
    
    def _rerank(self, query: str, chunks: List[Dict]) -> List[Dict]:
        """Rerank chunks using cross-encoder."""
        pairs = [[query, chunk["text"]] for chunk in chunks]
        scores = self.reranker.predict(pairs)
        
        for chunk, score in zip(chunks, scores):
            chunk["rerank_score"] = float(score)
        
        return sorted(chunks, key=lambda x: x.get("rerank_score", x["score"]), reverse=True)
    
    def format_context(self, chunks: List[Dict]) -> str:
        """Format retrieved chunks into a context string."""
        if not chunks:
            return "No relevant documents found."
        
        parts = []
        for i, chunk in enumerate(chunks, 1):
            source = chunk["metadata"].get("source", "Unknown")
            parts.append(f"[Source {i}: {source}]\n{chunk['text']}")
        
        return "\n\n---\n\n".join(parts)
```

### src/generation.py

```python
# src/generation.py
from typing import List, Dict, Any, Optional
import anthropic
from config import config

SYSTEM_PROMPT = """You are a helpful assistant that answers questions based on provided documents.

## Core Rules
1. ONLY use information from the provided context to answer questions
2. If the context doesn't contain the answer, say "I don't find information about that in the provided documents"
3. Always cite your sources using [Source N] format when referencing specific information
4. Be concise and direct — answer the question, don't pad

## Format
- Use bullet points for lists
- Use **bold** for key terms
- Keep responses under 300 words unless a detailed explanation is genuinely needed
- End with a follow-up question suggestion if appropriate"""


class RAGGenerator:
    """Generates answers using retrieved context."""
    
    def __init__(self):
        self.client = anthropic.Anthropic()
    
    def generate(
        self,
        query: str,
        context: str,
        conversation_history: Optional[List[Dict]] = None,
    ) -> str:
        """Generate a grounded answer for the query."""
        
        # Build the user message with context
        user_message = f"""<context>
{context}
</context>

<question>
{query}
</question>

Please answer the question using only the information in the context above. Cite sources where relevant."""
        
        # Build messages list
        messages = []
        
        # Add conversation history
        if conversation_history:
            messages.extend(conversation_history[-config.max_history_messages * 2:])
        
        messages.append({"role": "user", "content": user_message})
        
        response = self.client.messages.create(
            model=config.llm_model,
            max_tokens=1024,
            system=SYSTEM_PROMPT,
            messages=messages,
        )
        
        return response.content[0].text
    
    def generate_with_sources(
        self,
        query: str,
        chunks: List[Dict[str, Any]],
        conversation_history: Optional[List[Dict]] = None,
    ) -> Dict[str, Any]:
        """Generate answer and return with source metadata."""
        
        # Format context
        context_parts = []
        for i, chunk in enumerate(chunks, 1):
            source = chunk["metadata"].get("source", "Unknown")
            page = chunk["metadata"].get("page")
            label = f"{source}, page {page}" if page else source
            context_parts.append(f"[Source {i}: {label}]\n{chunk['text']}")
        context = "\n\n---\n\n".join(context_parts)
        
        # Generate answer
        answer = self.generate(query, context, conversation_history)
        
        # Extract referenced sources
        sources = []
        for i, chunk in enumerate(chunks, 1):
            if f"[Source {i}]" in answer or f"Source {i}" in answer:
                sources.append({
                    "number": i,
                    "file": chunk["metadata"].get("source"),
                    "page": chunk["metadata"].get("page"),
                    "chunk_index": chunk["metadata"].get("chunk_index"),
                    "relevance_score": chunk.get("score", 0)
                })
        
        return {
            "answer": answer,
            "sources": sources,
            "chunks_retrieved": len(chunks),
        }
```

---

## Phase 3: Conversation Memory

### src/memory.py

```python
# src/memory.py
from typing import List, Dict, Optional
from dataclasses import dataclass, field
import json
from pathlib import Path


@dataclass
class Message:
    role: str   # "user" or "assistant"
    content: str
    metadata: Dict = field(default_factory=dict)


class ConversationMemory:
    """Manages conversation history for multi-turn RAG."""
    
    def __init__(self, max_messages: int = 20, session_file: Optional[str] = None):
        self.max_messages = max_messages
        self.session_file = session_file
        self.messages: List[Message] = []
        
        if session_file and Path(session_file).exists():
            self.load(session_file)
    
    def add_user_message(self, content: str, **metadata):
        self.messages.append(Message(role="user", content=content, metadata=metadata))
        self._trim()
    
    def add_assistant_message(self, content: str, **metadata):
        self.messages.append(Message(role="assistant", content=content, metadata=metadata))
        self._trim()
    
    def get_history(self) -> List[Dict]:
        """Get history in Anthropic API format."""
        return [
            {"role": msg.role, "content": msg.content}
            for msg in self.messages
        ]
    
    def _trim(self):
        """Keep only the most recent messages."""
        if len(self.messages) > self.max_messages:
            # Always keep pairs (user + assistant)
            self.messages = self.messages[-self.max_messages:]
    
    def clear(self):
        self.messages = []
    
    def save(self, filepath: str):
        data = [
            {"role": m.role, "content": m.content, "metadata": m.metadata}
            for m in self.messages
        ]
        Path(filepath).write_text(json.dumps(data, indent=2))
    
    def load(self, filepath: str):
        data = json.loads(Path(filepath).read_text())
        self.messages = [
            Message(role=d["role"], content=d["content"], metadata=d.get("metadata", {}))
            for d in data
        ]
    
    def get_summary_for_context(self) -> str:
        """Brief summary of conversation so far."""
        if not self.messages:
            return ""
        
        topics = []
        for msg in self.messages:
            if msg.role == "user":
                topics.append(msg.content[:50] + "..." if len(msg.content) > 50 else msg.content)
        
        return f"Previous questions: {'; '.join(topics[-3:])}"
    
    @property
    def turn_count(self):
        return sum(1 for m in self.messages if m.role == "user")
```

---

## Phase 4: Gradio Web UI

### app.py

```python
# app.py
import gradio as gr
from pathlib import Path
from typing import List, Tuple, Optional

from src.ingestion import DocumentIndexer
from src.retrieval import Retriever
from src.generation import RAGGenerator
from src.memory import ConversationMemory


# Initialize components
indexer = DocumentIndexer()
retriever = Retriever(use_reranker=False)  # Set True if you have GPU
generator = RAGGenerator()
memory = ConversationMemory(max_messages=20)


def index_documents(files) -> str:
    """Handle file upload and indexing."""
    if not files:
        return "No files uploaded."
    
    results = []
    for file in files:
        result = indexer.index_document(file.name)
        status = result["status"]
        filename = result.get("filename", "unknown")
        
        if status == "indexed":
            results.append(f"✓ {filename}: {result['chunks']} chunks indexed")
        elif status == "already_indexed":
            results.append(f"⟳ {filename}: already indexed ({result['chunks']} chunks)")
        else:
            results.append(f"✗ {filename}: {result.get('message', 'failed')}")
    
    stats = indexer.get_stats()
    results.append(f"\nTotal in database: {stats['total_chunks']} chunks")
    return "\n".join(results)


def chat(
    message: str,
    history: List[Tuple[str, str]],
    show_sources: bool,
) -> Tuple[str, List[Tuple[str, str]]]:
    """Handle a chat message."""
    if not message.strip():
        return "", history
    
    # Get conversation history in API format
    api_history = []
    for user_msg, assistant_msg in history:
        api_history.append({"role": "user", "content": user_msg})
        api_history.append({"role": "assistant", "content": assistant_msg})
    
    # Retrieve relevant chunks
    chunks = retriever.retrieve(message)
    
    if not chunks:
        response = "I don't have any documents indexed yet. Please upload some documents first."
        history.append((message, response))
        return "", history
    
    # Generate answer
    result = generator.generate_with_sources(message, chunks, api_history)
    answer = result["answer"]
    
    # Add source citations to response if enabled
    if show_sources and result["sources"]:
        source_lines = ["\n\n**Sources:**"]
        for src in result["sources"]:
            page = src.get("page")
            label = f"{src['file']}, page {page}" if page else src["file"]
            source_lines.append(f"- {label}")
        answer += "\n".join(source_lines)
    
    history.append((message, answer))
    return "", history


def clear_chat():
    memory.clear()
    return [], ""


# Build Gradio UI
with gr.Blocks(title="RAG Chatbot", theme=gr.themes.Soft()) as demo:
    gr.Markdown("# 📚 RAG Chatbot\nChat with your documents using AI")
    
    with gr.Tab("💬 Chat"):
        chatbot = gr.Chatbot(height=500, label="Conversation")
        
        with gr.Row():
            msg_input = gr.Textbox(
                placeholder="Ask a question about your documents...",
                scale=4,
                label=""
            )
            send_btn = gr.Button("Send", variant="primary", scale=1)
        
        with gr.Row():
            show_sources = gr.Checkbox(label="Show sources", value=True)
            clear_btn = gr.Button("Clear chat")
        
        # Event handlers
        send_btn.click(
            fn=chat,
            inputs=[msg_input, chatbot, show_sources],
            outputs=[msg_input, chatbot]
        )
        msg_input.submit(
            fn=chat,
            inputs=[msg_input, chatbot, show_sources],
            outputs=[msg_input, chatbot]
        )
        clear_btn.click(fn=clear_chat, outputs=[chatbot, msg_input])
    
    with gr.Tab("📁 Documents"):
        gr.Markdown("Upload documents to chat with. Supports PDF, TXT, and Markdown.")
        
        file_upload = gr.File(
            label="Upload documents",
            file_count="multiple",
            file_types=[".pdf", ".txt", ".md"]
        )
        index_btn = gr.Button("Index Documents", variant="primary")
        index_output = gr.Textbox(label="Indexing results", lines=8)
        
        index_btn.click(fn=index_documents, inputs=[file_upload], outputs=[index_output])
    
    with gr.Tab("ℹ️ About"):
        gr.Markdown("""
## How this works

1. **Upload** your documents (PDF, TXT, or Markdown)
2. **Index** — documents are split into chunks and embedded as vectors
3. **Chat** — your question is embedded and matched against stored chunks
4. **Answer** — Claude generates an answer grounded in the retrieved context

## Tips
- Ask specific questions for better results
- The chatbot remembers your conversation history
- Check "Show sources" to see which documents were used
""")


if __name__ == "__main__":
    demo.launch(share=False, server_port=7860)
```

### Running the App

```bash
# Create .env file
echo "ANTHROPIC_API_KEY=your_key_here" > .env

# Run the app
python app.py

# Visit http://localhost:7860
```

---

## Phase 5: Evaluation and Improvement

### src/evaluation.py

```python
# src/evaluation.py
from typing import List, Dict, Any, Optional
import anthropic
from dataclasses import dataclass

@dataclass
class EvalResult:
    question: str
    answer: str
    ground_truth: Optional[str]
    faithfulness_score: float   # Is answer grounded in context?
    relevance_score: float      # Is answer relevant to question?
    completeness_score: float   # Does answer fully address question?
    
    @property
    def average_score(self):
        return (self.faithfulness_score + self.relevance_score + self.completeness_score) / 3


class RAGEvaluator:
    """Evaluates RAG system quality using LLM-as-judge."""
    
    def __init__(self):
        self.client = anthropic.Anthropic()
    
    def evaluate_answer(
        self,
        question: str,
        answer: str,
        context: str,
        ground_truth: Optional[str] = None,
    ) -> EvalResult:
        """Evaluate a single QA pair."""
        
        faithfulness = self._score_faithfulness(answer, context)
        relevance = self._score_relevance(question, answer)
        completeness = self._score_completeness(question, answer, ground_truth)
        
        return EvalResult(
            question=question,
            answer=answer,
            ground_truth=ground_truth,
            faithfulness_score=faithfulness,
            relevance_score=relevance,
            completeness_score=completeness,
        )
    
    def _score(self, prompt: str) -> float:
        """Ask the LLM to score 0-10 and return normalized 0-1 score."""
        response = self.client.messages.create(
            model="claude-haiku-4-5",  # Use fast/cheap model for eval
            max_tokens=100,
            messages=[{"role": "user", "content": prompt}]
        )
        text = response.content[0].text.strip()
        
        # Extract number from response
        import re
        numbers = re.findall(r'\b([0-9]|10)\b', text)
        if numbers:
            return int(numbers[0]) / 10.0
        return 0.5  # Default if parsing fails
    
    def _score_faithfulness(self, answer: str, context: str) -> float:
        prompt = f"""Rate from 0-10 how faithfully this answer is grounded in the provided context.
10 = completely grounded, every claim is in the context.
0 = makes up facts not in context.

Context: {context[:1000]}

Answer: {answer}

Score (just the number):"""
        return self._score(prompt)
    
    def _score_relevance(self, question: str, answer: str) -> float:
        prompt = f"""Rate from 0-10 how relevant this answer is to the question.
10 = directly answers the question.
0 = completely off-topic.

Question: {question}
Answer: {answer}

Score (just the number):"""
        return self._score(prompt)
    
    def _score_completeness(
        self, question: str, answer: str, ground_truth: Optional[str]
    ) -> float:
        if not ground_truth:
            # Without ground truth, check if question seems answered
            prompt = f"""Rate from 0-10 how completely this answer addresses the question.
10 = fully answers all aspects.
0 = doesn't address the question at all.

Question: {question}
Answer: {answer}

Score (just the number):"""
        else:
            prompt = f"""Rate from 0-10 how complete this answer is compared to the reference answer.
10 = covers everything in the reference.
0 = misses all key points.

Reference: {ground_truth}
Answer: {answer}

Score (just the number):"""
        return self._score(prompt)
    
    def evaluate_dataset(
        self, 
        qa_pairs: List[Dict[str, str]],
        retriever,
        generator,
    ) -> Dict[str, Any]:
        """Evaluate on a test dataset."""
        results = []
        
        for pair in qa_pairs:
            question = pair["question"]
            ground_truth = pair.get("answer")
            
            # Run RAG
            chunks = retriever.retrieve(question)
            context = "\n\n".join(c["text"] for c in chunks)
            gen_result = generator.generate_with_sources(question, chunks)
            answer = gen_result["answer"]
            
            # Evaluate
            result = self.evaluate_answer(question, answer, context, ground_truth)
            results.append(result)
            
            print(f"Q: {question[:50]}...")
            print(f"  Faithfulness: {result.faithfulness_score:.2f}")
            print(f"  Relevance:    {result.relevance_score:.2f}")
            print(f"  Completeness: {result.completeness_score:.2f}")
            print(f"  Average:      {result.average_score:.2f}\n")
        
        avg_scores = {
            "faithfulness": sum(r.faithfulness_score for r in results) / len(results),
            "relevance": sum(r.relevance_score for r in results) / len(results),
            "completeness": sum(r.completeness_score for r in results) / len(results),
            "average": sum(r.average_score for r in results) / len(results),
        }
        
        return {"per_question": results, "averages": avg_scores}
```

### Running Evaluations

```python
# evaluate.py
from src.retrieval import Retriever
from src.generation import RAGGenerator
from src.evaluation import RAGEvaluator

# Test questions (create 10-20 for your documents)
test_set = [
    {
        "question": "What are the system requirements for installation?",
        "answer": "At least 8GB RAM and 10GB disk space are required."  # optional ground truth
    },
    {
        "question": "How do I reset my password?",
    },
]

retriever = Retriever()
generator = RAGGenerator()
evaluator = RAGEvaluator()

results = evaluator.evaluate_dataset(test_set, retriever, generator)
print(f"\nOverall scores: {results['averages']}")
```

---

## Mini Extensions

### Extension 1: Multi-document Source Filtering (30 min)

```python
# Add to the Gradio UI: let users filter which documents to search

def chat_with_filter(message, history, show_sources, selected_docs):
    # Pass selected_docs to retriever
    chunks = retriever.retrieve(message, filter_sources=selected_docs)
    ...
```

### Extension 2: Conversation Export (30 min)

```python
def export_conversation(history, format="markdown"):
    """Export chat history to file."""
    if format == "markdown":
        lines = ["# Conversation Export\n"]
        for user, assistant in history:
            lines.append(f"**You:** {user}\n")
            lines.append(f"**Bot:** {assistant}\n\n---\n")
        return "\n".join(lines)
    elif format == "json":
        import json
        return json.dumps([
            {"user": u, "assistant": a} 
            for u, a in history
        ], indent=2)
```

### Extension 3: Streaming Responses (30 min)

```python
# Replace generate() with streaming version for better UX
def generate_streaming(self, query, context, history):
    """Stream tokens as they're generated."""
    with self.client.messages.stream(
        model=config.llm_model,
        max_tokens=1024,
        system=SYSTEM_PROMPT,
        messages=history + [{"role": "user", "content": f"Context:\n{context}\n\nQuestion: {query}"}]
    ) as stream:
        for text in stream.text_stream:
            yield text
```

---

## Exercises

1. **Chunking experiment:** Index the same document with chunk sizes 200, 500, and 1000. Compare the quality of answers for short vs. long questions.

2. **Add metadata filtering:** Modify the retriever to accept `filter_by_source` parameter so users can limit search to specific documents.

3. **Conversation-aware retrieval:** Modify the retrieval step to include recent conversation turns in the query (query expansion using conversation context).

4. **Hallucination detector:** Add a post-generation step that checks if every sentence in the answer can be traced to a retrieved chunk.

5. **Performance test:** Measure latency for: embedding query, ChromaDB search, LLM generation. Where is the bottleneck?

---

**[← Chapter 37: Context Engineering](37-context-engineering.md) | [Chapter 39: AI Agents Architecture →](39-ai-agents-architecture.md)**
