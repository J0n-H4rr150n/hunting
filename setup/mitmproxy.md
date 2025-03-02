# mitmproxy  

## Install  

install pipx:  
```
sudo apt update && sudo apt install pipx -y
```

install on ubuntu with pipx:  
```
pipx install mitmproxy
```  

ensurepath:  
```
pipx ensurepath
```  

check if it works:  
```
mitmproxy
```  

run without checking certs and save to the cloud  
```
mitmproxy --ssl-insecure -s save_to_cloud.py
```  
- Example script at: [https://github.com/J0n-H4rr150n/hunting/blob/main/scripts/mitmproxy_addons/save_to_cloud.py](https://github.com/J0n-H4rr150n/hunting/blob/main/scripts/mitmproxy_addons/save_to_cloud.py)  

## Links  

User Interface (cli) docs:  
[https://docs.mitmproxy.org/stable/mitmproxytutorial-userinterface/](https://docs.mitmproxy.org/stable/mitmproxytutorial-userinterface/)  

Example addons:  
[https://docs.mitmproxy.org/stable/addons-examples/](https://docs.mitmproxy.org/stable/addons-examples/)  

