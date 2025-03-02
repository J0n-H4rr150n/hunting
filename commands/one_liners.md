# one-line commands

## recon  

### assetfinder + anew  
```
cat domains.txt | assetfinder | sed 's/^\*\.//' | sort -u | anew assetfinder_domains.txt
```  
- `domains.txt` is a list of target domains to start with (one per line)  
- The `sed` command remove `*.` for wildcards  
- Make sure to have the necessary tools installed - [https://github.com/J0n-H4rr150n/hunting/blob/main/setup/tools.md](https://github.com/J0n-H4rr150n/hunting/blob/main/setup/tools.md)  
