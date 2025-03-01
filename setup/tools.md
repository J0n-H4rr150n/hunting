# installing tools on ubuntu  

# Basic Tools  

## nmap  

install:  
```
sudo apt update && sudo apt install nmap -y
```  

## golang  

Download go for linux:  
[https://go.dev/dl/](https://go.dev/dl/)  

Example (will change over time):  
```
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
```   

Unzip and put in /usr/local:  
```
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
```  

add to path:  
```
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

source the updated file:  
```
source ~/.bashrc
```

check if go works:  
```
go version
```

add the golang tool install location to path:  
```
echo 'export PATH=$PATH:/home/bug/go/bin' >> ~/.bashrc
```
- If the above doesn't work create the directory or wait until after installing a golang tool.

# TomNomNom Tools  

## anew  
[https://github.com/tomnomnom/anew](https://github.com/tomnomnom/anew)  

install  
```
go install -v github.com/tomnomnom/anew@latest
```  

check if it works:  
```
anew -h
```

## waybackurls  
[https://github.com/tomnomnom/waybackurls](https://github.com/tomnomnom/waybackurls)  

install  
```
go install github.com/tomnomnom/waybackurls@latest
```

check if it works:  
```
waybackurls -h
```  

## assetfinder  
[https://github.com/tomnomnom/assetfinder](https://github.com/tomnomnom/assetfinder)    

install
```
go install github.com/tomnomnom/assetfinder@latest
```

check if it works:  
```
assetfinder -h
```  

## httprobe  
[https://github.com/tomnomnom/httprobe](https://github.com/tomnomnom/httprobe)  

install  
```
go install github.com/tomnomnom/httprobe@latest
```

check if it works:  
```
httprobe -h
```  

## gron  
[https://github.com/tomnomnom/gron](https://github.com/tomnomnom/gron)  

install  
```
go install github.com/tomnomnom/gron@latest
```

check if it works:  
```
gron -h
```  

# ProjectDiscovery Tools  

## httpx  
[https://github.com/projectdiscovery/httpx](https://github.com/projectdiscovery/httpx)  

install
```
go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest
```

check if it works:  
```
httpx -h
```






