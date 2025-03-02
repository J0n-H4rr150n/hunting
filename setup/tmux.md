# tmux  

# install  

```
sudo apt update && sudo apt install tmux
```

check if it works:  
```
tmux
```

# usage  

## common commands  

prefix (start of commands) -> `CTRL+b`  

`Ctrl+b` + `D` - detach from the current tmux session  

`Ctrl+b` + `%` - split window horizontally  
`Ctrl+b` + `"` - split window vertically  
`Ctrl+b` + arrow key (left right, up, date) - move around between the window panes  
`Ctrl+b` + `x` - close the current pane

`Ctrl+b` + `c` - create a new window  
`Ctrl+b` + `n` - move to next window  
`Ctrl+b` + `p` - move to previous window

`Ctrl+B` + `:` - command line with tab completion  
`Ctrl+b` + `?` - view keybindings (press q to quit)  
`ctrl+b` + `w` - navigation panel

Ctrl+B : — Enter the command line to type commands. Tab completion is available.
Ctrl+B ? — View all keybindings. Press Q to exit.
Ctrl+B W — Open a panel to navigate across windows in multiple sessions.

## example  

Start a long-running process then you can detach from it to keep it running.

prefix command:  
`Ctrl+b`  

detach command (after prefix):  
`d`  

list active tmux sessions:  
```
tmux ls
```

You should be safe to disconnect from your ssh session.  

To get back to the long-running process...  
```
tmux ls
```

view the session from the list:  
```
tmux attach -t 0
```

to cancel:  
`Ctrl+c`  

