Markdown to slide

[![Build Status](https://ci.ronmi.tw/api/badges/ronmi/md-slider/status.svg)](https://ci.ronmi.tw/ronmi/md-slider)

# Synopsis

```bash
docker run -d --name slider -v /path/to/your/markdown/slides:/data -p 9527:8000 ronmi/md-slider
```

After that, see http://127.0.0.1:9527

```markdown
#+TITLE: MD-Slider
#+SUBTITLE: Simple Markdown-to-slides converter inspired by x/tool/present
#+AUTHOR: Ronmi Ren
#+TITLETEXT: **Nerdy** programmer
#+TITLETEXT: Another line of title
#+EMAIL: QAQ@example.com
#+URL: https://81k.today
#+TEXT: [81K Today](https://81k.today)
#+FOOTER: [MD-Slider](https://git.ronmi.tw/ronmi/md-slider)

# First page

## Second page

some text

- some
- list
- items

# Third page

some more text

![gopher](https://blog.golang.org/gopher/header.jpg)

[Copyright information of the gopher image](https://blog.golang.org/gopher)
```

# Markup

Level 1 or 2 headings are marks of page beginning.

All lines before first page mark are metadata.

## Metadata

Required tags:

- TITLE: Title of this slide, shows in first page.
- AUTHOR: Your name.
- EMAIL: Your email, shows in last page.

Optional tags:

- FOOTER: One line of markdown text shows at bottom of every page.
- SUBTITLE: Description of slide title, show in first page.
- FACEBOOK: Your Facebook username, shows in last page.
- TWITTER: Your Twitter username (Without `@`), shows in last page.
- URL: Some url about you, shows in last page.
- TEXT: One line of markdown text about you, shows in last page.
- TITLETEXT: One line of markdown text about you, shows in first page.
