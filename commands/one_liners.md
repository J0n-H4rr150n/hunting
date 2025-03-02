# one-line commands

## recon  

### assetfinder + anew  
```
cat domains.txt | assetfinder | sed 's/^\*\.//' | sort -u | anew assetfinder_domains.txt
```  
- `domains.txt` is a list of target domains to start with (one per line)  
- The `sed` command remove `*.` for wildcards  
- Make sure to have the necessary tools installed - [https://github.com/J0n-H4rr150n/hunting/blob/main/setup/tools.md](https://github.com/J0n-H4rr150n/hunting/blob/main/setup/tools.md)  

check if the domains are potentially alive on ports 80 and 443:  
```
xargs -I {} sh -c 'nc -zv -w 5 {} 80 && echo {} >> assetfinder_domains_alive.txt || nc -zv -w 5 {} 443 && echo {} >> assetfinder_domains_alive.txt' < assetfinder_domains.txt
```

a different check to see if the hostnames resolve to an ip:  
```
args -I {} sh -c 'if host {} > /dev/null; then echo {} >> assetfinder_domains_with_ips.txt; fi' < assetfinder_domains.txt
```  

try to detect if the hostname is behind a `waf`:  
```
xargs -I {} sh -c 'wafw00f -a {} | anew assetfinder_domains_alive_waf.txt' < assetfinder_domains_alive.txt
```  
- Make sure to have the necessary tools installed - [https://github.com/J0n-H4rr150n/hunting/tree/main/tools](https://github.com/J0n-H4rr150n/hunting/tree/main/tools)  
