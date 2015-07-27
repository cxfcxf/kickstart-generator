package main

import (
    "log"
    "fmt"
    "flag"
    "strconv"
    "net/http"
    "io/ioutil"
    "text/template"
)

var listen = flag.String("l", ":8888", "listen ip:port, default is 0.0.0.0:8888")

type QueryConfig struct {
    QueryData   map[string]string
}

func Atof(s string) float64 {
    f, _ := strconv.ParseFloat(s, 64)
    return f
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.Printf("request address %s", r.RemoteAddr)
    log.Printf("request uri %s", r.URL)
    if r.URL.Path != "/ks.cfg" {
        fmt.Fprintf(w, "please use /ks.cfg? to generate ks files")
    } else {

        err := r.ParseForm()
        if err != nil { panic(err) }

        var qc QueryConfig
        qc.QueryData = make(map[string]string)
        for k, v := range r.Form {
            qc.QueryData[k] = v[0]  //we only allow one key coresponding to one value (a key with multivlaue will be shrink to uri.Values[k][0])
        }

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
    flag.Parse()
    log.Println("Starting KS-Generator Web Service")
    http.HandleFunc("/", handler)
    http.ListenAndServe(*listen, nil)
}