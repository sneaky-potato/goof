# Control flow in goof

## if else

```pascal
5 0
2dup > if
    "hello world" puts
else
    "goofbye world" puts
end
```

```mermaid
stateDiagram
    [*] --> if
    if --> else+1: Jump
    if --> body
    body --> else
    else --> else+1
    else+1 --> end
    else --> end: Jump
    end --> end+1: Jump
```

## If elif else

```pascal
5 0
if 2dup > do
    "hello world" puts
elif 2dup = do
    "hello different world" puts
else
    "goofbye world" puts
end
```

```mermaid
stateDiagram
    [*] --> if
    if --> do1
    do1 --> elif+1: Jump
    elif --> do2
    elif --> else+1: Jump
    do2 --> else: Jump
```

## Loops

```pascal
5 0
while 2dup > do
    "hello world" puts
    1 +
end
```

```mermaid
stateDiagram
    [*] --> while
    while --> do
    do --> while: Jump
    end --> do: Jump
    do --> end: true
    do --> end+1: false
```
