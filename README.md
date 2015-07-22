#kickstart file generator

the url.Values will be parsed into map[string][]string[0]

when you send a querystring like
"http://127.0.0.1:8888/ks.cfg?os=centos&version=7.0&ondisk=sda&ipaddr=94.46.146.40&nm=255.255.255.0&gw=94.46.146.1&ns=8.8.8.8&hn=edge22-lhr-n.maxcdn.net&fstype=ext4&offdisk=sdb&tmpl=ks.tmpl"

you can send whatever querystring you like, as long as you remeber its name
the querystring will be in .QueryData map, so you can access it by {{index .QueryData "version"}} in the template

i also send in a function called Atof which converts string to float64 so you can compare it in template


the program will load ks.tmpl as a template file,
so you can define what ever structure file you want

ks.tmpl is the default location and serves as an example for how you create the template