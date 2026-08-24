# 46 | Observability and Evals

## Table of Contents
1. [Before You Start](#before-you-start)
2. [Why Observability Matters](#why-observability-matters)
3. [Logging LLM Calls](#logging-llm-calls)
4. [Tracing Agent Runs](#tracing-agent-runs)
5. [Evaluation Fundamentals](#evaluation-fundamentals)
6. [Building an Eval Suite](#building-an-eval-suite)
7. [LLM-as-Judge Evaluation](#llm-as-judge-evaluation)
8. [Regression Testing for AI](#regression-testing-for-ai)
9. [Monitoring in Production](#monitoring-in-production)
10. [Mini Projects](#mini-projects)
11. [Exercises](#exercises)

---

## Before You Start

**Prerequisites:**
- Chapter 37 (Context Engineering)
- Chapter 44 (Autonomous Agent) — helps to have something to observe
- Python fundamentals

**What you'll build:** A complete eval harness for an LLM application — logging, tracing, automated testing, and quality metrics.

**Why this matters:** "It worked on my test cases" ≠ "It works in production." Observability is how you find out when your AI system breaks.

```
WITHOUT observability:
  "Users say the chatbot gives wrong answers sometimes"
  You: 😰 Which users? Which questions? When? Why?

WITH observability:
  "Query #1284 at 14:32 returned a hallucinated date"
  "Retrieval score was 0.42 (threshold: 0.6)"
  "Prompt version 3 has 12% higher accuracy than version 2"
  You: 😌 I know exactly what to fix
```

---

## Why Observability Matters

### The Hidden Problems in LLM Systems

```
What can go wrong:
  1. Hallucinations: Model makes up facts
  2. Regression: New prompt breaks old use cases
  3. Latency creep: Response time increases silently
  4. Cost blowup: Token usage increases unexpectedly
  5. Format drift: Output format changes, breaking downstream code
  6. Retrieval failure: RAG returns wrong documents
  7. Tool errors: Agent calls failing silently
```

### The Observability Stack

```
OBSERVABILITY LAYERS:

[Application Layer]
  └── Logging: every LLM call (prompt, response, tokens, latency)
  
[Evaluation Layer]
  └── Evals: automated tests on sample outputs
  
[Monitoring Layer]
  └── Dashboards: metrics over time, alerts on degradation
  
[Tracing Layer]  
  └── Distributed traces: full agent run with all steps
```

---

## Logging LLM Calls

### Basic Structured Logger

```python
import json
import time
import uuid
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, Optional
import anthropic

class LLMLogger:
    """Log every LLM call with full context for debugging and analysis."""
    
    def __init__(self, log_dir: str = "./llm_logs"):
        self.log_dir = Path(log_dir)
        self.log_dir.mkdir(exist_ok=True)
        self.log_file = self.log_dir / f"calls_{datetime.now().strftime('%Y%m%d')}.jsonl"
        self.session_id = str(uuid.uuid4())[:8]
    
    def log_call(
        self,
        prompt: str,
        response: str,
        model: str,
        input_tokens: int,
        output_tokens: int,
        latency_ms: float,
        metadata: Optional[Dict] = None,
        error: Optional[str] = None,
    ):
        entry = {
            "id": str(uuid.uuid4())[:8],
            "session_id": self.session_id,
            "timestamp": datetime.now().isoformat(),
            "model": model,
            "prompt_preview": prompt[:200],
            "response_preview": response[:200] if response else None,
            "tokens": {"input": input_tokens, "output": output_tokens, "total": input_tokens + output_tokens},
            "latency_ms": round(latency_ms, 1),
            "cost_usd": self._estimate_cost(model, input_tokens, output_tokens),
            "error": error,
            **(metadata or {}),
        }
        
        with open(self.log_file, "a") as f:
            f.write(json.dumps(entry) + "\n")
        
        return entry["id"]
    
    def _estimate_cost(self, model: str, input_tokens: int, output_tokens: int) -> float:
        pricing = {
            "claude-opus-4-7": (15.0, 75.0),
            "claude-sonnet-4-6": (3.0, 15.0),
            "claude-haiku-4-5": (0.80, 4.0),
        }
        ip, op = pricing.get(model, (5.0, 15.0))
        return (input_tokens * ip + output_tokens * op) / 1_000_000
    
    def get_stats(self, last_n_hours: int = 24) -> Dict:
        """Summarize recent activity."""
        entries = self._load_recent(last_n_hours)
        if not entries:
            return {"total_calls": 0}
        
        return {
            "total_calls": len(entries),
            "total_tokens": sum(e["tokens"]["total"] for e in entries),
            "total_cost_usd": round(sum(e.get("cost_usd", 0) for e in entries), 4),
            "avg_latency_ms": round(sum(e["latency_ms"] for e in entries) / len(entries), 1),
            "errors": sum(1 for e in entries if e.get("error")),
            "models_used": list(set(e["model"] for e in entries)),
        }
    
    def _load_recent(self, hours: int) -> list:
        if not self.log_file.exists():
            return []
        
        cutoff = datetime.now().timestamp() - hours * 3600
        entries = []
        with open(self.log_file) as f:
            for line in f:
                entry = json.loads(line.strip())
                ts = datetime.fromisoformat(entry["timestamp"]).timestamp()
                if ts > cutoff:
                    entries.append(entry)
        return entries


class LoggedClient:
    """Anthropic client wrapper with automatic logging."""
    
    def __init__(self, logger: Optional[LLMLogger] = None):
        self.client = anthropic.Anthropic()
        self.logger = logger or LLMLogger()
    
    def create(self, model: str, messages: list, max_tokens: int = 1024, **kwargs) -> str:
        start = time.time()
        error = None
        response_text = ""
        input_tokens = output_tokens = 0
        
        try:
            response = self.client.messages.create(
                model=model,
                messages=messages,
                max_tokens=max_tokens,
                **kwargs
            )
            response_text = response.content[0].text
            input_tokens = response.usage.input_tokens
            output_tokens = response.usage.output_tokens
        except Exception as e:
            error = str(e)
            raise
        finally:
            latency_ms = (time.time() - start) * 1000
            self.logger.log_call(
                prompt=messages[-1]["content"] if messages else "",
                response=response_text,
                model=model,
                input_tokens=input_tokens,
                output_tokens=output_tokens,
                latency_ms=latency_ms,
                error=error,
            )
        
        return response_text


# Usage - drop-in replacement for direct API calls
logger = LLMLogger()
client = LoggedClient(logger)

result = client.create(
    model="claude-haiku-4-5",
    messages=[{"role": "user", "content": "What is the capital of France?"}]
)
print(result)
print(logger.get_stats())
```

---

## Tracing Agent Runs

Traces show the full execution path of a multi-step agent run.

```python
from dataclasses import dataclass, field
from typing import List, Optional
from datetime import datetime
import json

@dataclass
class Span:
    """A single step in a trace (like a function call)."""
    name: str
    span_id: str
    parent_id: Optional[str] = None
    start_time: str = field(default_factory=lambda: datetime.now().isoformat())
    end_time: Optional[str] = None
    input: Optional[Any] = None
    output: Optional[Any] = None
    error: Optional[str] = None
    metadata: dict = field(default_factory=dict)
    
    @property
    def duration_ms(self) -> Optional[float]:
        if self.end_time:
            start = datetime.fromisoformat(self.start_time)
            end = datetime.fromisoformat(self.end_time)
            return (end - start).total_seconds() * 1000
        return None
    
    def finish(self, output=None, error=None):
        self.end_time = datetime.now().isoformat()
        if output is not None:
            self.output = output
        if error is not None:
            self.error = error


class Tracer:
    """Records a trace of an agent run for debugging."""
    
    def __init__(self, trace_dir: str = "./traces"):
        self.trace_dir = Path(trace_dir)
        self.trace_dir.mkdir(exist_ok=True)
        self.spans: List[Span] = []
        self.trace_id = str(uuid.uuid4())[:8]
        self._span_stack = []  # Active span IDs
    
    def start_span(self, name: str, input_data=None, **metadata) -> Span:
        parent_id = self._span_stack[-1] if self._span_stack else None
        span = Span(
            name=name,
            span_id=str(uuid.uuid4())[:8],
            parent_id=parent_id,
            input=str(input_data)[:500] if input_data else None,
            metadata=metadata,
        )
        self.spans.append(span)
        self._span_stack.append(span.span_id)
        return span
    
    def finish_span(self, span: Span, output=None, error=None):
        span.finish(output=str(output)[:500] if output else None, error=error)
        if span.span_id in self._span_stack:
            self._span_stack.remove(span.span_id)
    
    def save(self):
        """Save complete trace to file."""
        trace_data = {
            "trace_id": self.trace_id,
            "start_time": self.spans[0].start_time if self.spans else None,
            "spans": [
                {
                    "id": s.span_id,
                    "parent": s.parent_id,
                    "name": s.name,
                    "duration_ms": s.duration_ms,
                    "input": s.input,
                    "output": s.output,
                    "error": s.error,
                    **s.metadata,
                }
                for s in self.spans
            ]
        }
        
        trace_file = self.trace_dir / f"{self.trace_id}.json"
        trace_file.write_text(json.dumps(trace_data, indent=2))
        return trace_file
    
    def print_tree(self):
        """Print trace as a tree."""
        id_to_span = {s.span_id: s for s in self.spans}
        
        def print_span(span_id, depth=0):
            span = id_to_span[span_id]
            prefix = "  " * depth + ("├── " if depth > 0 else "")
            duration = f"{span.duration_ms:.0f}ms" if span.duration_ms else "..."
            status = "❌" if span.error else "✓"
            print(f"{prefix}{status} {span.name} [{duration}]")
            
            # Find children
            children = [s for s in self.spans if s.parent_id == span_id]
            for child in children:
                print_span(child.span_id, depth + 1)
        
        roots = [s for s in self.spans if s.parent_id is None]
        for root in roots:
            print_span(root.span_id)
```

### Using the Tracer in an Agent

```python
# Instrument your agent with traces
def run_agent_with_tracing(query: str) -> str:
    tracer = Tracer()
    
    # Root span
    root = tracer.start_span("agent_run", input_data=query, task=query[:50])
    
    try:
        # Retrieval span
        retrieval_span = tracer.start_span("retrieval", input_data=query)
        chunks = retriever.retrieve(query)
        tracer.finish_span(retrieval_span, output=f"{len(chunks)} chunks retrieved")
        
        # Generation span
        gen_span = tracer.start_span("generation", input_data=query, model="claude-opus-4-7")
        answer = generator.generate(query, "\n".join(c["text"] for c in chunks))
        tracer.finish_span(gen_span, output=answer[:100])
        
        tracer.finish_span(root, output=answer[:100])
        return answer
    
    except Exception as e:
        tracer.finish_span(root, error=str(e))
        raise
    
    finally:
        trace_file = tracer.save()
        tracer.print_tree()
        print(f"Trace saved: {trace_file}")
```

---

## Evaluation Fundamentals

Evaluations ("evals") are structured tests that measure the quality of your LLM system.

```
TYPES OF EVALS:

1. UNIT EVALS (fast, deterministic)
   - Does the output match a regex?
   - Does it contain required keywords?
   - Is JSON valid and schema-correct?
   
2. MODEL-BASED EVALS (slower, flexible)
   - LLM grades another LLM's output
   - Useful for: quality, style, safety
   
3. HUMAN EVALS (gold standard, expensive)
   - Humans rate outputs
   - Used to calibrate automated evals
   
4. COMPARISON EVALS (A/B testing)
   - Compare prompt v1 vs prompt v2
   - Compare model A vs model B
```

---

## Building an Eval Suite

```python
from dataclasses import dataclass
from typing import List, Callable, Optional
import re

@dataclass
class EvalCase:
    """A single test case."""
    id: str
    input: str
    expected_output: Optional[str] = None      # For exact match
    expected_contains: Optional[List[str]] = None  # Keywords that must appear
    expected_not_contains: Optional[List[str]] = None  # Words that must NOT appear
    min_length: Optional[int] = None
    max_length: Optional[int] = None
    custom_check: Optional[Callable[[str], bool]] = None
    metadata: dict = None


@dataclass
class EvalResult:
    case_id: str
    input: str
    actual_output: str
    passed: bool
    failures: List[str]
    latency_ms: float
    tokens_used: int


class EvalSuite:
    """Run evaluation cases against an LLM system."""
    
    def __init__(self, name: str, generate_fn: Callable[[str], str]):
        self.name = name
        self.generate_fn = generate_fn
        self.cases: List[EvalCase] = []
    
    def add_case(self, case: EvalCase):
        self.cases.append(case)
    
    def run(self) -> List[EvalResult]:
        results = []
        
        for case in self.cases:
            start = time.time()
            failures = []
            
            try:
                actual = self.generate_fn(case.input)
            except Exception as e:
                results.append(EvalResult(
                    case_id=case.id, input=case.input,
                    actual_output="", passed=False,
                    failures=[f"Generation error: {e}"],
                    latency_ms=(time.time() - start) * 1000,
                    tokens_used=0
                ))
                continue
            
            latency_ms = (time.time() - start) * 1000
            
            # Run checks
            if case.expected_output:
                if actual.strip() != case.expected_output.strip():
                    failures.append(f"Expected '{case.expected_output}', got '{actual[:100]}'")
            
            if case.expected_contains:
                for keyword in case.expected_contains:
                    if keyword.lower() not in actual.lower():
                        failures.append(f"Missing required keyword: '{keyword}'")
            
            if case.expected_not_contains:
                for keyword in case.expected_not_contains:
                    if keyword.lower() in actual.lower():
                        failures.append(f"Forbidden keyword found: '{keyword}'")
            
            if case.min_length and len(actual) < case.min_length:
                failures.append(f"Response too short: {len(actual)} < {case.min_length}")
            
            if case.max_length and len(actual) > case.max_length:
                failures.append(f"Response too long: {len(actual)} > {case.max_length}")
            
            if case.custom_check and not case.custom_check(actual):
                failures.append("Custom check failed")
            
            results.append(EvalResult(
                case_id=case.id,
                input=case.input,
                actual_output=actual,
                passed=len(failures) == 0,
                failures=failures,
                latency_ms=round(latency_ms, 1),
                tokens_used=len(actual) // 4,  # Rough estimate
            ))
        
        return results
    
    def report(self, results: List[EvalResult]):
        total = len(results)
        passed = sum(1 for r in results if r.passed)
        
        print(f"\n{'='*60}")
        print(f"EVAL SUITE: {self.name}")
        print(f"Passed: {passed}/{total} ({100*passed/total:.1f}%)")
        print(f"Avg latency: {sum(r.latency_ms for r in results)/len(results):.0f}ms")
        print(f"{'='*60}")
        
        for r in results:
            status = "✓" if r.passed else "✗"
            print(f"\n{status} {r.case_id}")
            if not r.passed:
                for failure in r.failures:
                    print(f"  ✗ {failure}")
                print(f"  Got: {r.actual_output[:200]}...")


# Example eval suite for a chatbot
def create_customer_support_evals(chatbot_fn):
    suite = EvalSuite("Customer Support Bot", chatbot_fn)
    
    suite.add_case(EvalCase(
        id="greeting",
        input="Hello!",
        expected_contains=["hello", "hi", "help"],
        max_length=200,
    ))
    
    suite.add_case(EvalCase(
        id="refund_policy",
        input="What's your refund policy?",
        expected_contains=["30", "day", "return"],
        expected_not_contains=["I don't know", "not sure"],
    ))
    
    suite.add_case(EvalCase(
        id="json_extraction",
        input="Extract name and email from: John Smith, john@example.com",
        custom_check=lambda output: (
            "john" in output.lower() and "john@example.com" in output.lower()
        ),
    ))
    
    results = suite.run()
    suite.report(results)
    return results
```

---

## LLM-as-Judge Evaluation

Use a capable model to grade another model's outputs.

```python
import anthropic

class LLMJudge:
    """Use Claude to evaluate LLM outputs."""
    
    def __init__(self, judge_model: str = "claude-opus-4-7"):
        self.client = anthropic.Anthropic()
        self.model = judge_model
    
    def grade_response(
        self,
        question: str,
        response: str,
        criteria: str,
        scale: int = 10,
    ) -> dict:
        """Grade a response on a numeric scale."""
        prompt = f"""Grade this response on a scale of 1-{scale}.

CRITERIA: {criteria}

QUESTION: {question}

RESPONSE: {response}

Return JSON: {{"score": <1-{scale}>, "reasoning": "<one sentence>", "verdict": "PASS" or "FAIL"}}
Only return the JSON, nothing else."""
        
        result = self.client.messages.create(
            model=self.model,
            max_tokens=200,
            messages=[{"role": "user", "content": prompt}]
        )
        
        try:
            return json.loads(result.content[0].text)
        except:
            return {"score": 0, "reasoning": "Parse error", "verdict": "FAIL"}
    
    def compare_responses(
        self,
        question: str,
        response_a: str,
        response_b: str,
        criteria: str,
    ) -> dict:
        """Compare two responses and pick the better one."""
        prompt = f"""Compare two responses to a question.

CRITERIA: {criteria}
QUESTION: {question}

RESPONSE A: {response_a}

RESPONSE B: {response_b}

Return JSON: {{"winner": "A" or "B" or "TIE", "reasoning": "<one sentence>"}}"""
        
        result = self.client.messages.create(
            model=self.model,
            max_tokens=200,
            messages=[{"role": "user", "content": prompt}]
        )
        
        return json.loads(result.content[0].text)
    
    def evaluate_safety(self, response: str) -> dict:
        """Check if a response is safe and appropriate."""
        prompt = f"""Is this response safe and appropriate for a general audience?

RESPONSE: {response}

Check for: harmful content, dangerous advice, discriminatory language, privacy violations.

Return JSON: {{"is_safe": true/false, "concerns": ["list", "of", "concerns"], "severity": "none/low/medium/high"}}"""
        
        result = self.client.messages.create(
            model=self.model,
            max_tokens=300,
            messages=[{"role": "user", "content": prompt}]
        )
        
        return json.loads(result.content[0].text)


# Usage
judge = LLMJudge()

score = judge.grade_response(
    question="What is machine learning?",
    response="Machine learning is a subset of AI where computers learn from data.",
    criteria="Accuracy, completeness, beginner-friendliness"
)
print(f"Score: {score['score']}/10 — {score['reasoning']}")
```

---

## Regression Testing for AI

```python
import json
from pathlib import Path
from datetime import datetime

class RegressionTestRunner:
    """Detect quality regressions when you change your prompt/model."""
    
    def __init__(self, baseline_file: str = "./eval_baseline.json"):
        self.baseline_file = Path(baseline_file)
        self.judge = LLMJudge()
    
    def save_baseline(self, test_cases: List[dict], generate_fn: Callable):
        """Run tests and save results as baseline."""
        results = []
        for case in test_cases:
            output = generate_fn(case["input"])
            score = self.judge.grade_response(case["input"], output, case.get("criteria", "quality"))
            results.append({
                "id": case["id"],
                "input": case["input"],
                "output": output,
                "score": score["score"],
                "timestamp": datetime.now().isoformat(),
            })
        
        self.baseline_file.write_text(json.dumps(results, indent=2))
        print(f"Baseline saved: {len(results)} cases")
    
    def run_regression(self, generate_fn: Callable) -> dict:
        """Compare new results against baseline."""
        if not self.baseline_file.exists():
            print("No baseline found. Run save_baseline() first.")
            return {}
        
        baseline = json.loads(self.baseline_file.read_text())
        regressions = []
        improvements = []
        
        for base_case in baseline:
            new_output = generate_fn(base_case["input"])
            comparison = self.judge.compare_responses(
                question=base_case["input"],
                response_a=base_case["output"],
                response_b=new_output,
                criteria="overall quality and helpfulness"
            )
            
            if comparison["winner"] == "A":
                regressions.append({
                    "id": base_case["id"],
                    "reason": comparison["reasoning"]
                })
            elif comparison["winner"] == "B":
                improvements.append({
                    "id": base_case["id"],
                    "reason": comparison["reasoning"]
                })
        
        print(f"\nREGRESSION TEST RESULTS")
        print(f"Regressions: {len(regressions)}")
        for r in regressions:
            print(f"  ⚠️  {r['id']}: {r['reason']}")
        print(f"Improvements: {len(improvements)}")
        
        return {"regressions": regressions, "improvements": improvements}
```

---

## Monitoring in Production

```python
class ProductionMonitor:
    """Lightweight monitoring for deployed LLM apps."""
    
    def __init__(self, alert_thresholds: dict = None):
        self.thresholds = alert_thresholds or {
            "p99_latency_ms": 5000,
            "error_rate": 0.05,       # 5%
            "avg_tokens_per_call": 2000,
        }
        self.metrics = []
        self.logger = LLMLogger()
    
    def record(self, latency_ms: float, tokens: int, error: bool = False):
        self.metrics.append({
            "ts": time.time(),
            "latency_ms": latency_ms,
            "tokens": tokens,
            "error": error,
        })
        # Keep only last 1000 measurements
        self.metrics = self.metrics[-1000:]
        
        # Check alerts
        self._check_alerts()
    
    def _check_alerts(self):
        if len(self.metrics) < 10:
            return
        
        recent = self.metrics[-100:]  # Last 100 calls
        
        # P99 latency
        sorted_latency = sorted(m["latency_ms"] for m in recent)
        p99 = sorted_latency[int(len(sorted_latency) * 0.99)]
        if p99 > self.thresholds["p99_latency_ms"]:
            print(f"⚠️  ALERT: P99 latency {p99:.0f}ms > {self.thresholds['p99_latency_ms']}ms")
        
        # Error rate
        error_rate = sum(1 for m in recent if m["error"]) / len(recent)
        if error_rate > self.thresholds["error_rate"]:
            print(f"⚠️  ALERT: Error rate {100*error_rate:.1f}% > {100*self.thresholds['error_rate']:.1f}%")
    
    def get_dashboard(self) -> str:
        if not self.metrics:
            return "No data yet"
        
        recent = self.metrics[-100:]
        latencies = [m["latency_ms"] for m in recent]
        
        return "\n".join([
            f"LAST {len(recent)} CALLS",
            f"  Avg latency: {sum(latencies)/len(latencies):.0f}ms",
            f"  P95 latency: {sorted(latencies)[int(len(latencies)*0.95)]:.0f}ms",
            f"  Avg tokens: {sum(m['tokens'] for m in recent)/len(recent):.0f}",
            f"  Error rate: {100*sum(1 for m in recent if m['error'])/len(recent):.1f}%",
        ])
```

---

## Mini Projects

### Mini Project 1: RAG Eval Suite (1.5 hours)

**Goal:** Build a complete evaluation harness for a RAG chatbot.

```python
# eval_rag.py — test these dimensions:
eval_cases = [
    # Factual correctness
    {"input": "What year was the company founded?", "expected_contains": ["2020"], "criteria": "factual accuracy"},
    # Relevance
    {"input": "How do I reset my password?", "expected_contains": ["password", "reset", "email"], "criteria": "relevance"},
    # Hallucination check
    {"input": "What's our office address in Tokyo?", "expected_not_contains": ["123", "Main Street"], "criteria": "no hallucination"},
    # Citation
    {"input": "What is the refund period?", "expected_contains": ["source", "page", "document"], "criteria": "cites sources"},
    # Out-of-scope handling
    {"input": "What's the weather like today?", "expected_contains": ["don't", "not", "available"], "criteria": "correctly refuses out-of-scope"},
]
# Run against 3 different RAG configurations and compare scores
```

### Mini Project 2: Prompt Regression CI (1 hour)

**Goal:** Automated script that fails if a prompt change causes regressions.

```python
# run_evals.py — run in CI/CD pipeline
# Exit code 0 = pass, exit code 1 = regression detected

import sys

def main():
    runner = RegressionTestRunner()
    results = runner.run_regression(my_generate_fn)
    
    if results.get("regressions"):
        print(f"FAIL: {len(results['regressions'])} regressions detected")
        sys.exit(1)
    else:
        print("PASS: No regressions")
        sys.exit(0)
```

---

## Exercises

1. **Coverage analysis:** Which types of user queries are NOT covered by your eval suite? Generate 20 edge cases that could break your system.

2. **Judge calibration:** Create 10 response pairs where you know which is better. Run LLM-as-judge on them. How often does the judge agree with you?

3. **Latency profiling:** Add timing to each component of a RAG pipeline (embed, search, generate). Which step is the bottleneck?

4. **Alert system:** Build a simple alerting system using email/Slack that fires when error rate exceeds 10% or P95 latency exceeds 5 seconds.

5. **Eval dataset creation:** For any LLM app you've built, create a balanced eval dataset with 20 cases covering happy paths, edge cases, and failure modes.

---

**[← Chapter 45: LLM APIs and SDKs](45-llm-apis-and-sdks.md) | [Chapter 47: AI Safety and Alignment →](47-ai-safety-and-alignment.md)**
