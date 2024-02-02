" gpt_command.vim
function! CallExternalGPTAndInsert(arg)
    let l:command = "cgptcodegen " . shellescape(a:arg)
    let l:output = system(l:command)

    " Check if the command was successful
    if v:shell_error == 0
        " Insert the output at the current cursor position
        execute "normal! i" . l:output
    else
        echo "Error: Command failed"
    endif
endfunction

command! -nargs=1 GPT call CallExternalGPTAndInsert(<f-args>)

