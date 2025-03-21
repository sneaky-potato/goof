# Control flow in goof

## if elif else

```pascal
5 0
if 2dup > do
    "hello world" puts
elif 2dup = do
    "goofbye world" puts
else
    "else world" puts
end
```

```mermaid
stateDiagram
    state "2dup >" as ifCondition
    state "do" as ifDo
    state "hello world" as ifBody
    state "2dup =" as elifCondition
    state "do" as elifDo
    state "goofbye world" as elifBody
    state "else world" as elseBody
    [*] --> if
    if --> ifCondition
    ifCondition --> ifDo
    ifDo --> elifCondition: Jump
    ifDo --> ifBody
    ifBody --> elif
    elifCondition --> elifDo
    elifDo --> elseBody: Jump
    elifDo --> elifBody
    elifBody --> else
    elif --> end
    else --> end
    elseBody --> end
    end --> [*]
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
    state "2dup >" as whileCondition
    state "hello world" as whileBody
    [*] --> while
    while --> whileCondition
    whileCondition --> do
    do --> [*]: Jump
    do --> whileBody
    whileBody --> end
    end --> while
```

