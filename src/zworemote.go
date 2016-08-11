// Copyright 2016 Ted Goddard. All rights reserved.
// Use of this source code is governed the MIT
// license that can be found in the LICENSE file.

package main

import "fmt"
import "flag"
//import "net"
//import "io/ioutil"
//import "os"
//import "os/signal"
//import "os/exec"
//import "path/filepath"
//import "strings"
//import "bufio"
import "net/http"
import "log"
import "zwoasi"
import "zworemote"
//import "encoding/json"

func main() {
    fmt.Println("Starting server")
//    cameraNumber := flag.String("cam", "", "device number of camera")
    flag.Parse()

    http.HandleFunc("/zworemote/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, zworemote.ClientHTML)
    })

    http.HandleFunc("/zworemote/img.png", func(w http.ResponseWriter, r *http.Request) {
//        object := r.FormValue("o")
        w.Header().Set("Content-Type", "image/png")
        zwoasi.WriteImage(640, 480, w)

    })

    log.Print("http.ListenAndServe")
    log.Fatal(http.ListenAndServe(":8080", nil))


    zwoasi.GetImage("pleides")
}
