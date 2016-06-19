"contents of minimal .vimrc"
execute pathogen#infect()
syntax on
filetype plugin indent on
"to activate NERDTREE
"use this command                   -> :NERDTreeToggle
"to bookmark                        -> :Bookmark bookmarkname
"NERDTree keybinding to CTRL + n    -> : nmap <c-n> :NERDTreeToggle<cr>
nmap<c-n> :NERDTreeToggle<cr>
"TagBar keybinding                  -> :nmap <f8> :TagbarToggle<cr>
nmap<f8> :TagbarToggle<cr>
"pangloss vim-javascript
"to set the foldmethod              -> :set foldmethod=syntax
let g:javascript_enable_domhtmlcss = 1
let g:javascript_ignore_javaScriptdoc = 1
