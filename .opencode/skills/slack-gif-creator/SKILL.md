---
name: slack-gif-creator
description: Knowledge and utilities for creating animated GIFs optimized for Slack. Provides constraints, validation tools, and animation concepts. Use when users request animated GIFs for Slack like "make me a GIF of X doing Y for Slack."
---

# Slack GIF Creator

A toolkit providing utilities and knowledge for creating animated GIFs optimized for Slack.

## Slack Requirements

**Dimensions:**
- Emoji GIFs: 128x128 (recommended)
- Message GIFs: 480x480

**Parameters:**
- FPS: 10-30 (lower is smaller file size)
- Colors: 48-128 (fewer = smaller file size)
- Duration: Keep under 3 seconds for emoji GIFs

## Core Workflow

```python
from core.gif_builder import GIFBuilder
from PIL import Image, ImageDraw

builder = GIFBuilder(width=128, height=128, fps=10)

for i in range(12):
    frame = Image.new('RGB', (128, 128), (240, 248, 255))
    draw = ImageDraw.Draw(frame)
    # Draw your animation using PIL primitives
    builder.add_frame(frame)

builder.save('output.gif', num_colors=48, optimize_for_emoji=True)
```

## Drawing Graphics

### Working with User-Uploaded Images
If a user uploads an image, consider whether they want to use it directly or as inspiration.

### Drawing from Scratch
Use PIL ImageDraw primitives:

```python
draw = ImageDraw.Draw(frame)
draw.ellipse([x1, y1, x2, y2], fill=(r, g, b), outline=(r, g, b), width=3)
draw.polygon(points, fill=(r, g, b), outline=(r, g, b), width=3)
draw.line([(x1, y1), (x2, y2)], fill=(r, g, b), width=5)
```

**Don't use:** Emoji fonts (unreliable across platforms).

### Making Graphics Look Good
- Use thicker lines (width=2 or higher)
- Add visual depth with gradients and layered shapes
- Use vibrant, complementary colors
- Be creative and detailed

## Available Utilities

### GIFBuilder (`core.gif_builder`)
```python
builder = GIFBuilder(width=128, height=128, fps=10)
builder.add_frame(frame)
builder.add_frames(frames)
builder.save('out.gif', num_colors=48, optimize_for_emoji=True, remove_duplicates=True)
```

### Validators (`core.validators`)
```python
from core.validators import validate_gif, is_slack_ready
passes, info = validate_gif('my.gif', is_emoji=True, verbose=True)
```

### Easing Functions (`core.easing`)
```python
from core.easing import interpolate
y = interpolate(start=0, end=400, t=t, easing='ease_out')
# Available: linear, ease_in, ease_out, ease_in_out, bounce_out, elastic_out, back_out
```

### Frame Helpers (`core.frame_composer`)
```python
from core.frame_composer import (
    create_blank_frame, create_gradient_background, draw_circle, draw_text, draw_star
)
```

## Animation Concepts

- **Shake/Vibrate**: Offset position with oscillation using `math.sin()`
- **Pulse/Heartbeat**: Scale size rhythmically with sine wave
- **Bounce**: Use `bounce_out` easing for landing
- **Spin/Rotate**: `image.rotate(angle, resample=Image.BICUBIC)`
- **Fade In/Out**: Use RGBA images or `Image.blend()`
- **Slide**: Move from off-screen to position with `ease_out`
- **Zoom**: Scale and crop center
- **Explode/Particle Burst**: Create particles with random velocities

## Optimization Strategies

1. Fewer frames - Lower FPS or shorter duration
2. Fewer colors - `num_colors=48`
3. Smaller dimensions - 128x128
4. Remove duplicates - `remove_duplicates=True`
5. Emoji mode - `optimize_for_emoji=True`

## Dependencies
```bash
pip install pillow imageio numpy
```
