# VPS

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
