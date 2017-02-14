# gomiruleaker cli

download, parse, and post emails to miru-leaks

## Dev Notes

Uses concurrent pipeline pattern

```generator -> parse[] -> accumulator -> poster```
