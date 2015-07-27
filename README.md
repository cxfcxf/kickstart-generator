# kickstart file generator

i call it kickstart generator since i mainly use it to generate kickstart file and provide them for pxebooting.

Due to the feature it has, the program is more like quertstring to web text convertor.

### Notice 
the url.Values will be parsed into map[string][]string[0] which means:
```
ks.cfg?os=centos&os=ubuntu
```
will result in:
```
.QueryData["os"] = centos
```
the other multipart for same key will be stirpped out

# How to use it:
when you send a querystring to browser
```
"http://127.0.0.1:8888/ks.cfg?version=7&ondisk=sda&ipaddr=94.46.146.40&nm=255.255.255.0&gw=94.46.146.1&ns=8.8.8.8&hn=www.ksgenerator.come&fstype=ext4&offdisk=sdb&tmpl=ks.tmpl"
```
you will get a webtext in your browser depending on your self-defined template

#### in this case, its ks.tmpl (ks.tmpl will be the default template file)

you can send whatever querystring you like, as long as you remeber its name

the querystring can be accessed by  {{.QueryData}} map inside template you define.

### for example
you can access it by
```
{{index .QueryData "key"}}
```
in the ks.tmpl

### Notice:
the defuat template is ks.tmpl, but you can use multiple template for different purposes as long as you have
```
&tmpl=ks.tmpl inside your querystring
```
### Also:
the template can be dynamicly loaded since the program load it for each request.

you can define whatever structure you want template as long as you remember the querystring you passed

### additional functions can be used in template
i also send in a function called Atof Ato converts string to float64

so you can compare float64 number inside template which is handy for version comparson


### Config.json
removed not necessary

# License
Mit