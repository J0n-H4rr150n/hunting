# metasploit  
[https://github.com/rapid7/metasploit-framework](https://github.com/rapid7/metasploit-framework)  

Using metasploit
[https://docs.metasploit.com/docs/using-metasploit/basics/using-metasploit.html](https://docs.metasploit.com/docs/using-metasploit/basics/using-metasploit.html)  

install on ubuntu:  
```
curl https://raw.githubusercontent.com/rapid7/metasploit-omnibus/master/config/templates/metasploit-framework-wrappers/msfupdate.erb > msfinstall && \
  chmod 755 msfinstall && \
  ./msfinstall
```

check if it works:  
```
msfconsole
```
