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
import "time"
import "bufio"
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

    http.HandleFunc("/zworemote/series", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        handleSeries(r)
        fmt.Fprintf(w, "series complete")
    })

    http.HandleFunc("/zworemote/cam.json", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        zwoasi.WriteStats(w)

    })

    log.Print("http.ListenAndServe")
    log.Fatal(http.ListenAndServe(":8080", nil))

}

func handleSeries(r *http.Request) {
    n, err := strconv.Atoi(r.FormValue("n"))
    x, err := strconv.Atoi(r.FormValue("x"))
    y, err := strconv.Atoi(r.FormValue("y"))
    origin := image.Point{x, y}
    width, err := strconv.Atoi(r.FormValue("w"))
    height, err := strconv.Atoi(r.FormValue("h"))
    depth, err := strconv.Atoi(r.FormValue("d"))
    if (depth == 0) {
        depth = 8
    }
    e, err := strconv.ParseFloat(r.FormValue("e"), 64)
    g, err := strconv.ParseFloat(r.FormValue("g"), 64)
    log.Print(err)
    for i := 0; i < n; i++ {
        now := time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
        fileName := fmt.Sprintf("/tmp/camseries%x.png", now)
        f, err := os.Create(fileName)
        log.Print(err)
        bufWriter := bufio.NewWriter(f)
        zwoasi.WritePNGImage(origin, width, height, depth, e, g, bufWriter)
        bufWriter.Flush()
    }
}

func handleImageRequest(writerFunc func(origin image.Point, width int, height int, depth int, exposure float64, gain float64, imageWriter io.Writer) image.Image, w http.ResponseWriter, r *http.Request) {
    x, err := strconv.Atoi(r.FormValue("x"))
    y, err := strconv.Atoi(r.FormValue("y"))
//save image or not
//    s := r.FormValue("s")
    origin := image.Point{x, y}
    width, err := strconv.Atoi(r.FormValue("w"))
    height, err := strconv.Atoi(r.FormValue("h"))
    e, err := strconv.ParseFloat(r.FormValue("e"), 64)
    g, err := strconv.ParseFloat(r.FormValue("g"), 64)
    depth := 8
log.Print("image size ", width, " ", height)
    log.Print(err)
    writerFunc(origin, width, height, depth, e, g, w)
}
