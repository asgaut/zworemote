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
import "image"
//import "image/color"
import "image/png"
import "image/jpeg"
//import "encoding/binary"
import "encoding/json"
import "unsafe"

/*
#cgo CFLAGS: -I.
//#cgo LDFLAGS: -lstdc++ –framework Foundation -lobjc.A -lusb-1.0 -L/Users/goddards/Documents/development/zworemote/src/zwoasi -lASICamera2 -v
#cgo LDFLAGS: -lstdc++ -L/usr/local/lib -lusb-1.0 -L${SRCDIR} -lASICamera2 -v

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include "ASICamera2.h"

#define MAX_CONTROL 7

char* bayer[] = {"RG","BG","GR","GB"};
char* controls[MAX_CONTROL] = {"Exposure", "Gain", "Gamma", "WB_R", "WB_B", "Brightness", "USB Traffic"};

int CamNum = 0;

void asiOpenCamera()  {
    int numDevices = ASIGetNumOfConnectedCameras();
    ASI_CAMERA_INFO ASICameraInfo;

    for(int i = 0; i < numDevices; i++) {
        ASIGetCameraProperty(&ASICameraInfo, i);
        printf("%d %s\n",i, ASICameraInfo.Name);
    }

    if(ASIOpenCamera(CamNum) != ASI_SUCCESS) {
        printf("OpenCamera error\n");
    }
    printf("%s information\n",ASICameraInfo.Name);
    int iMaxWidth, iMaxHeight;
    iMaxWidth = ASICameraInfo.MaxWidth;
    iMaxHeight =  ASICameraInfo.MaxHeight;
    printf("resolution:%dX%d\n", iMaxWidth, iMaxHeight);
    if(ASICameraInfo.IsColorCam) {
        printf("Color Camera: bayer pattern:%s\n",bayer[ASICameraInfo.BayerPattern]);
    } else {
        printf("Mono camera\n");
    }

    ASI_CONTROL_CAPS ControlCaps;
    int iNumOfCtrl = 0;
    ASIGetNumOfControls(CamNum, &iNumOfCtrl);
    for ( int i = 0; i < iNumOfCtrl; i++)
    {
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

unsigned char* asiGetImage(char *fileName, int x, int y, int width, int height, double exposure, double gain, int* widthFound, int* heightFound, int* len)  {

    bool bresult;

    ASI_CAMERA_INFO ASICameraInfo;

    ASIGetCameraProperty(&ASICameraInfo, CamNum);

    if (width == 0) {
        width = ASICameraInfo.MaxWidth;
        height = ASICameraInfo.MaxHeight;
        printf("setting max wxh %dx%d\n", width, height);
    }

    ASISetROIFormat(CamNum, width, height,  1, ASI_IMG_RAW8);
    ASISetStartPos(CamNum, x, y);
    printf("setting origin %dx%d wxh %dx%d\n", x, y, width, height);

    int imageWidth;
    int imageHeight;
    int bin = 1;
    int imageType;

    ASIGetROIFormat(CamNum, &imageWidth, &imageHeight, &bin, (ASI_IMG_TYPE*)&imageType);
    printf("ASIGetROIFormat %dx%d %d\n", imageWidth, imageHeight, imageType);

    *widthFound = imageWidth;
    *heightFound = imageHeight;
    *len = imageWidth * imageHeight;

    int imageSize = ASICameraInfo.MaxWidth * ASICameraInfo.MaxHeight; //Assume RAW8
    printf("Image size %d\n", imageSize);
//    unsigned char imageData[imageSize];
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
    usleep(10000);//10ms
    status = ASI_EXP_WORKING;
    while(status == ASI_EXP_WORKING) {
        ASIGetExpStatus(CamNum, &status);
        usleep(10000);//10ms
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

*/
import "C"

var mutex = &sync.Mutex{}

func OpenCamera() {
    C.asiOpenCamera()
}

func CloseCamera() {
    C.asiCloseCamera()
}

func GetTemperature() float64 {
    var valueC C.long
    mutex.Lock()
    valueC = C.asiGetTemperature()
    mutex.Unlock()
    value := float64(valueC)
    return value
}

func GetStats() map[string]string {
    stats := map[string]string{}

    stats["temperature"] = strconv.FormatFloat(GetTemperature(), 'E', -1, 32)

    return stats
}


func GetImage(x int, y int, width int, height int, exposure float64, gain float64) image.Image  {
    var widthC C.int
    var heightC C.int
    var lenC C.int

    mutex.Lock()
    greyCBytes := C.asiGetImage(C.CString(""), C.int(x), C.int(y), C.int(width), C.int(height), C.double(exposure), C.double(gain), &widthC, &heightC, &lenC)
    mutex.Unlock()
    greyBytes := C.GoBytes(unsafe.Pointer(greyCBytes), lenC)

    widthFound := int(widthC)
    heightFound := int(heightC)
//    len := int(lenC)

    img := image.NewGray(image.Rect(0, 0, widthFound, heightFound))
    img.Pix = greyBytes

    C.free(unsafe.Pointer(greyCBytes))
    
    return img
}

func WritePNGImage(origin image.Point, width int, height int, exposure float64, gain float64, imageWriter io.Writer) image.Image {
    return writeEncodedImage(encodePNG, origin, width, height, exposure, gain, imageWriter)
}

func WriteJPGImage(origin image.Point, width int, height int, exposure float64, gain float64, imageWriter io.Writer) image.Image {
    return writeEncodedImage(encodeJPG, origin, width, height, exposure, gain, imageWriter)
}

func writeEncodedImage(encoder func (imageWriter io.Writer, image image.Image), origin image.Point, width int, height int, exposure float64, gain float64, imageWriter io.Writer) image.Image {
    if (exposure == 0.0) {
        exposure = 300.0
    }
    if (gain == 0.0) {
        gain = 1.0
    }
    x := origin.X
    y := origin.Y

    greyImage := GetImage(x, y, width, height, exposure, gain)

    fmt.Println("Contrast ", GetContrast(greyImage))

    bufWriter := bufio.NewWriter(imageWriter)
    encoder(bufWriter, greyImage)
    bufWriter.Flush()
    return greyImage
}

func GetContrast(image image.Image) float64 {
    total := 0.0
    max := 0.0
    min := 65535.0
    bounds := image.Bounds()
    pixelCount := (bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X)
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
//            r, g, b, a := image.At(x, y).RGBA()
            r, _, _, _ := image.At(x, y).RGBA()
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
