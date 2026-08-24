# Chapter 55: Modern Agent Frameworks — LangGraph, CrewAI, and OpenAI Agents SDK

*Chapter 39 taught you what agents are and how to build one from scratch. This chapter teaches you the frameworks professionals use to build production agents: LangGraph for complex workflows, CrewAI for multi-agent teams, and OpenAI's Agents SDK. You will know when to use each one.*

---

## Table of Contents

1. Why Frameworks Exist — The Problem with Scratch-Built Agents
2. The Landscape: Five Frameworks Worth Knowing
3. LangGraph — Agents as State Machines
4. CrewAI — Multi-Agent Teams
5. OpenAI Agents SDK — Minimal and Powerful
6. Claude as an Agent with Tool Use
7. Choosing the Right Framework
8. Mini Project: Research Agent in All Three Frameworks
9. Summary
10. Exercises

---

## 1. Why Frameworks Exist

When you build an agent from scratch (Chapter 39), you quickly hit the same problems:

**Problem 1: No branching.** ReAct agents go in a straight line. Real tasks need conditional logic: "if the search returns nothing, try a different query; if it returns too much, filter it."

**Problem 2: No resumability.** If the agent crashes on step 7 of 12, you restart from step 1. Real production agents need to resume from where they stopped.

**Problem 3: Coordination is hard.** When multiple agents need to collaborate — one searches, one reads, one writes — coordinating them manually becomes spaghetti.

**Problem 4: Debugging is invisible.** When an agent makes the wrong decision, you have no way to replay the reasoning or inspect state at each step.

Frameworks solve these problems. They add structure without taking away flexibility.

---

## 2. The Landscape: Frameworks Worth Knowing

| Framework | Who made it | Best for | Learning curve |
|-----------|-------------|----------|----------------|
| **LangGraph** | LangChain | Complex workflows with branching, state, cycles | Medium |
| **CrewAI** | CrewAI Inc. | Multi-agent teams with defined roles | Low |
| **OpenAI Agents SDK** | OpenAI | Lightweight, production-ready, any model | Low |
| **AutoGen** | Microsoft | Multi-agent conversations and code execution | Medium |
| **Pydantic AI** | Pydantic | Type-safe agents with strong validation | Medium |

We will cover LangGraph, CrewAI, and OpenAI Agents SDK — the three most widely used in 2025.

---

## 3. LangGraph — Agents as State Machines

LangGraph models your agent as a **graph**: nodes are actions (LLM calls, tool calls), edges are transitions between actions (possibly conditional), and state is passed between nodes.

This gives you precise control over agent flow — something you simply cannot achieve with a simple ReAct loop.

### Installing LangGraph

```bash
pip install langgraph langchain-anthropic
```

### The core concepts

```
STATE:   A dictionary that flows through the graph.
         Every node reads from state and writes back to it.

NODE:    A function that processes state and returns updated state.
         Can be an LLM call, a tool, a condition check — anything.

EDGE:    A connection between nodes.
         Can be unconditional (always go to next node)
         or conditional (go to different nodes based on state).

GRAPH:   The complete flow — nodes + edges.
```

### Example: A Research Agent with LangGraph

```python
from typing import TypedDict, Annotated
from langgraph.graph import StateGraph, END
from langgraph.graph.message import add_messages
from langchain_anthropic import ChatAnthropic
from langchain_core.messages import HumanMessage, AIMessage, ToolMessage
from langchain_core.tools import tool
import json

# Step 1: Define the state
class ResearchState(TypedDict):
    messages: Annotated[list, add_messages]   # conversation history
    query: str                                 # original research query
    sources: list[str]                         # URLs found
    findings: list[str]                        # extracted findings
    report: str                                # final report
    attempts: int                              # how many search attempts

# Step 2: Define tools (Chapter 40 covers this in depth)
@tool
def web_search(query: str) -> str:
    """Search the web for information about a topic."""
    # In production: use Tavily, Serper, or Brave Search API
    # For this example, we return a simulated result
    return f"Found 3 results for '{query}': [simulated results]"

@tool
def read_page(url: str) -> str:
    """Read the content of a web page."""
    return f"Content of {url}: [simulated page content]"

tools = [web_search, read_page]
tools_by_name = {t.name: t for t in tools}

# Step 3: Create the LLM with tools bound to it
llm = ChatAnthropic(model="claude-sonnet-4-6")
llm_with_tools = llm.bind_tools(tools)

# Step 4: Define the nodes

def search_node(state: ResearchState) -> ResearchState:
    """Node: search the web based on the query."""
    messages = state["messages"]
    
    # Ask the LLM what to search for
    response = llm_with_tools.invoke(messages)
    return {"messages": [response], "attempts": state.get("attempts", 0) + 1}

def tool_node(state: ResearchState) -> ResearchState:
    """Node: execute any tool calls the LLM made."""
    messages = state["messages"]
    last_message = messages[-1]
    
    results = []
    for tool_call in last_message.tool_calls:
        tool_fn = tools_by_name[tool_call["name"]]
        result = tool_fn.invoke(tool_call["args"])
        results.append(
            ToolMessage(content=str(result), tool_call_id=tool_call["id"])
        )
    
    return {"messages": results}

def write_report_node(state: ResearchState) -> ResearchState:
    """Node: write the final research report."""
    messages = state["messages"]
    messages.append(HumanMessage(content=
        "Based on all the research above, write a comprehensive report."
    ))
    response = llm.invoke(messages)
    return {"messages": [response], "report": response.content}

# Step 5: Define routing logic
def should_continue(state: ResearchState) -> str:
    """Decide what to do after each LLM call."""
    last_message = state["messages"][-1]
    
    # If the LLM wants to use tools, run the tools
    if hasattr(last_message, "tool_calls") and last_message.tool_calls:
        return "use_tools"
    
    # If we have enough information or too many attempts, write the report
    if state.get("attempts", 0) >= 3:
        return "write_report"
    
    # Otherwise, keep searching
    return "search"

# Step 6: Build the graph
graph_builder = StateGraph(ResearchState)

graph_builder.add_node("search", search_node)
graph_builder.add_node("tools", tool_node)
graph_builder.add_node("write_report", write_report_node)

graph_builder.set_entry_point("search")

graph_builder.add_conditional_edges(
    "search",
    should_continue,
    {
        "use_tools": "tools",
        "search": "search",
        "write_report": "write_report",
    }
)
graph_builder.add_edge("tools", "search")        # after tools, search again
graph_builder.add_edge("write_report", END)      # after report, done

graph = graph_builder.compile()

# Step 7: Run the agent
result = graph.invoke({
    "messages": [HumanMessage(content="Research the latest advances in transformer efficiency")],
    "query": "transformer efficiency advances",
    "sources": [],
    "findings": [],
    "report": "",
    "attempts": 0,
})

print(result["report"])
```

### Why LangGraph shines here

The graph makes the flow explicit and debuggable:
```
START → search → (tool calls?) → tools → search → ... → write_report → END
                              ↗                 ↗
                  (conditional routing based on state)
```

You can visualize this graph, inspect state at each node, and add checkpoints to resume after crashes.

---

## 4. CrewAI — Multi-Agent Teams

CrewAI organizes agents as a team with defined roles. Each agent has a name, role, goal, and backstory. A "crew" is the team. A "process" defines how they work together.

### Installing CrewAI

```bash
pip install crewai crewai-tools
```

### Example: A Content Creation Crew

```python
from crewai import Agent, Task, Crew, Process
from crewai_tools import SerperDevTool, ScrapeWebsiteTool
from langchain_anthropic import ChatAnthropic

llm = ChatAnthropic(model="claude-sonnet-4-6", temperature=0.7)
search_tool = SerperDevTool()    # requires SERPER_API_KEY

# Step 1: Define your agents — each is a specialist
researcher = Agent(
    role="Senior Research Analyst",
    goal="Find comprehensive, accurate information on any given topic.",
    backstory="""You are an expert researcher who has worked in 
                 academic publishing for 15 years. You are thorough, 
                 cite sources, and only report verified facts.""",
    tools=[search_tool],
    llm=llm,
    verbose=True,
)

writer = Agent(
    role="Technical Content Writer",
    goal="Write clear, engaging technical content that non-experts can understand.",
    backstory="""You are a former professor turned technical writer. 
                 You explain complex topics with analogies and examples.""",
    llm=llm,
    verbose=True,
)

editor = Agent(
    role="Senior Editor",
    goal="Improve clarity, accuracy, and structure of written content.",
    backstory="""You are a sharp editor with 20 years of experience. 
                 You catch errors, improve flow, and ensure consistency.""",
    llm=llm,
    verbose=True,
)

# Step 2: Define tasks — each is assigned to an agent
research_task = Task(
    description="""Research the topic: {topic}
                   Gather key facts, recent developments, statistics, and expert opinions.
                   Produce a structured research summary.""",
    expected_output="A structured summary with 5-8 key findings, statistics, and sources.",
    agent=researcher,
)

writing_task = Task(
    description="""Using the research summary, write a 500-word article about {topic}.
                   Make it accessible to a general audience. Include an introduction, 
                   3 main sections, and a conclusion.""",
    expected_output="A well-structured 500-word article ready for publication.",
    agent=writer,
    context=[research_task],   # this task reads the output of research_task
)

editing_task = Task(
    description="""Edit the article for clarity, accuracy, and flow.
                   Fix any errors. Improve weak sentences. Verify all claims 
                   are supported by the research.""",
    expected_output="A polished, publication-ready article.",
    agent=editor,
    context=[research_task, writing_task],
)

# Step 3: Create and run the crew
crew = Crew(
    agents=[researcher, writer, editor],
    tasks=[research_task, writing_task, editing_task],
    process=Process.sequential,   # tasks run one after another
    verbose=True,
)

result = crew.kickoff(inputs={"topic": "How transformers work in AI"})
print(result)
```

### CrewAI process types

```python
Process.sequential   # tasks run in order, each can read prior outputs
Process.hierarchical # a manager agent assigns and oversees tasks
                     # (requires manager_llm parameter)
```

### When CrewAI works best

CrewAI shines when:
- You have distinct roles (researcher, writer, coder)
- Tasks have natural dependencies (write after researching)
- You want to simulate a team dynamic

It is easier to get started with than LangGraph, but less flexible for complex branching logic.

---

## 5. OpenAI Agents SDK — Minimal and Powerful

The OpenAI Agents SDK (formerly "Swarm") is the simplest and most production-ready framework. Despite the name, it works with any model including Claude.

### Installing

```bash
pip install openai-agents
```

### Core concepts

```
AGENT:    An LLM with a name, instructions, and optional tools.
          Simple. No roles, no backstory, no crew.

HANDOFF:  An agent can "hand off" to another agent.
          This is how multi-agent coordination happens.
          The current agent passes its context to a different agent.

RUN:      Execute an agent with a user message.
          The SDK handles the tool calling loop automatically.
```

### Example: Customer Support System

*(The `openai-agents` package's API surface moves fast — verify decorator/parameter names against its current docs before relying on this in production.)*

```python
from agents import Agent, Runner, handoff, function_tool
import asyncio

# Define tools — the decorator is `function_tool`, not `tool`
@function_tool
def look_up_order(order_id: str) -> dict:
    """Look up details for an order by ID."""
    # In production: query your database
    return {
        "order_id": order_id,
        "status": "shipped",
        "items": ["Widget A", "Widget B"],
        "estimated_delivery": "2025-06-17",
    }

@function_tool
def issue_refund(order_id: str, reason: str) -> str:
    """Process a refund for an order."""
    return f"Refund processed for order {order_id}. Reason: {reason}. 3-5 business days."

@function_tool
def escalate_to_human(issue_summary: str) -> str:
    """Escalate a complex issue to a human agent."""
    return f"Escalated to human team. Ticket created. Summary: {issue_summary}"

# Define agents
billing_agent = Agent(
    name="Billing Specialist",
    instructions="""You handle billing issues: refunds, charges, payment problems.
                    You have access to order lookup and refund tools.
                    Be empathetic and efficient.""",
    tools=[look_up_order, issue_refund],
)

technical_agent = Agent(
    name="Technical Support",
    instructions="""You solve technical problems with the product.
                    Ask clarifying questions to diagnose the issue.
                    Provide step-by-step solutions.""",
    tools=[escalate_to_human],
)

# Triage agent decides who to hand off to — handoffs are a separate `handoffs`
# parameter, not merged into `tools`
triage_agent = Agent(
    name="Customer Service Triage",
    instructions="""You are the first point of contact for customer support.
                    Greet the customer, understand their issue, and route them 
                    to the right specialist.
                    - Billing/refund issues → Billing Specialist
                    - Technical problems → Technical Support
                    - Unclear → ask a clarifying question""",
    handoffs=[
        handoff(billing_agent, tool_description_override="Transfer to billing specialist"),
        handoff(technical_agent, tool_description_override="Transfer to technical support"),
    ],
)

# Run the system
async def main():
    runner = Runner()
    
    result = await runner.run(
        triage_agent,
        "I was charged twice for my order #12345 and need a refund.",
    )
    
    print(f"Response: {result.final_output}")
    print(f"Handled by: {result.last_agent.name}")

asyncio.run(main())
```

### Handoffs make routing explicit

The key insight in OpenAI Agents SDK is that **routing is just another tool**. When the triage agent calls `handoff(billing_agent)`, it is like calling any other tool — except the tool transfers execution to a different agent with the full conversation context.

```
User: "I was charged twice" 
→ Triage agent reads message
→ Decides to call handoff(billing_agent)
→ Billing agent receives the full conversation
→ Billing agent calls look_up_order("12345")
→ Billing agent calls issue_refund("12345", "duplicate charge")
→ Billing agent responds to user
```

---

## 6. Claude as an Agent with Tool Use Directly

You do not always need a framework. For simple agents, the Anthropic API's tool use is all you need:

```python
import anthropic
import json

client = anthropic.Anthropic()

# Define tools in Anthropic's format
tools = [
    {
        "name": "calculator",
        "description": "Perform mathematical calculations",
        "input_schema": {
            "type": "object",
            "properties": {
                "expression": {
                    "type": "string",
                    "description": "Mathematical expression to evaluate, e.g. '2 + 2 * 3'"
                }
            },
            "required": ["expression"]
        }
    },
    {
        "name": "get_weather",
        "description": "Get current weather for a city",
        "input_schema": {
            "type": "object",
            "properties": {
                "city": {"type": "string"}
            },
            "required": ["city"]
        }
    }
]

def execute_tool(name: str, inputs: dict) -> str:
    if name == "calculator":
        try:
            result = eval(inputs["expression"])  # in production: use safe_eval
            return str(result)
        except Exception as e:
            return f"Error: {e}"
    elif name == "get_weather":
        # In production: call a real weather API
        return f"Weather in {inputs['city']}: 22°C, sunny"
    return "Unknown tool"

def run_agent(user_message: str) -> str:
    """Simple agent loop using Claude directly."""
    messages = [{"role": "user", "content": user_message}]
    
    while True:
        response = client.messages.create(
            model="claude-sonnet-4-6",
            max_tokens=1000,
            tools=tools,
            messages=messages,
        )
        
        # Claude responded without using tools → done
        if response.stop_reason == "end_turn":
            return response.content[0].text
        
        # Claude wants to use tools
        if response.stop_reason == "tool_use":
            # Add Claude's response to conversation
            messages.append({"role": "assistant", "content": response.content})
            
            # Execute each tool call
            tool_results = []
            for block in response.content:
                if block.type == "tool_use":
                    result = execute_tool(block.name, block.input)
                    tool_results.append({
                        "type": "tool_result",
                        "tool_use_id": block.id,
                        "content": result,
                    })
            
            # Add tool results to conversation
            messages.append({"role": "user", "content": tool_results})
            # Continue the loop — Claude will process the results

print(run_agent("What is 127 * 83, and what is the weather in London?"))
```

---

## 7. Choosing the Right Framework

```
You need:                                       Use:
────────────────────────────────────────────────────────────
A simple agent with 1-5 tools                  Claude direct tool use
Complex branching / state machines              LangGraph
Multi-agent teams with distinct roles          CrewAI
Production reliability, clean API              OpenAI Agents SDK
Maximum flexibility, low boilerplate           OpenAI Agents SDK
Long-running workflows with checkpointing      LangGraph
Dynamic agent collaboration                    AutoGen
```

### Decision flowchart

```
Does your agent have branching logic?
├── Yes → LangGraph
└── No: Does it involve multiple specialized agents?
    ├── Yes, with team roles → CrewAI
    └── Yes, with routing → OpenAI Agents SDK
        No: Just use Claude with tools directly
```

---

## 8. Mini Project: Research Agent in All Three Frameworks

Build the same agent three ways to see the differences:

**The agent:** Given a research question, search the web, read relevant pages, synthesize findings, and produce a 300-word summary.

**Implementation 1 (LangGraph):** Use a state machine with search → read → synthesize nodes and conditional routing.

**Implementation 2 (CrewAI):** Use a researcher agent and a writer agent, connected sequentially.

**Implementation 3 (OpenAI Agents SDK):** Use a single agent with handoff to a writer agent.

**What to compare:**
- How many lines of code?
- How easy was it to debug when something went wrong?
- How would you add a new step (e.g., fact-checking)?
- Which would you use for a production system?

---

## Summary

- **LangGraph**: best for complex workflows with state, branching, and cycles. More code, more control.
- **CrewAI**: best for multi-agent teams with distinct roles. Most beginner-friendly.
- **OpenAI Agents SDK**: best for production simplicity. Works with any LLM including Claude.
- **Direct Claude tool use**: always an option for simple agents. No framework needed.

---

## Exercises

**Easy:**

1. Install all three frameworks. Run the "hello world" example from each framework's documentation.

2. Using CrewAI, create a 2-agent "blog post crew": one agent researches a topic, the other writes a short post from the research. Pick any topic you are curious about.

**Medium:**

3. Build the research agent from Section 8 using the OpenAI Agents SDK. Use real tools (Tavily or Serper for search). Time how long it takes to produce a response.

4. Extend the CrewAI example with an additional "Fact Checker" agent who reviews the written content and flags any claims that need verification.

5. Add a human-in-the-loop step to the LangGraph research agent: after gathering information but before writing the report, pause and ask the user: "I found these sources. Should I proceed? (y/n)"

**Hard:**

6. Build the same agent (a coding assistant that can read files, write code, and run tests) in all three frameworks. Document the trade-offs you observed. Which would you choose for a production deployment and why?

7. Add memory (from Chapter 54) to any of the three framework agents. The agent should remember the user's preferences and past projects across sessions.

---

**[← Chapter 42: MCP Servers](42-mcp-servers.md) | [Chapter 43: Multi-Agent Systems →](43-multi-agent-systems.md)**
