# init  

# Basic Tools  

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


# ProjectDiscovery Tools
