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

Start a long-running process then you can detach from it to keep it running.

prefix command:  
`Ctrl+B`  

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
`CTRL+C`  
