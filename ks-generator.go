package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func locateMirror(ip []string) string {
	network := fmt.Sprintf("%s.%s.%s", ip[0], ip[1], ip[2])

	f, err := ioutil.ReadFile("./mirrors.json")
	if err != nil { panic(err) }

	var config map[string]string
	json.Unmarshal(f, &config)

	for k, v := range config {
		if network == k {
			return v
		}
	}
	if len(config["default"]) > 0 {
		return config["default"]
	} else {
		return "http://mirror.centos.org/centos-6/"
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.RemoteAddr, r.URL)
	if r.URL.Path != "/ks.cfg" {
		fmt.Fprintf(w, "please use /ks.cfg? to generate ks files")
	} else {
		shadows := "your password shadows"

		err := r.ParseForm()
		if err != nil { panic(err) }
		uri := r.Form

		if len(r.RemoteAddr) < 12 { r.RemoteAddr = "127.0.0.1:8888" }

		ip := strings.Split(r.RemoteAddr, ".")
		ipaddr := strings.Split(r.RemoteAddr, ":")[0]
		nm := "255.255.255.0"
		gw := fmt.Sprintf("%s.%s.%s.1", ip[0], ip[1], ip[2])
		mirror := locateMirror(ip)

		version, fstype, ondisk, dorder := "6.5", "ext4", "", ""
		if uri.Get("version") != "" {version = uri.Get("version")}
		if uri.Get("fstype") != "" {fstype = uri.Get("fstype")}
		if uri.Get("ondisk") != "" {
			ondisk = "--ondisk=" + uri.Get("ondisk")
			dorder = " --driveorder=" + uri.Get("ondisk")
		}
		
		if uri.Get("ipaddr") != "" {ipaddr = uri.Get("ipaddr")}
		if uri.Get("nm") != "" {nm = uri.Get("nm")}
		
		// if you specify ip you must specify gw also, otherwise it will be 127.0.0.1
		if uri.Get("gw") != "" {gw = uri.Get("gw")}

		fmt.Fprintf(w, "# kickstart generator\n")
		fmt.Fprintf(w, "install\n")
		fmt.Fprintf(w, "url --url=%s%s/os/x86_64\n\n", mirror, version)
		fmt.Fprintf(w, "rootpw --iscrypted %s\n\n", shadows)
		fmt.Fprintf(w, "network --onboot yes --device eth0 --mtu=1500 --bootproto static --ip %s --netmask %s --gateway %s --nameserver 8.8.8.8\n", ipaddr, nm, gw)
		fmt.Fprintf(w, "auth --useshadow --passalgo=sha512 --enablefingerprint\ntext\nkeyboard us\nlang en_US\nselinux --disabled\nfirewall --disabled\nskipx\nlogging --level=info\ntimezone --utc Etc/GMT\n\nzerombr\nclearpart --all\n\n")
		fmt.Fprintf(w, "part / --fstype=%s --size=1 --grow %s --asprimary\n", fstype, ondisk)
		fmt.Fprintf(w, "part swap --fstype=swap --size=4096 %s\n", ondisk)
		fmt.Fprintf(w, "bootloader --location=mbr%s --append=\"crashkernel=auto rhgb quiet\"\n", dorder)
		fmt.Fprintf(w, "firstboot --disable\nreboot\n\n")
		fmt.Fprintf(w, "%%pre\n#/bin/sh\ntouch /tmp/part.cfg\n%%end\n\n")
		fmt.Fprintf(w, "%%packages\n@base\n@core\n\n")
		fmt.Fprintf(w, "%%post\n/usr/sbin/pwconv\n%%end")
	}
}

func main() {
	log.Println("Starting KS-Generator Web Service")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8888", nil)
}