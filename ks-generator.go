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
    Shadow      string
    Listen      string
    Mirror      map[string]string
}

type QueryConfig struct {
    QueryData   map[string]string
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

func Atof(s string) float64 {
    f, _ := strconv.ParseFloat(s, 64)
    return f
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
    log.Printf("request address %s", r.RemoteAddr)
    log.Printf("request uri %s", r.URL)
    if r.URL.Path != "/ks.cfg" {
        fmt.Fprintf(w, "please use /ks.cfg? to generate ks files")
    } else {
        // load config file config.json
        config := loadConfig("config.json")

        err := r.ParseForm()
        if err != nil { panic(err) }

        var qc QueryConfig
        qc.QueryData = make(map[string]string)
        for k, v := range r.Form {
            qc.QueryData[k] = v[0]  //we only allow one key coresponding to one value (a key with multivlaue will be shrink to uri.Values[k][0])
        }

        if len(r.RemoteAddr) < 12 { r.RemoteAddr = "127.0.0.1:8888" }

        ipaddr := strings.Split(r.RemoteAddr, ":")[0]
        mirror := locateMirror(ipaddr, config.Mirror)

        qc.QueryData["mirror"] = mirror
        qc.QueryData["password"] = config.Shadow

        tmpl := "ks.tmpl"
        if len(qc.QueryData["tmpl"]) > 0 {
            tmpl = qc.QueryData["tmpl"]
        }

        tp, err := ioutil.ReadFile(tmpl)
        if err != nil {
            log.Println("fail to parse template, please check if template file exists")
        }

        funcMap := template.FuncMap{
            "Atof": Atof,
        }

        t := template.Must(template.New("ks-generator").Funcs(funcMap).Parse(string(tp)))
        t.Execute(w, qc)
    }
}

func main() {
    // this will only load once and lock config.Listen value
    config := loadConfig("config.json")
    log.Println("Starting KS-Generator Web Service")
    http.HandleFunc("/", handler)
    http.ListenAndServe(config.Listen, nil)
}