package model

type Operation struct {
    Op        int
    Value     interface{}
    Jump      int
    FilePath  string
    Row       int
}

type Word struct {
    Type     int
    Value    interface{}
    Expanded int
}

type Token struct {
    FilePath  string
    Row       int
    TokenWord Word
}
