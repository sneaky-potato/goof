" Vim syntax file
" Language: Goof

" Usage
" Put this file in 
" ~/.vim/syntax/goof.vim OR ~/.config/nvim/syntax/goof.vim
" and add in your .vimrc file the following line
" autocmd BufRead,BufNewFile *.goof set filetype=goof

if exists("b:current_syntax")
  finish
endif

set iskeyword=a-z,A-Z,-,*,_,!,@
syntax keyword goofTodos TODO XXX FIXME NOTE

" Language keywords
syntax keyword goofKeywords if if* else while do include mem end here macro

syntax keyword goofStatements dup rot over swap

" Comments
syntax region goofCommentLine start="//" end="$"   contains=goofTodos

" String literals
syntax region goofString start=/\v"/ skip=/\v\\./ end=/\v"/ contains=goofEscapes

" Char literals
syntax region goofChar start=/\v'/ skip=/\v\\./ end=/\v'/ contains=goofEscapes

" Escape literals \n, \r, ....
syntax match goofEscapes display contained "\\[nr\"']"

" Number literals
syntax region goofNumber start=/\s\d/ skip=/\d/ end=/\s/

" Set highlights
highlight default link goofTodos Todo
highlight default link goofKeywords Keyword
highlight default link goofStatements Statement
highlight default link goofCommentLine Comment
highlight default link goofString String
highlight default link goofNumber Number
highlight default link goofChar Character
highlight default link goofEscapes SpecialChar

let b:current_syntax = "goof"
