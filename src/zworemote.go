// Copyright 2016 Ted Goddard. All rights reserved.
// Use of this source code is governed the MIT
// license that can be found in the LICENSE file.

package main

import "fmt"
import "flag"
//import "net"
//import "io/ioutil"
import "io"
import "image"
import "os"
import "os/signal"
//import "os/exec"
//import "path/filepath"
//import "strings"
//import "bufio"
import "strconv"
import "net/http"
import "log"
import "zwoasi"
import "zworemote"
//import "encoding/json"

func main() {
    fmt.Println("Starting server")
//    cameraNumber := flag.String("cam", "", "device number of camera")
    flag.Parse()

    zwoasi.OpenCamera()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt)
    go func(){
        for _ = range sigChan {
            zwoasi.CloseCamera()
            fmt.Println("Camera closed.")
            os.Exit(0)
        }
    }()

    http.HandleFunc("/zworemote/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, zworemote.ClientHTML)
    })

    http.HandleFunc("/zworemote/cam.png", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "image/png")
        handleImageRequest(zwoasi.WritePNGImage, w, r)
    })

    http.HandleFunc("/zworemote/cam.jpg", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "image/jpeg")
        handleImageRequest(zwoasi.WriteJPGImage, w, r)
    })

    http.HandleFunc("/zworemote/cam.json", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        zwoasi.WriteStats(w)

    })

    log.Print("http.ListenAndServe")
    log.Fatal(http.ListenAndServe(":8080", nil))

}

func handleImageRequest(writerFunc func(origin image.Point, width int, height int, exposure float64, imageWriter io.Writer), w http.ResponseWriter, r *http.Request) {
    x, err := strconv.Atoi(r.FormValue("x"))
    y, err := strconv.Atoi(r.FormValue("y"))
    origin := image.Point{x, y}
    width, err := strconv.Atoi(r.FormValue("w"))
    height, err := strconv.Atoi(r.FormValue("h"))
    e, err := strconv.ParseFloat(r.FormValue("e"), 64)
    log.Print(err)
    writerFunc(origin, width, height, e, w)
}
