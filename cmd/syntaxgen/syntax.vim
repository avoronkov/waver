" This file was generated by syntaxgen.

" Quit when a syntax file was already loaded.
if exists('b:current_syntax') | finish|  endif

" Keywords
syn keyword celPragmaSt{{ range .Pragmas }} {{ . }}{{ end }}

syn keyword celFunction{{ range .StdFunctions }} {{ . }}{{ end }}
syn keyword celRepeat{{ range .Functions }} {{ . }}{{ end }}
{{ range .FunctionOperators }}
syn match celRepeat "{{.}}"{{ end }}

{{ range .ModifierOperators }}
syn match celOperator '{{.}}'{{ end }}

{{ range .Modifiers }}
syn match celOperator "\<{{ . }}\>"{{ end }}

syn match celSignal '->'

syn match celPragma '%'
syn match celPragma '%%'

{{ range .Identifiers }}
syn match celIdent '\<{{.}}\>'{{ end }}

syn match celString "'.*'"
syn match celString '".*"'

syn match celNumber '\<\d\+\>'
syn match celNumber '\<\d\+\.\d\+\>'
syn match celNumber '\<[ABCDEFG][bs]\?\d\>'

syn match celForbiddedTab '\t'

syn keyword celTodo contained TODO FIXME XXX NOTE
syn match celComment "#.*$" contains=celTodo

" instruments options
syn keyword celFilter{{ range .Filters }} {{.}}{{ end }}
{{ range .FilterOptions }}
syn keyword celFilterOption {{.}}{{ end }}

" Define highlighting
hi def link celPragmaSt Conditional

hi def link celNumber Number
hi def link celString String
hi def link celFunction Function
hi def link celRepeat Function
hi def link celOperator Operator
hi def link celIdent Identifier
hi def link celTodo Todo
hi def link celComment Comment
hi def link celPragma Type
hi def link celSignal Special
hi def link celForbiddedTab Error

hi def link celFilter StorageClass
hi def link celFilterOption Structure

let b:current_syntax = 'pelia'