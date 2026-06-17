---
name: mcp-builder
description: Guide for creating high-quality MCP (Model Context Protocol) servers that enable LLMs to interact with external services through well-designed tools. Use when building MCP servers to integrate external APIs or services, whether in Python (FastMCP) or Node/TypeScript (MCP SDK).
---

# MCP Server Development Guide

## Overview

Create MCP (Model Context Protocol) servers that enable LLMs to interact with external services through well-designed tools. The quality of an MCP server is measured by how well it enables LLMs to accomplish real-world tasks.

---

# Process

## High-Level Workflow

Creating a high-quality MCP server involves four main phases:

### Phase 1: Deep Research and Planning

#### Understand Modern MCP Design

**API Coverage vs. Workflow Tools:**
Balance comprehensive API endpoint coverage with specialized workflow tools. Workflow tools can be more convenient for specific tasks, while comprehensive coverage gives agents flexibility to compose operations.

**Tool Naming and Discoverability:**
Clear, descriptive tool names help agents find the right tools quickly. Use consistent prefixes (e.g., `github_create_issue`, `github_list_repos`) and action-oriented naming.

**Context Management:**
Agents benefit from concise tool descriptions and the ability to filter/paginate results.

**Actionable Error Messages:**
Error messages should guide agents toward solutions with specific suggestions.

#### Study MCP Protocol Documentation

Navigate the MCP specification: Start with the sitemap at `https://modelcontextprotocol.io/sitemap.xml`, then fetch specific pages with `.md` suffix.

Key pages to review:
- Specification overview and architecture
- Transport mechanisms (streamable HTTP, stdio)
- Tool, resource, and prompt definitions

#### Study Framework Documentation

**Recommended stack:**
- **Language**: TypeScript (high-quality SDK support)
- **Transport**: Streamable HTTP for remote servers, stdio for local servers

**Load framework documentation:**
- **MCP Best Practices** — read `reference/mcp_best_practices.md`
- For TypeScript: fetch `https://raw.githubusercontent.com/modelcontextprotocol/typescript-sdk/main/README.md`
- For Python: fetch `https://raw.githubusercontent.com/modelcontextprotocol/python-sdk/main/README.md`

#### Plan Your Implementation

Review the service's API documentation. Prioritize comprehensive API coverage. List endpoints starting with the most common operations.

---

### Phase 2: Implementation

#### Set Up Project Structure

See language-specific guides.

#### Implement Core Infrastructure

Create shared utilities:
- API client with authentication
- Error handling helpers
- Response formatting (JSON/Markdown)
- Pagination support

#### Implement Tools

For each tool:

**Input Schema:**
- Use Zod (TypeScript) or Pydantic (Python)
- Include constraints and clear descriptions
- Add examples in field descriptions

**Tool Description:**
- Concise summary of functionality
- Parameter descriptions
- Return type schema

**Implementation:**
- Async/await for I/O operations
- Proper error handling with actionable messages
- Support pagination where applicable

**Annotations:**
- `readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`

---

### Phase 3: Review and Test

#### Code Quality
Review for: no duplicated code, consistent error handling, full type coverage, clear tool descriptions.

#### Build and Test
- TypeScript: `npm run build`, test with MCP Inspector (`npx @modelcontextprotocol/inspector`)
- Python: verify syntax, test with MCP Inspector

---

### Phase 4: Create Evaluations

After implementing your MCP server, create comprehensive evaluations to test its effectiveness.

1. **Tool Inspection**: List available tools and understand their capabilities
2. **Content Exploration**: Use READ-ONLY operations to explore available data
3. **Question Generation**: Create 10 complex, realistic questions
4. **Answer Verification**: Solve each question yourself to verify answers

Each question should be: independent, read-only, complex, realistic, verifiable, stable.
