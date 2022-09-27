"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
" Execute validgen to generate code for any validators found in the current file.
"
"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
function! ft#go#RunValidgenForCurrentFile()
	let project_dir = trim(system('git rev-parse --show-toplevel'))
	let current_file = expand("%:p")
	let out_file = substitute(current_file, '.go$', '_valid.go', '')

	exec ":!validgen -f=" . current_file . ' -c=' . project_dir . '/.valid.yaml'
	exec ":e " . out_file
endfunction
