// Copyright 2016 Ted Goddard. All rights reserved.
// Use of this source code is governed the MIT
// license that can be found in the LICENSE file.

package zwoasi

//import "bytes"
//import "fmt"
//import "os"
import "io"
import "bufio"
//import "math"
import "strconv"
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
    if(ASIOpenCamera(CamNum) != ASI_SUCCESS) {
        printf("OpenCamera error\n");
    }

    long ltemp = 0;
    ASI_BOOL bAuto = ASI_FALSE;
    ASI_ERROR_CODE err;
    err = ASIGetControlValue(CamNum, ASI_TEMPERATURE, &ltemp, &bAuto);

    ASICloseCamera(CamNum);

    return ltemp;
}

unsigned char* asiGetImage(char *fileName, int x, int y, int width, int height, double exposure, int* widthFound, int* heightFound, int* len)  {

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

    printf("Setting exposure %f\n", exposure);

    ASISetControlValue(CamNum, ASI_GAIN, 500, ASI_FALSE);
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

func OpenCamera() {
    C.asiOpenCamera()
}

func CloseCamera() {
    C.asiCloseCamera()
}

func GetTemperature() float64 {
    var valueC C.long
    valueC = C.asiGetTemperature()
    value := float64(valueC)
    return value
}

func GetStats() map[string]string {
    stats := map[string]string{}

    stats["temerature"] = strconv.FormatFloat(GetTemperature(), 'E', -1, 32)

    return stats
}


func GetImage(x int, y int, width int, height int, exposure float64) image.Image  {
    var widthC C.int
    var heightC C.int
    var lenC C.int

    greyCBytes := C.asiGetImage(C.CString(""), C.int(x), C.int(y), C.int(width), C.int(height), C.double(exposure), &widthC, &heightC, &lenC)
    greyBytes := C.GoBytes(unsafe.Pointer(greyCBytes), lenC)

    widthFound := int(widthC)
    heightFound := int(heightC)
//    len := int(lenC)

    img := image.NewGray(image.Rect(0, 0, widthFound, heightFound))
    img.Pix = greyBytes

    C.free(unsafe.Pointer(greyCBytes))
    
    return img
}

func WritePNGImage(origin image.Point, width int, height int, exposure float64, imageWriter io.Writer)  {
    writeEncodedImage(encodePNG, origin, width, height, exposure, imageWriter)
}

func WriteJPGImage(origin image.Point, width int, height int, exposure float64, imageWriter io.Writer)  {
    writeEncodedImage(encodeJPG, origin, width, height, exposure, imageWriter)
}

func writeEncodedImage(encoder func (imageWriter io.Writer, image image.Image), origin image.Point, width int, height int, exposure float64, imageWriter io.Writer)  {
    if (exposure == 0.0) {
        exposure = 300.0
    }
    x := origin.X
    y := origin.Y
    greyImage := GetImage(x, y, width, height, exposure)

    bufWriter := bufio.NewWriter(imageWriter)
    encoder(bufWriter, greyImage)
    bufWriter.Flush()
}

func encodePNG(imageWriter io.Writer, image image.Image) {
    png.Encode(imageWriter, image)
}

func encodeJPG(imageWriter io.Writer, image image.Image) {
    jpeg.Encode(imageWriter, image, &jpeg.Options{Quality: 90})
}

func WriteStats(jsonWriter io.Writer)  {
    camStats := GetStats()
    camStatsJSON, _ := json.Marshal(camStats)
    bufWriter := bufio.NewWriter(jsonWriter)
    bufWriter.WriteString(string(camStatsJSON))
    bufWriter.Flush()
}
