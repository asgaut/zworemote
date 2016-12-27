// Copyright 2016 Ted Goddard. All rights reserved.
// Use of this source code is governed the MIT
// license that can be found in the LICENSE file.

package zwoasi

//import "bytes"
import "fmt"
//import "os"
import "io"
import "bufio"
//import "math"
import "strconv"
import "sync"
import "time"
import "image"
import "image/draw"
import "image/color"
import "image/png"
import "image/jpeg"
//import "encoding/binary"
import "encoding/json"
import "unsafe"

/*
#cgo CFLAGS: -I.
//#cgo LDFLAGS: -lstdc++ –framework Foundation -lobjc.A -lusb-1.0 -L/Users/goddards/Documents/development/zworemote/src/zwoasi -lASICamera2 -v
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit -lstdc++ -L/usr/local/lib  -lusb-1.0 -L${SRCDIR} -lASICamera2 -lEFWFilter -v

#ifdef WIN32
#include <windows.h>
#elif _POSIX_C_SOURCE >= 199309L
#include <time.h>   // for nanosleep
#else
#include <unistd.h> // for usleep
#endif

void sleep_ms(int milliseconds) { // cross-platform sleep function
#ifdef WIN32
    Sleep(milliseconds);
#elif _POSIX_C_SOURCE >= 199309L
    struct timespec ts;
    ts.tv_sec = milliseconds / 1000;
    ts.tv_nsec = (milliseconds % 1000) * 1000000;
    nanosleep(&ts, NULL);
#else
    usleep(milliseconds * 1000);
#endif
}

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include "ASICamera2.h"
#include "EFW_filter.h"

#define MAX_CONTROL 7

char* bayer[] = {"RG","BG","GR","GB"};
char* controls[MAX_CONTROL] = {"Exposure", "Gain", "Gamma", "WB_R", "WB_B", "Brightness", "USB Traffic"};

int CamNum = 0;
int FilterNum = 0;

void asiOpenCamera()  {
    int numDevices = ASIGetNumOfConnectedCameras();
    ASI_CAMERA_INFO ASICameraInfo;

    for (int i = 0; i < numDevices; i++) {
        ASIGetCameraProperty(&ASICameraInfo, i);
        printf("%d %s\n",i, ASICameraInfo.Name);
    }

    if (ASIOpenCamera(CamNum) != ASI_SUCCESS) {
        printf("OpenCamera error\n");
    }
    printf("%s information\n",ASICameraInfo.Name);
    int iMaxWidth, iMaxHeight;
    iMaxWidth = ASICameraInfo.MaxWidth;
    iMaxHeight =  ASICameraInfo.MaxHeight;
    printf("resolution:%dX%d\n", iMaxWidth, iMaxHeight);
    if (ASICameraInfo.IsColorCam) {
        printf("Color Camera: bayer pattern:%s\n",bayer[ASICameraInfo.BayerPattern]);
    } else {
        printf("Mono camera\n");
    }

    ASI_CONTROL_CAPS ControlCaps;
    int iNumOfCtrl = 0;
    ASIGetNumOfControls(CamNum, &iNumOfCtrl);
    for (int i = 0; i < iNumOfCtrl; i++) {
        ASIGetControlCaps(CamNum, i, &ControlCaps);
        printf("%s\n", ControlCaps.Name);
    }
}

void asiCloseCamera()  {
    ASICloseCamera(CamNum);
}

long asiGetTemperature()  {

    long ltemp = 0;
    ASI_BOOL bAuto = ASI_FALSE;
    ASI_ERROR_CODE err;
    err = ASIGetControlValue(CamNum, ASI_TEMPERATURE, &ltemp, &bAuto);

    return ltemp / 10;
}

unsigned char* asiGetImage(char *fileName, int x, int y, int width, int height, int depth, double exposure, double gain, int* widthFound, int* heightFound, int* len)  {

    bool bresult;

    ASI_CAMERA_INFO ASICameraInfo;

    ASIGetCameraProperty(&ASICameraInfo, CamNum);

    if (width == 0) {
        width = ASICameraInfo.MaxWidth;
        height = ASICameraInfo.MaxHeight;
        printf("setting max wxh %dx%d\n", width, height);
    }

    int depthFactor = (depth == 8) ? 1 : 2;
    ASI_IMG_TYPE imageType = (depth == 8) ? ASI_IMG_RAW8 : ASI_IMG_RAW16;

    ASISetROIFormat(CamNum, width, height,  1, imageType);
    ASISetStartPos(CamNum, x, y);
    printf("setting origin %dx%d wxh %dx%d\n", x, y, width, height);

    int imageWidth;
    int imageHeight;
    int bin = 1;

    ASIGetROIFormat(CamNum, &imageWidth, &imageHeight, &bin, (ASI_IMG_TYPE*)&imageType);
    printf("ASIGetROIFormat %dx%d %d\n", imageWidth, imageHeight, imageType);

    *widthFound = imageWidth;
    *heightFound = imageHeight;

    int imageSize = ASICameraInfo.MaxWidth * ASICameraInfo.MaxHeight * depthFactor;
    *len = imageSize;

    printf("Image size %d\n", imageSize);
    unsigned char* imageData;
    imageData = (unsigned char*) malloc(imageSize);

    printf("Setting exposure %f gain %f\n", exposure, gain);

    ASISetControlValue(CamNum, ASI_GAIN, gain, ASI_FALSE);
    ASISetControlValue(CamNum, ASI_EXPOSURE, exposure * 1000, ASI_FALSE);
    ASISetControlValue(CamNum, ASI_BANDWIDTHOVERLOAD, 45, ASI_FALSE);

    printf("Taking exposure\n");

    ASI_EXPOSURE_STATUS status;
    int loopCount = 0;
    ASIStartExposure(CamNum, ASI_FALSE);
    sleep_ms(10);
    status = ASI_EXP_WORKING;
    while(status == ASI_EXP_WORKING) {
        ASIGetExpStatus(CamNum, &status);
        sleep_ms(10);
        loopCount += 1;
    }
    printf("exposure took %d ms\n", loopCount * 10);

    if (status == ASI_EXP_SUCCESS) {
        ASIGetDataAfterExp(CamNum, (unsigned char*)imageData, imageSize);
    }

    ASIStopExposure(CamNum);

    printf("returning imageData\n");

    return(imageData);
}

EFW_INFO EFWInfo;

void asiOpenFilter()  {
    int numFilters = EFWGetNum();
    if (numFilters < 1) {
        printf("No filter wheel connected\n");
    }

    printf("Filter Wheels:\n");
    for (int i = 0; i < numFilters; i++) {
        EFWGetID(i, &EFWInfo.ID);
        EFWGetProperty(EFWInfo.ID, &EFWInfo);
        printf("%d: %s\n", i, EFWInfo.Name);
    }

    if (EFWOpen(FilterNum) != EFW_SUCCESS) {
        printf("Unable to open filter wheel.\n");
    }

}

void asiSetFilterPosition(int position)  {

    position -= 1;

    EFW_ERROR_CODE err;
    while (true) {
        err = EFWGetProperty(FilterNum, &EFWInfo);
        if (err != EFW_ERROR_MOVING ) {
            break;
        }
        sleep_ms(500);
    }
    printf("%d slots: ", EFWInfo.slotNum);
    for (int i = 0; i < EFWInfo.slotNum; i++) {
        printf("%d ", i + 1);
    }
    int currentSlot;
    while(true) {
        err = EFWGetPosition(FilterNum, &currentSlot);
        if (err != EFW_SUCCESS || currentSlot != -1 ) {
            break;
        }
        sleep_ms(500);
    }
    printf("\ncurrent position: %d\n", currentSlot + 1);

    err = EFWSetPosition(FilterNum, position);
    if (err == EFW_SUCCESS) {
        printf("\nMoving...\n");
    } else {
        printf("Failed to move filter wheel.\n");
        return;
    }
    while(true) {
        err = EFWGetPosition(FilterNum, &currentSlot);
        if (err != EFW_SUCCESS || currentSlot != -1 ) {
            break;
        }
        sleep_ms(500);
    }
    printf("\ncurrent position: %d\n", currentSlot + 1);

}

void asiCloseFilter()  {
    EFWClose(FilterNum);
}

*/
import "C"

var mutex = &sync.Mutex{}

func OpenCamera() {
    C.asiOpenCamera()
    C.asiOpenFilter()
}

func CloseCamera() {
    C.asiCloseCamera()
    C.asiCloseFilter()
}

func GetTemperature() float64 {
    var valueC C.long
    mutex.Lock()
    valueC = C.asiGetTemperature()
    mutex.Unlock()
    value := float64(valueC)
    return value
}

func GetFriendlyTime(time float64) string {
    seconds := time / 1000.0
    if seconds < 60.0 {
        return fmt.Sprintf("%.1fs", seconds)
    }
    minutes := int(seconds / 60.0)
    extraSeconds := seconds - float64(minutes) * 60.0
    extraSecondsClause := ""
    if extraSeconds > 5.0 {
        extraSecondsClause = fmt.Sprintf("%.0fs", extraSeconds)
    }
    return fmt.Sprintf("%dmin %s", minutes, extraSecondsClause)
}

func GetStats() map[string]string {
    stats := map[string]string{}

    stats["temperature"] = strconv.FormatFloat(GetTemperature(), 'E', -1, 32)

    return stats
}


func GetImage(x int, y int, width int, height int, depth int, exposure float64, gain float64) image.Image  {
    var widthC C.int
    var heightC C.int
    var lenC C.int

    start := time.Now()
    mutex.Lock()
    greyCBytes := C.asiGetImage(C.CString(""), C.int(x), C.int(y), C.int(width), C.int(height), C.int(depth), C.double(exposure), C.double(gain), &widthC, &heightC, &lenC)
    mutex.Unlock()
    elapsed := time.Since(start)
    fmt.Printf("Exposure took %s\n", elapsed)

    greyBytes := C.GoBytes(unsafe.Pointer(greyCBytes), lenC)

    widthFound := int(widthC)
    heightFound := int(heightC)

    var resultImage image.Image
    if (depth == 8) {
        img := image.NewGray(image.Rect(0, 0, widthFound, heightFound))
        img.Pix = greyBytes
        resultImage = img
    } else {
        swab(greyBytes)
        img := image.NewGray16(image.Rect(0, 0, widthFound, heightFound))
        img.Pix = greyBytes
        resultImage = img
    }

    C.free(unsafe.Pointer(greyCBytes))
    
    return resultImage
}

func SetFilterPosition(position int) {
    C.asiSetFilterPosition(C.int(position))
}

func swab(bytes []byte) {
    for i := 0; i < len(bytes); i = i + 2 {
        tempByte := bytes[i];
        bytes[i] = bytes[i + 1]
        bytes[i + 1] = tempByte
    }
}

type CaptureConfig struct {
    Graphs string
}

func WritePNGImage(origin image.Point, width int, height int, depth int, exposure float64, gain float64, config CaptureConfig, imageWriter io.Writer) image.Image {
    return writeEncodedImage(encodePNG, origin, width, height, depth, exposure, gain, config,  imageWriter)
}

func WriteJPGImage(origin image.Point, width int, height int, depth int, exposure float64, gain float64, config CaptureConfig, imageWriter io.Writer) image.Image {
    return writeEncodedImage(encodeJPG, origin, width, height, depth, exposure, gain, config, imageWriter)
}

func writeEncodedImage(encoder func (imageWriter io.Writer, image image.Image), origin image.Point, width int, height int, depth int, exposure float64, gain float64, config CaptureConfig, imageWriter io.Writer) image.Image {
    if (exposure == 0.0) {
        exposure = 300.0
    }
    if (gain == 0.0) {
        gain = 1.0
    }
    x := origin.X
    y := origin.Y

    greyImage := GetImage(x, y, width, height, depth, exposure, gain)

    fmt.Println("Contrast ", GetContrast(greyImage))

    if ("all" == config.Graphs) {
        start := time.Now()
        greyImage = MarkStars(greyImage)
        elapsed := time.Since(start)
        fmt.Printf("MarkStars took %s\n", elapsed)
    }

    bufWriter := bufio.NewWriter(imageWriter)
    encoder(bufWriter, greyImage)
    bufWriter.Flush()
    return greyImage
}

func GetContrast(img image.Image) float64 {
    total := 0.0
    max := 0.0
    min := 65535.0
    bounds := img.Bounds()
    pixelCount := (bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X)
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
//            r, g, b, a := image.At(x, y).RGBA()
            r, _, _, _ := img.At(x, y).RGBA()
            v := float64(r)
            if (v < min) {
                min = v
            }
            if (v > max) {
                max = v
            }
//image.Set(x, y, color.RGBA{r, 0, 0, a})

            total += v
        }
    }
fmt.Println("                                      min: ", min, " max: ", max)
    return (max - min) / float64(pixelCount)
}

func MarkStars(img image.Image) image.Image {
    bounds := img.Bounds()
    outImage := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
    draw.Draw(outImage, outImage.Bounds(), img, bounds.Min, draw.Src)

    centerX := int64(0)
    centerY := int64(0)
    pixelTotal := int64(0)
    imageMax := 0
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, a := outImage.At(x, y).RGBA()
            if (int(r) > imageMax) {
                imageMax = int(r)
            }
            if (r > 25000) {
                outImage.Set(x, y, color.RGBA64{uint16(r), 0, 0, uint16(a)})
                centerX += int64(x)
                centerY += int64(y)
                pixelTotal += 1

            } else {
                outImage.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
            }
        }
    }

    red := color.RGBA{255, 0, 0, 255}
    blue := color.RGBA{0, 0, 255, 255}
    green := color.RGBA{0, 255, 0, 255}
    swath := 20

    if (pixelTotal > 0) {
        cX := int(centerX / pixelTotal)
        cY := int(centerY / pixelTotal)

        currentMax := bounds.Max.Y
        currentMaxX := 0
        halfMaxDelta := bounds.Max.Y
        halfMax := 0
        halfMaxX := 0
        halfMaxQuest := true
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            accumY := 0
            for dY := -swath; dY < swath; dY++ {
                r, g, b, a := outImage.At(x, cY + dY).RGBA()
                g = g
                b = b
                a = a
                accumY += int(r)
            }
            scale := 500
            if (imageMax > 0) {
                scale = 2 * swath * imageMax / (bounds.Max.Y - bounds.Min.Y)
            }
            gY := bounds.Max.Y - int(accumY / scale)
            if (gY < currentMax) {
                //max is reversed in graph
                currentMax = gY
                halfMax = gY / 2
                currentMaxX = x
                halfMaxQuest = true

            }
            if (halfMaxQuest) {
                if (halfMax - gY < halfMaxDelta) {
                    halfMaxDelta = halfMax - gY
                    halfMaxX = x
                } else {
//                    halfMaxQuest = false
                }
            }
            draw.Draw(outImage, image.Rect(x, gY - 1, x + 1, gY + 1), &image.Uniform{green}, image.ZP, draw.Src)

            draw.Draw(outImage, image.Rect(x - 1, cY - swath, x + 1, cY - swath + 1), &image.Uniform{red}, image.ZP, draw.Src)
            draw.Draw(outImage, image.Rect(x - 1, cY + swath, x + 1, cY + swath + 1), &image.Uniform{red}, image.ZP, draw.Src)
            draw.Draw(outImage, image.Rect(x - 1, cY - 1, x + 1, cY + 1), &image.Uniform{red}, image.ZP, draw.Src)
        }

        draw.Draw(outImage, image.Rect(cX - 5, cY - 5, cX + 5, cY + 5), &image.Uniform{red}, image.ZP, draw.Src)
        draw.Draw(outImage, image.Rect(currentMaxX, cY - 5, halfMaxX, cY + 5), &image.Uniform{blue}, image.ZP, draw.Src)
    }

    return outImage

}


func encodePNG(imageWriter io.Writer, image image.Image) {
    png.Encode(imageWriter, image)
}

func encodeJPG(imageWriter io.Writer, image image.Image) {
    quality := 30
    bounds := image.Bounds()
    if (bounds.Max.Y - bounds.Min.Y < 1000) {
        quality = 90
    }
    jpeg.Encode(imageWriter, image, &jpeg.Options{Quality: quality})
}

func WriteStats(jsonWriter io.Writer)  {
    camStats := GetStats()
    camStatsJSON, _ := json.Marshal(camStats)
    bufWriter := bufio.NewWriter(jsonWriter)
    bufWriter.WriteString(string(camStatsJSON))
    bufWriter.Flush()
}
