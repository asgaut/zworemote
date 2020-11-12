// Copyright 2016 Ted Goddard. All rights reserved.
// Use of this source code is governed the MIT
// license that can be found in the LICENSE file.

package main

import "fmt"
import "flag"
import "os"
import "os/signal"
import "strconv"
import "net/http"
import "log"
import "zworemote/zwoefw"
//import "encoding/json"

func main() {
    fmt.Println("Starting server")
//    cameraNumber := flag.String("cam", "", "device number of camera")
    flag.Parse()

    zwoefw.OpenFilter()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt)
    go func(){
        for _ = range sigChan {
            zwoefw.CloseFilter()
            fmt.Println("Filter Wheel closed.")
            os.Exit(0)
        }
    }()

    http.HandleFunc("/efwremote/filter", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "filter position")
        f := formInt(r, "f")
        zwoefw.SetFilterPosition(f)
    })

    log.Print("http.ListenAndServe")
    log.Fatal(http.ListenAndServe(":8081", nil))

}

func formInt(r *http.Request, name string) int {
    value, err := strconv.Atoi(r.FormValue(name))
    if nil != err {
        log.Printf("Failed to decode %s %s %s.", name, r.FormValue(name), err.Error())
        return 0
    }
    return value
}
