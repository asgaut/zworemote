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
        handleSeries(w, r)
        fmt.Fprintf(w, "series complete")
    })

    http.HandleFunc("/zworemote/cam.json", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        zwoasi.WriteStats(w)
    })

    log.Print("http.ListenAndServe")
    log.Fatal(http.ListenAndServe(":8080", nil))

}

func handleSeries(w http.ResponseWriter, r *http.Request) {
    n := formInt(r, "n")
    x := formInt(r, "x")
    y := formInt(r, "y")
    origin := image.Point{x, y}
    width := formInt(r, "w")
    height := formInt(r, "h")
    depth := formInt(r, "d")
    if (depth == 0) {
        depth = 8
    }
    e := formFloat(r, "e")
    g := formFloat(r, "g")
    fw, fok := w.(http.Flusher)

    i := 0
    camTemperature := 0.0
    friendlyExposure := zwoasi.GetFriendlyTime(e)
    for i = 0; i < n; i++ {
        f := getStampedFile()
        defer f.Close()
        bufWriter := bufio.NewWriter(f)
        zwoasi.WritePNGImage(origin, width, height, depth, e, g, bufWriter)
        bufWriter.Flush()
        camTemperature = zwoasi.GetTemperature()
        _, err := fmt.Fprintf(w, "Image %d/%d Exposure %s Gain %.0f Temperature %.0f\u00b0\n", i, n, friendlyExposure, g, camTemperature)
        if fok {
            fw.Flush()
        }
        if nil != err {
            break
        }
    }
    log.Printf("Took %d/%d images Exposure %s Gain %.0f Temperature %.0f\u00b0\n", i, n, friendlyExposure, g, camTemperature)

}

func getStampedFile() *os.File {
    now := time.Now()
    fileName := fmt.Sprintf("/tmp/camseries%d%02d%02d-%02d%02d%02d-%x.png",
            now.Year(), now.Month(), now.Day(),
            now.Hour(), now.Minute(), now.Second(), now.Nanosecond() / 1000)
    f, err := os.Create(fileName)
    log.Print(err)
    return f
}

func handleImageRequest(writerFunc func(origin image.Point, width int, height int, depth int, exposure float64, gain float64, imageWriter io.Writer) image.Image, w http.ResponseWriter, r *http.Request) {
    x := formInt(r, "x")
    y := formInt(r, "y")
    origin := image.Point{x, y}
    width := formInt(r, "w")
    height := formInt(r, "h")
    e := formFloat(r, "e")
    g := formFloat(r, "g")
    depth := 8
    writerFunc(origin, width, height, depth, e, g, w)
}

func formInt(r *http.Request, name string) int {
    value, err := strconv.Atoi(r.FormValue(name))
    if nil != err {
        log.Printf("Failed to decode %s %s %s.", name, r.FormValue(name), err.Error())
        return 0
    }
    return value
}

func formFloat(r *http.Request, name string) float64 {
    value, err := strconv.ParseFloat(r.FormValue(name), 64)
    if nil != err {
        log.Printf("Failed to decode %s %s %s.", name, r.FormValue(name), err.Error())
        return 0.0
    }
    return value
}
