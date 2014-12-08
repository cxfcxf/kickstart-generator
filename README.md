#kickstart file generator

in your pxelinux menu you can do something like

```
append ks=http://addressofthisserver/ks.cfg?version=6.5&ondisk=sdb
```

only support couple of variable right now,
version, ondisk, ip, gw, nm,

the mirrorlist maps your network into different mirrors which you can define yourself.

