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
//import "image/jpeg"
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

int CamNum = 0;

long asiGetTemperature()  {
	if(ASIOpenCamera(CamNum) != ASI_SUCCESS) {
		printf("OpenCamera error\n");
    }

	long ltemp = 0;
	ASI_BOOL bAuto = ASI_FALSE;
    ASI_ERROR_CODE err;
	err = ASIGetControlValue(CamNum, ASI_TEMPERATURE, &ltemp, &bAuto);
printf("%d error\n",err);

	ASICloseCamera(CamNum);

	return ltemp;
}

unsigned char* asiGetImage(char *fileName, int x, int y, int* width, int* height, int* len)  {
	char* bayer[] = {"RG","BG","GR","GB"};
	char* controls[MAX_CONTROL] = {"Exposure", "Gain", "Gamma", "WB_R", "WB_B", "Brightness", "USB Traffic"};

	bool bresult;
    unsigned char *pixelss;
    pixelss = (unsigned char *) malloc(*len);

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

    ASISetROIFormat(CamNum, 640, 480,  1, ASI_IMG_RAW8);
    ASISetStartPos(CamNum, x, y);

	int imageWidth;
	int imageHeight;
	int bin = 1;
    int imageType;

    ASIGetROIFormat(CamNum, &imageWidth, &imageHeight, &bin, (ASI_IMG_TYPE*)&imageType);
    printf("ASIGetROIFormat %dx%d %d\n", imageWidth, imageHeight, imageType);

    *width = imageWidth;
    *height = imageHeight;
    *len = imageWidth * imageHeight;

    int imageSize = ASICameraInfo.MaxWidth * ASICameraInfo.MaxHeight; //Assume RAW8
    printf("Image size %d\n", imageSize);
//    unsigned char imageData[imageSize];
    unsigned char* imageData;
    imageData = (unsigned char*) malloc(imageSize);

    printf("Setting exposure\n");

	ASISetControlValue(CamNum, ASI_GAIN, 500, ASI_FALSE);
	ASISetControlValue(CamNum, ASI_EXPOSURE, 1 * 1000*1000, ASI_FALSE);
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
	ASICloseCamera(CamNum);

    printf("returning imageData\n");

    return(imageData);
}

*/
import "C"

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


func GetImage(x int, y int, widthOut int, heightOut int) image.Image  {
    var widthC C.int
    var heightC C.int
    var lenC C.int
    greyCBytes := C.asiGetImage(C.CString(""), C.int(x), C.int(y), &widthC, &heightC, &lenC)
    greyBytes := C.GoBytes(unsafe.Pointer(greyCBytes), lenC)

    width := int(widthC)
    height := int(heightC)
//    len := int(lenC)

    img := image.NewGray(image.Rect(0, 0, width, height))
    img.Pix = greyBytes

    C.free(unsafe.Pointer(greyCBytes))
    
    return img
}

func WriteImage(x int, y int, width int, height int, imageWriter io.Writer)  {
    greyImage := GetImage(x, y, width, height)

    bufWriter := bufio.NewWriter(imageWriter)
    png.Encode(bufWriter, greyImage)
    bufWriter.Flush()
}

func WriteStats(jsonWriter io.Writer)  {
    camStats := GetStats()
    camStatsJSON, _ := json.Marshal(camStats)
    bufWriter := bufio.NewWriter(jsonWriter)
    bufWriter.WriteString(string(camStatsJSON))
    bufWriter.Flush()
}
