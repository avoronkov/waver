" This file was generated by syntaxgen.

" Quit when a syntax file was already loaded.
if exists('b:current_syntax') | finish|  endif

" Keywords
syn keyword celPragmaSt define

syn keyword celStdFunction - dup not and or + * drop top stack
syn keyword celFunction Play PlayBack FF Pos Len NPlay Goto

syn match celNumber '\<\d\+\>'
syn match celNumber '\<\d\+\.\d\+\>'

syn match celForbiddedTab '\t'

syn keyword celTodo contained TODO FIXME XXX NOTE
syn match celComment "#.*$" contains=celTodo

syn match celString '".*"'

" Define highlighting
hi def link celPragmaSt Special
hi def link celStdFunction Operator
hi def link celFunction Identifier
hi def link celNumber Number
hi def link celString String

hi def link celForbiddedTab Error

hi def link celTodo Todo
hi def link celComment Comment

let b:current_syntax = 'surfer'
