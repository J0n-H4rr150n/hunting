# installing tools on ubuntu  

# Basic Tools  

## locate  

install:  
```
sudo apt install plocate
```

check if it works:  
```
locate -h
```  

## node  
[https://github.com/nvm-sh/nvm?tab=readme-ov-file#installing-and-updating](https://github.com/nvm-sh/nvm?tab=readme-ov-file#installing-and-updating)    

WARNING: The command below will automatically run the `install.sh` file (check it first if you want).
```
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash
```

refresh:  
```
source ~/.bashrc
```

see available node versions:  
```
nvm list-remote
```  

install a version of node:  
```
nvm install lts/jod
```  
- At the time of this writing `lts/jod` was for `node v22.14.0 (npm v10.9.2)`.  

check node version:  
```
node --version
```

check npm version:  
```
npm --version
```  

## nmap  

install:  
```
sudo apt update && sudo apt install nmap -y
```  

check if it works:  
```
nmap -h
```

## proxychains4  

install:  
```
sudo apt update && sudo apt install proxychains4 -y
```

check if it works:  
```
proxychains4 --help
```  

add tor support (install tor):  
```
sudo apt update && sudo apt install tor -y
```  

enable tor:  
```
sudo systemctl start tor
```  

check if tor is active:  
```
sudo systemctl status tor
```  

update the proxychains4 config file:  
```
sudo nano /etc/proxychains4.conf
```
- enable `dynamic_chain` (remove the comment `#`)  
- disable `strict_chain` (add a command '#' in front of it)  
- enable `random_chain`  (remove the comment `#`)
- enable `proxy_dns` (remove the comment `#`)  

Go to the very last empty row in the file and add:  
```
socks5  127.0.0.1 9050
```  

Check if it works (get your normal ip first):  
```
ip=$(curl -s https://api.ipify.org); echo "Normal ip: $ip";
```  

Check if it works (get your proxychains4 ip):  
```
ip=$(proxychains4 curl -s https://api.ipify.org); echo "proxychains4 ip: $ip";
```  

Setup tor exit node to United States:  
```
sudo nano /etc/tor/torrc
```  
- Add the following at the bottom of the file:
```
ExitNodes {us}
StrictNodes 1
GeoIPExcludeUnknown 1
AllowSingleHopCircuits 0
```  

restart the tor service:  
```
sudo systemctl restart tor
```  

check if tor is enabled:  
```
sudo systemctl status tor
```  

Check if it works again (get your proxychains4 ip):  
```
ip=$(proxychains4 curl -s https://api.ipify.org); echo "proxychains4 ip: $ip";
```
- Check the ip returned with another tool like [https://www.iplocation.net/ip-lookup](https://www.iplocation.net/ip-lookup)  
- Try to reboot if it doesn't seem to work.  

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
- Note change `bug` above to whatever username you are using.
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

## subfinder  
[https://github.com/projectdiscovery/subfinder](https://github.com/projectdiscovery/subfinder)  

install  
```
go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
```

check if it works:  
```
subfinder -h
```

## katana  
[https://github.com/projectdiscovery/katana](https://github.com/projectdiscovery/katana)  

prereq:  
```
sudo apt update
sudo snap refresh
sudo apt install zip curl wget git
sudo snap install golang --classic
wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | sudo apt-key add - 
sudo sh -c 'echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list'
sudo apt update 
sudo apt install google-chrome-stable
```  

install  
```
go install github.com/projectdiscovery/katana/cmd/katana@latest
```  

check if it works:  
```
katana -h
```  

### troubleshooting  

Might need a different version of go...  

make all of these items are installed:  
```
sudo apt-get install curl git mercurial make binutils bison gcc build-essential
```

go version manager (gvm):  
[https://github.com/moovweb/gvm](https://github.com/moovweb/gvm)  

prereq:  
```
sudo apt update && sudo apt install binutils gcc make bison -y
```  

install:  
```
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
```  

katana mentioned needing go1.18+ (but `go1.24.0` might have breaking changes):  
```
gvm install go1.22.12
```

use `go1.22`:  
```
gvm use go1.22.12
```  

check the go version:  
```
go version
```  

try to install katana again...  
```
go install github.com/projectdiscovery/katana/cmd/katana@latest
```  

check if it works:  
```
katana -h
```

# BishopFox Tools  

## jsluice  
[https://github.com/BishopFox/jsluice](https://github.com/BishopFox/jsluice)  

install  
```
go install github.com/BishopFox/jsluice/cmd/jsluice@latest
```

check if it works:  
```
jsluice -h
```
