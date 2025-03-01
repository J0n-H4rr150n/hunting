# wordlists  

Setup a folder similar to kali:  
```
sudo mkdir /usr/share/wordlists
```  

## SecLists  

git the data:  
```
git clone https://github.com/danielmiessler/SecLists.git
```

move `SecLists` to the wordlists folder:  
```
sudo mv SecLists/ /usr/share/wordlists
```  

go to:
```
/usr/share/wordlists/SecLists/Passwords/Leaked-Databases
```  

unzip the rockyou file:  
```
sudo tar -xzf rockyou.txt.tar.gz
```  

check the `rockyou.txt` file extracted from above:  
```
cat rockyou.txt | wc -l
```  

## Assetnote Wordlists  

Select the wordlists you want from:  
[https://wordlists.assetnote.io/](https://wordlists.assetnote.io/)  

