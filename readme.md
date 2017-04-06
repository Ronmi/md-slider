Markdown to slide

# Synopsis

```bash
docker run -d --name slider -v /path/to/your/markdown/slides:/data -p 9527:8000 ronmi/md-slider
```

After that, see http://127.0.0.1:9527

# 支援的 markup

- headings: `#`/`##`/`###`
- lists: `*`/`-`/`+`/`1.`
- code block: triple back-quote with hint
- formats: 
  - fixed: back-quote
  - bold: `*`
  - italic: `_`
  - strike: `~~`
  - link: `[text](url)`

## 有特殊處理的 hint

- `html`

## 會特殊處理的 markup

`#` 和 `##` 都會開一張新的投影片

投影片沒有內文的時候，標題會上下左右置中
