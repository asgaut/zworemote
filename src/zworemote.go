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
import "os/exec"
//import "path/filepath"
//import "strings"
import "regexp"
import "time"
import "bufio"
import "strconv"
import "net/http"
import "log"
import "zwoasi"
import "zwoefw"
import "zworemote"
//import "encoding/json"

var validPrefix = regexp.MustCompile(`^[[:alnum:]]+$`)

func main() {
    fmt.Println("Starting server")
//    cameraNumber := flag.String("cam", "", "device number of camera")
    flag.Parse()

    zwoasi.OpenCamera()
    zwoefw.OpenFilter()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt)
    go func(){
        for _ = range sigChan {
            zwoasi.CloseCamera()
            zwoefw.CloseFilter()
            fmt.Println("Camera closed.")
            os.Exit(0)
        }
    }()

    http.HandleFunc("/zworemote/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, zworemote.ClientHTML)
    })

    http.HandleFunc("/zworemote/cam.png", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "image/png")
        w.Header().Set("Cache-Control", "no-store")
        handleImageRequest(zwoasi.WritePNGImage, w, r)
    })

    http.HandleFunc("/zworemote/cam.jpg", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "image/jpeg")
        w.Header().Set("Cache-Control", "no-store")
        handleImageRequest(zwoasi.WriteJPGImage, w, r)
    })

    http.HandleFunc("/zworemote/cam.mp4", func(w http.ResponseWriter, r *http.Request) {
        //starting point for video stream piped to ffmpeg
        w.Header().Set("Content-Type", "video/mp4")
        w.Header().Set("Cache-Control", "no-store")
        cmd := exec.Command("fgrep", "cat")
        stdin, err := cmd.StdinPipe()
        if nil != err { log.Println(err.Error()) }
        stdout, err := cmd.StdoutPipe()
        if nil != err { log.Println(err.Error()) }
        cmd.Start()
        fmt.Fprintf(stdin, "commands are\n")
        fmt.Fprintf(stdin, "concatenated\n")
        fmt.Fprintf(stdin, "regularly\n")
        stdin.Close()
        _, err = io.Copy(w, stdout)
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

    http.HandleFunc("/zworemote/filter", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "filter position")
        f := formInt(r, "f")
        zwoefw.SetFilterPosition(f)
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
    prefix := r.FormValue("p")
    if (!validPrefix.MatchString(prefix)) {
        prefix = "c"
    }
    f := formInt(r, "f")
    if (f != 0) {
        zwoefw.SetFilterPosition(f)
    }

    config := zwoasi.CaptureConfig{}
    config.Graphs = r.FormValue("graphs")

    fw, fok := w.(http.Flusher)

    i := 0
    camTemperature := 0.0
    friendlyExposure := zwoasi.GetFriendlyTime(e)
    for i = 0; i < n; i++ {
        f := getStampedFile(prefix)
        defer f.Close()
        bufWriter := bufio.NewWriter(f)
        zwoasi.WritePNGImage(origin, width, height, depth, e, g, config ,bufWriter)
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

func getStampedFile(prefix string) *os.File {
    now := time.Now()
    fileName := fmt.Sprintf("/tmp/%s_%d%02d%02d-%02d.%02d.%02d-%x.png",
            prefix,
            now.Year(), now.Month(), now.Day(),
            now.Hour(), now.Minute(), now.Second(), now.Nanosecond() / 1000)
    f, err := os.Create(fileName)
    log.Print(err)
    return f
}

func handleImageRequest(writerFunc func(origin image.Point, width int, height int, depth int, exposure float64, gain float64, config zwoasi.CaptureConfig, imageWriter io.Writer) image.Image, w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    x := formInt(r, "x")
    y := formInt(r, "y")
    origin := image.Point{x, y}
    width := formInt(r, "w")
    height := formInt(r, "h")
    e := formFloat(r, "e")
    g := formFloat(r, "g")
    depth := 8
    config := zwoasi.CaptureConfig{}
    config.Graphs = r.FormValue("graphs")

    writerFunc(origin, width, height, depth, e, g, config, w)
    elapsed := time.Since(start)
    fmt.Printf("handleImageRequest took %s\n", elapsed)
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
