---
name: xlsx
description: "Use this skill any time a spreadsheet file is the primary input or output. This means any task where the user wants to: open, read, edit, or fix an existing .xlsx, .xlsm, .csv, or .tsv file; create a new spreadsheet from scratch or from other data sources; or convert between tabular file formats. The deliverable must be a spreadsheet file."
---

# Requirements for Outputs

## All Excel files

### Professional Font
Use a consistent, professional font (e.g., Arial, Times New Roman) for all deliverables unless otherwise instructed.

### Zero Formula Errors
Every Excel model MUST be delivered with ZERO formula errors (#REF!, #DIV/0!, #VALUE!, #N/A, #NAME?).

### Preserve Existing Templates
Study and EXACTLY match existing format, style, and conventions when modifying files.

## Financial models

### Color Coding Standards
- **Blue text (RGB: 0,0,255)**: Hardcoded inputs
- **Black text (RGB: 0,0,0)**: ALL formulas and calculations
- **Green text (RGB: 0,128,0)**: Links from other worksheets
- **Red text (RGB: 255,0,0)**: External links to other files
- **Yellow background (RGB: 255,255,0)**: Key assumptions

### Formula Construction Rules
Place ALL assumptions in separate assumption cells. Use cell references instead of hardcoded values in formulas.

# XLSX creation, editing, and analysis

## CRITICAL: Use Formulas, Not Hardcoded Values

Always use Excel formulas instead of calculating values in Python and hardcoding them.

### ❌ WRONG
```python
sheet['B10'] = total  # Hardcodes 5000
```

### ✅ CORRECT
```python
sheet['B10'] = '=SUM(B2:B9)'
```

## Common Workflow
1. Choose tool: pandas for data, openpyxl for formulas/formatting
2. Create/Load workbook
3. Modify data, formulas, and formatting
4. Save to file
5. Recalculate formulas: `python scripts/recalc.py output.xlsx`
6. Verify and fix any errors

### Creating new Excel files

```python
from openpyxl import Workbook
from openpyxl.styles import Font, PatternFill, Alignment

wb = Workbook()
sheet = wb.active
sheet['A1'] = 'Hello'
sheet['B1'] = 'World'
sheet['B2'] = '=SUM(A1:A10)'
sheet['A1'].font = Font(bold=True, color='FF0000')
sheet['A1'].fill = PatternFill('solid', start_color='FFFF00')
sheet.column_dimensions['A'].width = 20
wb.save('output.xlsx')
```

### Editing existing Excel files

```python
from openpyxl import load_workbook

wb = load_workbook('existing.xlsx')
sheet = wb.active
sheet['A1'] = 'New Value'
sheet.insert_rows(2)
sheet.delete_cols(3)
wb.save('modified.xlsx')
```

## Recalculating formulas

```bash
python scripts/recalc.py output.xlsx 30
```

The script recalculates all formulas, scans for errors, and returns JSON with detailed error locations.

## Formula Verification Checklist
- [ ] Test 2-3 sample references before building full model
- [ ] Confirm Excel columns match (e.g., column 64 = BL, not BK)
- [ ] Remember Excel rows are 1-indexed
- [ ] Check for null values with `pd.notna()`
- [ ] Verify all cell references point to intended cells
- [ ] Check denominators before using `/` in formulas

## Best Practices

### Library Selection
- **pandas**: Best for data analysis, bulk operations, simple data export
- **openpyxl**: Best for complex formatting, formulas, Excel-specific features

### Working with openpyxl
- Cell indices are 1-based (row=1, column=1 = A1)
- Use `data_only=True` to read calculated values
- If saved with `data_only=True`, formulas are replaced with values permanently
- Use `read_only=True` for large files

### Working with pandas
```python
import pandas as pd
df = pd.read_excel('file.xlsx')
df.to_excel('output.xlsx', index=False)
```
