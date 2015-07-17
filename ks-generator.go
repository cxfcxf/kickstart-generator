package main

import (
    "os"
    "log"
    "fmt"
    "strconv"
    "strings"
    "net/http"
    "io/ioutil"
    "text/template"
    "encoding/json"
)

type Config struct {
    Template    string
    Shadow      string
    Resolver    string
    Mirror      map[string]string
}

type QueryConfig struct {
    Version     float64
    Mirror      string
    Password    string
    Ipaddr      string
    Netmask     string
    Gateway     string
    Nameserver  string
    Hostname    string
    Fstype      string
    Ondisk      string
    Offdisk     string
}

func loadConfig(file string) Config{
    f, err := ioutil.ReadFile(file)
    if err != nil {
        log.Println("failed to read config.json")
        os.Exit(1)
    }

    var config Config

    err = json.Unmarshal(f, &config)
    if err != nil {
        log.Println("failed to Parse config")
        os.Exit(1)
    }
    return config
}


func locateMirror(ipaddr string, mirror map[string]string) string {

    for k, v := range mirror {
        if strings.Contains(ipaddr, k) {
            return v
        }
    }
    if len(mirror["default"]) > 0 {
        return mirror["default"]
    } else {
        return "http://mirror.centos.org/centos-6/"
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.Printf("%s %s", r.RemoteAddr, r.URL)
    if r.URL.Path != "/ks.cfg" {
        fmt.Fprintf(w, "please use /ks.cfg? to generate ks files")
    } else {
        // load config file config.json
        config := loadConfig("config.json")

        err := r.ParseForm()
        if err != nil { panic(err) }
        uri := r.Form

        if len(r.RemoteAddr) < 12 { r.RemoteAddr = "127.0.0.1:8888" }

        ipaddr := strings.Split(r.RemoteAddr, ":")[0]
        netmask := "255.255.255.0"
        ip := strings.Split(r.RemoteAddr, ".")
        gateway := fmt.Sprintf("%s.%s.%s.1", ip[0], ip[1], ip[2])
        mirror := locateMirror(ipaddr, config.Mirror)

        version, fstype, ondisk, offdisk, hostname := 6.6, "ext4", "", "", ""
        if uri.Get("version") != "" {version, _ = strconv.ParseFloat(uri.Get("version"), 64)}
        if uri.Get("fstype") != "" {fstype = uri.Get("fstype")}
        if uri.Get("ondisk") != "" {ondisk = uri.Get("ondisk")}
        if uri.Get("offdisk") != "" {offdisk = uri.Get("offdisk")}
        if uri.Get("ipaddr") != "" {ipaddr = uri.Get("ipaddr")}
        if uri.Get("nm") != "" {netmask = uri.Get("nm")}
        
        // if you specify ip you must specify gw also, otherwise it will be 127.0.0.1
        if uri.Get("gw") != "" {gateway = uri.Get("gw")}
        if uri.Get("hostname") != "" {hostname = uri.Get("hostname")}

        fmt.Println(version)
        qc := &QueryConfig{
            version,
            mirror,
            config.Shadow,
            ipaddr,
            netmask,
            gateway,
            config.Resolver,
            hostname,
            fstype,
            ondisk,
            offdisk,
        }

        tp, _ := ioutil.ReadFile(config.Template)

        t := template.Must(template.New("ks-generator").Parse(string(tp)))
        t.Execute(w, qc)
    }
}

func main() {
    log.Println("Starting KS-Generator Web Service")
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8888", nil)
}