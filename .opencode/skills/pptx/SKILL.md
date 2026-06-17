---
name: pptx
description: "Use this skill any time a .pptx file is involved in any way — as input, output, or both. This includes: creating slide decks, pitch decks, or presentations; reading, parsing, or extracting text from any .pptx file; editing, modifying, or updating existing presentations; combining or splitting slide files; working with templates, layouts, speaker notes, or comments. Trigger whenever the user mentions \"deck,\" \"slides,\" \"presentation,\" or references a .pptx filename."
---

# PPTX Skill

## Quick Reference

| Task | Guide |
|------|-------|
| Read/analyze content | `python -m markitdown presentation.pptx` |
| Edit or create from template | Read `editing.md` |
| Create from scratch | Read `pptxgenjs.md` |

---

## Reading Content

```bash
python -m markitdown presentation.pptx
python scripts/thumbnail.py presentation.pptx
python scripts/office/unpack.py presentation.pptx unpacked/
```

---

## Editing Workflow

Read `editing.md` for full details.

1. Analyze template with `thumbnail.py`
2. Unpack → manipulate slides → edit content → clean → pack

---

## Creating from Scratch

Read `pptxgenjs.md` for full details.

Use when no template or reference presentation is available.

---

## Design Ideas

Don't create boring slides. Consider ideas from this list.

### Before Starting
- Pick a bold, content-informed color palette
- One color should dominate (60-70%), with 1-2 supporting tones and one sharp accent
- Consider dark/light contrast: dark for title + conclusion, light for content
- Commit to a visual motif and carry it across every slide

### Color Palettes

| Theme | Primary | Secondary | Accent |
|-------|---------|-----------|--------|
| Midnight Executive | `1E2761` (navy) | `CADCFC` (ice blue) | `FFFFFF` |
| Forest & Moss | `2C5F2D` (forest) | `97BC62` (moss) | `F5F5F5` |
| Coral Energy | `F96167` (coral) | `F9E795` (gold) | `2F3C7E` |
| Warm Terracotta | `B85042` | `E7E8D1` (sand) | `A7BEAE` |
| Charcoal Minimal | `36454F` | `F2F2F2` | `212121` |

### For Each Slide
**Every slide needs a visual element.** Layout options: two-column, icon+text rows, 2x2 grid, half-bleed image. Data display: large stat callouts, comparison columns, timeline.

### Typography

| Header Font | Body Font |
|-------------|-----------|
| Georgia | Calibri |
| Arial Black | Arial |
| Cambria | Calibri |

| Element | Size |
|---------|------|
| Slide title | 36-44pt bold |
| Body text | 14-16pt |
| Captions | 10-12pt |

### Avoid
- Don't repeat the same layout — vary columns, cards, callouts
- Don't center body text — left-align paragraphs and lists
- Don't default to blue — pick colors for the specific topic
- Don't create text-only slides
- NEVER use accent lines under titles (hallmark of AI-generated slides)

---

## QA (Required)

**Assume there are problems. Your job is to find them.**

### Content QA
```bash
python -m markitdown output.pptx
grep -iE "xxxx|lorem|ipsum|placeholder" output.md
```

### Visual QA
Convert slides to images, then use subagents to inspect for:
- Overlapping elements
- Text overflow
- Low-contrast text/icons
- Uneven spacing
- Elements too close together

### Verification Loop
1. Generate slides → Convert to images → Inspect
2. List issues found
3. Fix issues
4. Re-verify affected slides
5. Repeat until no new issues

---

## Converting to Images

```bash
python scripts/office/soffice.py --headless --convert-to pdf output.pptx
pdftoppm -jpeg -r 150 output.pdf slide
```

---

## Dependencies

- `pip install "markitdown[pptx]"` - text extraction
- `pip install Pillow` - thumbnail grids
- `npm install -g pptxgenjs` - creating from scratch
- LibreOffice (`soffice`) - PDF conversion
- Poppler (`pdftoppm`) - PDF to images
