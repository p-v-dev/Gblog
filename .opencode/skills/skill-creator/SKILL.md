---
name: skill-creator
description: Create new skills, modify and improve existing skills, and measure skill performance. Use when users want to create a skill from scratch, edit, or optimize an existing skill, run evals to test a skill, benchmark skill performance with variance analysis, or optimize a skill's description for better triggering accuracy.
---

# Skill Creator

A skill for creating new skills and iteratively improving them.

## Communicating with the user

The skill creator may be used by people across a wide range of familiarity with coding. Pay attention to context cues to understand how to phrase your communication. Briefly explain terms if in doubt.

---

## Creating a skill

### Capture Intent

Start by understanding the user's intent. Extract answers from the conversation history. The user may need to fill gaps, and should confirm before proceeding.

1. What should this skill enable Claude to do?
2. When should this skill trigger?
3. What's the expected output format?
4. Should we set up test cases?

### Interview and Research

Ask questions about edge cases, input/output formats, example files, success criteria, and dependencies.

### Write the SKILL.md

Based on the user interview, fill in:
- **name**: Skill identifier
- **description**: When to trigger, what it does. Make descriptions slightly "pushy" to combat undertriggering.
- The rest of the skill content

### Skill Writing Guide

#### Anatomy of a Skill
```
skill-name/
├── SKILL.md (required)
├── scripts/    - Executable code
├── references/ - Docs loaded as needed
└── assets/     - Templates, icons, fonts
```

#### Progressive Disclosure
1. **Metadata** (name + description) - Always in context (~100 words)
2. **SKILL.md body** - In context when skill triggers (<500 lines ideal)
3. **Bundled resources** - As needed

#### Writing Patterns
Prefer imperative form. Explain why things are important rather than using heavy-handed MUSTs.

### Test Cases

After writing the skill draft, come up with 2-3 realistic test prompts. Save test cases to `evals/evals.json`.

## Running and evaluating test cases

Spawn subagents with and without the skill for each test case. While runs are in progress, draft assertions. Grade results, aggregate into benchmark, and launch the viewer.

## Improving the skill

1. Generalize from feedback — avoid overfitting to specific examples
2. Keep the prompt lean — remove things not pulling their weight
3. Explain the why behind instructions
4. Look for repeated work across test cases to bundle as scripts

## Description Optimization

Generate 20 eval queries (mix of should-trigger and should-not-trigger), review with user, then run the optimization loop.

## Claude.ai-specific instructions

No subagents means run test cases one at a time. Skip baseline runs. Skip browser viewer — present results directly.

## Reference files

- `agents/grader.md` — Evaluating assertions
- `agents/comparator.md` — Blind A/B comparison
- `agents/analyzer.md` — Analyzing results
- `references/schemas.md` — JSON structures
