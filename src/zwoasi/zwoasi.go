// Copyright 2016 Ted Goddard. All rights reserved.
// Use of this source code is governed the MIT
// license that can be found in the LICENSE file.

package zwoasi

import "bytes"
import "fmt"
import "os"
import "io"
import "bufio"
//import "math"
import "image"
//import "image/color"
import "image/png"
//import "image/jpeg"
import "encoding/binary"
import "unsafe"

/*
#cgo CFLAGS: -I.
//#cgo LDFLAGS: -lstdc++ –framework Foundation -lobjc.A -lusb-1.0 -L/Users/goddards/Documents/development/zworemote/src/zwoasi -lASICamera2 -v
#cgo LDFLAGS: -lstdc++ -lusb-1.0 -L/Users/goddards/Documents/development/zworemote/src/zwoasi -lASICamera2 -v

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include "ASICamera2.h"

#define MAX_CONTROL 7

unsigned char* asiGetImage(char *fileName, int* width, int* height, int* len)  {
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
    int CamNum = 0;
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

	int imageWidth;
	int imageHeight;
	int bin = 1;
    int imageType;

    ASIGetROIFormat(CamNum, &imageWidth, &imageHeight, &bin, (ASI_IMG_TYPE*)&imageType);
    printf("ASIGetROIFormat %dx%d %d\n", imageWidth, imageHeight, imageType);

    *width = imageWidth;
    *height = imageHeight;
    *len = imageWidth * imageHeight;
 
	long ltemp = 0;
	ASI_BOOL bAuto = ASI_FALSE;
	ASIGetControlValue(CamNum, ASI_TEMPERATURE, &ltemp, &bAuto);
	printf("sensor temperature:%02f\n", (float)ltemp/10.0);

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



type GrayFloat64 struct  {
    //maximum pixel value to allow scaling to 16 bits
    MaxPixel float64
    Pix []float64
    Stride int
    Rect image.Rectangle
}


func GetImage(fileName string) image.Image  {
    var widthC C.int
    var heightC C.int
    var lenC C.int
    greyCBytes := C.asiGetImage(C.CString(fileName), &widthC, &heightC, &lenC)
    greyBytes := C.GoBytes(unsafe.Pointer(greyCBytes), lenC)

    width := int(widthC)
    height := int(heightC)
    len := int(lenC)


f, err := os.OpenFile("x.png", os.O_CREATE|os.O_WRONLY, 0666)
if err != nil {
    fmt.Println("file open error: ", err)
}
img := image.NewGray(image.Rect(0, 0, width, height))
img.Pix = greyBytes

err = png.Encode(f, img)
if err != nil {
    fmt.Println("png.Encode error: ",err)
}
fmt.Println("png written")

    maxPixel := 0.0
    var pixel float64
    floatLen := len / 8
    buf := bytes.NewReader(greyBytes)
    var greyFloats = make([]float64, floatLen)
    for i := 0; i < floatLen; i++ {
        err := binary.Read(buf, binary.LittleEndian, &pixel)
        if err != nil {
            fmt.Println("binary.Read failed:", err)
        }
        greyFloats[i] = pixel
        if (pixel > maxPixel)  {
            maxPixel = pixel
        }
    }

/*
    greyImage := GrayFloat64{
        MaxPixel:maxPixel,
        Pix:greyFloats,
        Stride:width,
        Rect:image.Rect(0, 0, width, height),
    }
*/

    C.free(unsafe.Pointer(greyCBytes))
    
    return img
}

func WriteImage(width int, height int, imageWriter io.Writer)  {
    greyImage := GetImage("whatever")

    bufWriter := bufio.NewWriter(imageWriter)
    png.Encode(bufWriter, greyImage)
    bufWriter.Flush()
}

