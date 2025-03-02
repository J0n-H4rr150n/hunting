# VPS

Why use a VPS?  
- Some tools take a long time to run and it is helpful if it can keep going even if you turn of your main device.  
- Some tools can essentially block all other connections while the tool runs (i.e. `amass` if not tuned properly).  
- Etc.  

Setup a new VM:  
- At least 4 GB of memory  
- At least 40 GB of disk space  

ssh into the VM:  
```
sudo apt update && sudo apt upgrade -y
```

create a new sudo user:  
```
sudo adduser bug
```

add user to sudo group:  
```
sudo usermod -aG sudo bug
```  

verify sudo access:  
```
su - bug
```
```
sudo -l
```
