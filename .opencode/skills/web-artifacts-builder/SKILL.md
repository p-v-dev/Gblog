---
name: web-artifacts-builder
description: Suite of tools for creating elaborate, multi-component HTML artifacts using modern frontend web technologies (React, Tailwind CSS, shadcn/ui). Use for complex artifacts requiring state management, routing, or shadcn/ui components — not for simple single-file HTML artifacts.
---

# Web Artifacts Builder

To build powerful frontend artifacts, follow these steps:
1. Initialize the frontend repo using `scripts/init-artifact.sh`
2. Develop your artifact by editing the generated code
3. Bundle all code into a single HTML file using `scripts/bundle-artifact.sh`
4. Display artifact to user
5. (Optional) Test the artifact

**Stack**: React 18 + TypeScript + Vite + Parcel (bundling) + Tailwind CSS + shadcn/ui

## Design & Style Guidelines

To avoid "AI slop", avoid using excessive centered layouts, purple gradients, uniform rounded corners, and Inter font.

## Quick Start

### Step 1: Initialize Project

```bash
bash scripts/init-artifact.sh <project-name>
cd <project-name>
```

This creates a fully configured project with React + TypeScript, Tailwind CSS, shadcn/ui, path aliases, Parcel bundler.

### Step 2: Develop Your Artifact

Edit the generated files. See **Common Development Tasks** below for guidance.

### Step 3: Bundle to Single HTML File

```bash
bash scripts/bundle-artifact.sh
```

Creates `bundle.html` — a self-contained artifact with all JavaScript, CSS, and dependencies inlined.

### Step 4: Share Artifact with User

Share the bundled HTML file in conversation so they can view it as an artifact.

### Step 5: Testing/Visualizing (Optional)

Only perform if necessary or requested. Use Playwright or Puppeteer.

## Reference

- **shadcn/ui components**: https://ui.shadcn.com/docs/components
