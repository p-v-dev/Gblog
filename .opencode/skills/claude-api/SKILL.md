---
name: claude-api
description: |-
  Reference for the Claude API / Anthropic SDK — model ids, pricing, params, streaming, tool use, MCP, agents, caching, token counting, model migration.
  TRIGGER — read BEFORE opening the target file; don't skip because it "looks like a one-liner" — whenever: the prompt names Claude/Anthropic in any form (Claude, Anthropic, Fable, Opus, Sonnet, Haiku, `anthropic`, `@anthropic-ai`, `claude-*`, `us.anthropic.*`, `[1m]`); the user asks about an LLM (pricing/model choice/limits/caching) — never answer from memory; OR the task is LLM-shaped with provider unstated (agent/MCP/tool-definition/multi-agent/RAG/LLM-judge/computer-use; generate/summarize/extract/classify/rewrite/converse over NL; debugging refusals/cutoffs/streaming/tool-calls/tokens).
  SKIP only when another provider is being worked on (overrides all triggers): OpenAI/GPT/Gemini/Llama/Mistral/Cohere/Ollama named in the query; OR `grep -rE 'openai|langchain_openai|google.generativeai|genai|mistralai|cohere|ollama'` over the project hits (run this grep FIRST if no provider named — don't Read the file).
---

# Building LLM-Powered Applications with Claude

This skill helps you build LLM-powered applications with Claude. Choose the right surface based on your needs, detect the project language, then read the relevant language-specific documentation.

## Before You Start

Scan the target file (or, if no target file, the prompt and project) for non-Anthropic provider markers. If you find any, stop and tell the user that this skill produces Claude/Anthropic SDK code; ask whether they want to switch the file to Claude or want a non-Claude implementation. Do not edit a non-Anthropic file with Anthropic SDK calls.

## Output Requirement

When the user asks you to add, modify, or implement a Claude feature, your code must call Claude through one of:

1. **The official Anthropic SDK** for the project's language (`anthropic`, `@anthropic-ai/sdk`, `com.anthropic.*`, etc.). This is the default whenever a supported SDK exists for the project.
2. **Raw HTTP** (`curl`, `requests`, `fetch`, `httpx`, etc.) — only when the user explicitly asks for cURL/REST/raw HTTP, the project is a shell/cURL project, or the language has no official SDK.

Never mix the two — don't reach for `requests`/`fetch` in a Python or TypeScript project just because it feels lighter. Never fall back to OpenAI-compatible shims.

**Never guess SDK usage.** Function names, class names, namespaces, method signatures, and import paths must come from explicit documentation — either the `{lang}/` files in this skill or the official SDK repositories or documentation links listed in `shared/live-sources.md`. If the binding you need is not explicitly documented in the skill files, WebFetch the relevant SDK repo from `shared/live-sources.md` before writing code. Do not infer Ruby/Java/Go/PHP/C# APIs from cURL shapes or from another language's SDK.

## Defaults

Unless the user requests otherwise:

For the Claude model version, please use Claude Opus 4.8, which you can access via the exact model string `claude-opus-4-8`. Please default to using adaptive thinking (`thinking: {type: "adaptive"}`) for anything remotely complicated. And finally, please default to streaming for any request that may involve long input, long output, or high `max_tokens` — it prevents hitting request timeouts. Use the SDK's `.get_final_message()` / `.finalMessage()` helper to get the complete response if you don't need to handle individual stream events

---

## Subcommands

If the User Request at the bottom of this prompt is a bare subcommand string (no prose), search every **Subcommands** table in this document — including any in sections appended below — and follow the matching Action column directly. This lets users invoke specific flows via `/claude-api <subcommand>`. If no table in the document matches, treat the request as normal prose.

| Subcommand | Action |
|---|---|
| `migrate` | Migrate existing Claude API code to a newer model. **Read `shared/model-migration.md` immediately** and follow it in order: Step 0 (confirm scope — ask which files/directories before any edit), Step 1 (classify each file), then the per-target breaking-changes section. Do not summarize the guide — execute it. If the user did not name a target model, ask which model to migrate to in the same turn as the scope question. |

---

## Language Detection

Before reading code examples, determine which language the user is working in:

1. **Look at project files** to infer the language:

   - `*.py`, `requirements.txt`, `pyproject.toml`, `setup.py`, `Pipfile` → **Python** — read from `python/`
   - `*.ts`, `*.tsx`, `package.json`, `tsconfig.json` → **TypeScript** — read from `typescript/`
   - `*.js`, `*.jsx` (no `.ts` files present) → **TypeScript** — JS uses the same SDK, read from `typescript/`
   - `*.java`, `pom.xml`, `build.gradle` → **Java** — read from `java/`
   - `*.kt`, `*.kts`, `build.gradle.kts` → **Java** — Kotlin uses the Java SDK, read from `java/`
   - `*.scala`, `build.sbt` → **Java** — Scala uses the Java SDK, read from `java/`
   - `*.go`, `go.mod` → **Go** — read from `go/`
   - `*.rb`, `Gemfile` → **Ruby** — read from `ruby/`
   - `*.cs`, `*.csproj` → **C#** — read from `csharp/`
   - `*.php`, `composer.json` → **PHP** — read from `php/`

2. **If multiple languages detected** (e.g., both Python and TypeScript files):

   - Check which language the user's current file or question relates to
   - If still ambiguous, ask: "I detected both Python and TypeScript files. Which language are you using for the Claude API integration?"

3. **If language can't be inferred** (empty project, no source files, or unsupported language):

   - Use AskUserQuestion with options: Python, TypeScript, Java, Go, Ruby, cURL/raw HTTP, C#, PHP
   - If AskUserQuestion is unavailable, default to Python examples and note: "Showing Python examples. Let me know if you need a different language."

4. **If unsupported language detected** (Rust, Swift, C++, Elixir, etc.):

   - Suggest cURL/raw HTTP examples from `curl/` and note that community SDKs may exist
   - Offer to show Python or TypeScript examples as reference implementations

5. **If user needs cURL/raw HTTP examples**, read from `curl/`.

### Language-Specific Feature Support

| Language   | Tool Runner | Managed Agents | Notes                                 |
| ---------- | ----------- | -------------- | ------------------------------------- |
| Python     | Yes (beta)  | Yes (beta)     | Full support — `@beta_tool` decorator |
| TypeScript | Yes (beta)  | Yes (beta)     | Full support — `betaZodTool` + Zod    |
| Java       | Yes (beta)  | Yes (beta)     | Beta tool use with annotated classes  |
| Go         | Yes (beta)  | Yes (beta)     | `BetaToolRunner` in `toolrunner` pkg  |
| Ruby       | Yes (beta)  | Yes (beta)     | `BaseTool` + `tool_runner` in beta    |
| C#         | Yes (beta)  | Yes (beta)     | `BetaToolRunner` + raw JSON schema    |
| PHP        | Yes (beta)  | Yes (beta)     | `BetaRunnableTool` + `toolRunner()`   |
| cURL       | N/A         | Yes (beta)     | Raw HTTP, no SDK features             |

---

## Which Surface Should I Use?

> **Start simple.** Default to the simplest tier that meets your needs. Single API calls and workflows handle most use cases — only reach for agents when the task genuinely requires open-ended, model-driven exploration.

| Use Case                                        | Tier            | Recommended Surface       | Why                                                          |
| ----------------------------------------------- | --------------- | ------------------------- | ------------------------------------------------------------ |
| Classification, summarization, extraction, Q&A  | Single LLM call | **Claude API**            | One request, one response                                    |
| Batch processing or embeddings                  | Single LLM call | **Claude API**            | Specialized endpoints                                        |
| Multi-step pipelines with code-controlled logic | Workflow        | **Claude API + tool use** | You orchestrate the loop                                     |
| Custom agent with your own tools                | Agent           | **Claude API + tool use** | Maximum flexibility                                          |
| Server-managed stateful agent with workspace    | Agent           | **Managed Agents**        | Anthropic runs the loop and hosts the tool-execution sandbox |
| Persisted, versioned agent configs              | Agent           | **Managed Agents**        | Agents are stored objects; sessions pin to a version         |
| Long-running multi-turn agent with file mounts  | Agent           | **Managed Agents**        | Per-session containers, SSE event stream, Skills + MCP       |

---

## Current Models (cached: 2026-06-04)

| Model             | Model ID            | Context        | Input $/1M | Output $/1M |
| ----------------- | ------------------- | -------------- | ---------- | ----------- |
| Claude Fable 5    | `claude-fable-5`      | 1M             | $10.00     | $50.00      |
| Claude Mythos 5 (Project Glasswing only) | `claude-mythos-5` | 1M | $10.00     | $50.00      |
| Claude Opus 4.8   | `claude-opus-4-8`   | 1M             | $5.00      | $25.00      |
| Claude Opus 4.7   | `claude-opus-4-7`   | 1M             | $5.00      | $25.00      |
| Claude Opus 4.6   | `claude-opus-4-6`   | 1M             | $5.00      | $25.00      |
| Claude Sonnet 4.6 | `claude-sonnet-4-6` | 1M             | $3.00      | $15.00      |
| Claude Haiku 4.5  | `claude-haiku-4-5`  | 200K           | $1.00      | $5.00       |

**ALWAYS use `claude-opus-4-8` unless the user explicitly names a different model.** This is non-negotiable. Do not use `claude-sonnet-4-6`, `claude-sonnet-4-5`, or any other model unless the user literally says "use sonnet" or "use haiku". Never downgrade for cost — that's the user's decision, not yours. Use `claude-fable-5` only when the user explicitly asks for Claude Fable 5, "fable", or Anthropic's most capable model — it has different API behavior than the Opus family (see below) and pricing that exceeds Opus-tier.

## Thinking & Effort (Quick Reference)

**Fable 5 / Opus 4.8 / 4.7 — Adaptive thinking only:** Use `thinking: {type: "adaptive"}`. `thinking: {type: "enabled", budget_tokens: N}` returns a 400 — adaptive is the only on-mode. On Opus 4.8 and 4.7, `{type: "disabled"}` and omitting `thinking` both work; on Fable 5, an explicit `{type: "disabled"}` returns a 400 — omit the `thinking` param entirely instead. Sampling parameters (`temperature`, `top_p`, `top_k`) are also removed and will 400.
**Opus 4.6 — Adaptive thinking (recommended):** Use `thinking: {type: "adaptive"}`. Claude dynamically decides when and how much to think. No `budget_tokens` needed.
**Effort parameter (GA, no beta header):** Controls thinking depth and overall token spend via `output_config: {effort: "low"|"medium"|"high"|"max"}`.

## Reading Guide

After detecting the language, read the relevant files based on what the user needs:

### Quick Task Reference

**Single text classification/summarization/extraction/Q&A:**
→ Read only `{lang}/claude-api/README.md`

**Chat UI or real-time response display:**
→ Read `{lang}/claude-api/README.md` + `{lang}/claude-api/streaming.md`

**Long-running conversations (may exceed context window):**
→ Read `{lang}/claude-api/README.md` — see Compaction section

**Migrating to a newer model (Fable 5 / Opus 4.8 / Opus 4.7 / Opus 4.6 / Sonnet 4.6) or replacing a retired model:**
→ Read `shared/model-migration.md`

**Prompting or tuning Fable 5 (long turns, effort, verbosity, autonomous runs, sub-agents):**
→ Read `shared/model-migration.md` → Migrating to Fable 5 → Behavioral shifts (prompt-tunable) + Long-running agent recommendations

**Prompt caching / optimize caching / "why is my cache hit rate low":**
→ Read `shared/prompt-caching.md` + `{lang}/claude-api/README.md` (Prompt Caching section)

**Count tokens in a file / prompt / diff ("how many tokens is X"):**
→ Read `shared/token-counting.md` — use `messages.count_tokens`, never `tiktoken`

**Function calling / tool use / agents:**
→ Read `{lang}/claude-api/README.md` + `shared/tool-use-concepts.md` + `{lang}/claude-api/tool-use.md`

**Agent design (tool surface, context management, caching strategy):**
→ Read `shared/agent-design.md`

**Batch processing (non-latency-sensitive):**
→ Read `{lang}/claude-api/README.md` + `{lang}/claude-api/batches.md`

**File uploads across multiple requests:**
→ Read `{lang}/claude-api/README.md` + `{lang}/claude-api/files-api.md`

**Managed Agents (server-managed stateful agents with workspace):**
→ Read `shared/managed-agents-overview.md` + the rest of the `shared/managed-agents-*.md` files.

### Claude API (Full File Reference)

Read the **language-specific Claude API folder** (`{language}/claude-api/`):

1. **`{language}/claude-api/README.md`** — Read this first.
2. **`shared/tool-use-concepts.md`** — Read when the user needs function calling, code execution, memory, or structured outputs.
3. **`shared/agent-design.md`** — Read when designing an agent.
4. **`{language}/claude-api/tool-use.md`** — Read for language-specific tool use code examples.
5. **`{language}/claude-api/streaming.md`** — Read when building chat UIs.
6. **`{language}/claude-api/batches.md`** — Read when processing many requests offline.
7. **`{language}/claude-api/files-api.md`** — Read when sending the same file across multiple requests.
8. **`shared/prompt-caching.md`** — Read when adding or optimizing prompt caching.
9. **`shared/error-codes.md`** — Read when debugging HTTP errors.
10. **`shared/model-migration.md`** — Read when upgrading to newer models.
11. **`shared/live-sources.md`** — WebFetch URLs for fetching the latest official documentation.

## When to Use WebFetch

Use WebFetch to get the latest documentation when:

- User asks for "latest" or "current" information
- Cached data seems incorrect
- User asks about features not covered here

Live documentation URLs are in `shared/live-sources.md`.

## Common Pitfalls

- Don't truncate inputs when passing files or content to the API
- **Fable 5 / Opus 4.8 / 4.7 thinking:** Adaptive only. `thinking: {type: "enabled", budget_tokens: N}` returns 400
- **Prefill removed (Fable 5 and the 4.6/4.7/4.8 family):** Assistant message prefills return a 400 error
- **Fable 5 `refusal` stop reason:** Safety classifiers may decline a request — check `stop_reason` before reading `response.content[0]`
- **Fable 5 tokenizer:** ~30% more tokens for the same content vs Opus-tier models
- **`max_tokens` defaults:** Default to ~16000 for non-streaming, ~64000 for streaming
- **128K output tokens:** Supported on Fable 5, Opus 4.6+ — requires streaming for large values
- **Tool call JSON parsing:** Always parse tool inputs with `json.loads()` / `JSON.parse()`
- **Structured outputs:** Use `output_config: {format: {...}}` instead of deprecated `output_format`
- **Don't reimplement SDK functionality:** The SDK provides high-level helpers — use them
- **Don't define custom types for SDK data structures:** The SDK exports types for all API objects
