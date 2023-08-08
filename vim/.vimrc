set clipboard=unnamedplus
set noerrorbells
set visualbell
set smartcase
set number
set inccommand=split
set scrolloff=5
set incsearch

"""""" Mappings
let mapleader=" "

"" Editing
inoremap jk <ESC>
vnoremap < <gv
vnoremap > >gv
nnoremap <A-j> :m .+1<CR>==
vnoremap <A-j> :m .+1<CR>==
nnoremap <A-k> :m .-2<CR>==
vnoremap <A-k> :m .-2<CR>==

xnoremap J :move '>+1<CR>gv
xnoremap K :move '<-2<CR>gv
xnoremap <A-j> :move '>+1<CR>gv-gvx
xnoremap <A-K> :move '<-2<CR>gv-gv
inoremap <A-j> <Esc>:m .+1<CR>==gi
inoremap <A-K> <Esc>:m .-2<CR>==gi

"" Buffers and Windows
" Close all buffers except current
nnoremap <Leader>bo <cmd>%bd|e#<cr>
nnoremap <C-h> <C-w>h
nnoremap <C-l> <C-w>l
nnoremap <C-j> <C-w>j
nnoremap <C-k> <C-w>k

"" Navigation
nnoremap <S-l> :bnext<CR>
nnoremap <S-h> :bprevious<CR>
nnoremap <C-x> :bdelete<CR>
