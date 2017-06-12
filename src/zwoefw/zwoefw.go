// Copyright 2017 Ted Goddard. All rights reserved.
// Use of this source code is governed the MIT
// license that can be found in the LICENSE file.

package zwoefw

/*
#cgo CFLAGS: -std=c11 -I${SRCDIR}/include
#cgo darwin LDFLAGS: -framework CoreFoundation -framework IOKit -L${SRCDIR}/lib/mac
#cgo linux,arm LDFLAGS: -L${SRCDIR}/lib/armv7
#cgo linux,amd64 LDFLAGS: -L/lib/x86_64-linux-gnu -lm -L${SRCDIR}/lib/x64
#cgo LDFLAGS: -L/usr/local/lib  -lusb-1.0 -lEFWFilter -lstdc++ -v
#cgo linux LDFLAGS: -ludev


#ifdef WIN32
#include <windows.h>
#elif _POSIX_C_SOURCE >= 199309L
#include <time.h>   // for nanosleep
#else
#include <unistd.h> // for usleep
#endif

static void sleep_ms(int milliseconds) { // cross-platform sleep function
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
#include "EFW_filter.h"

#define MAX_CONTROL 7

int FilterNum = 0;

EFW_INFO EFWInfo;

void efwOpenFilter()  {
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

void efwSetFilterPosition(int position)  {

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

void efwCloseFilter()  {
    EFWClose(FilterNum);
}

*/
import "C"

func OpenFilter() {
    C.efwOpenFilter()
}

func CloseFilter() {
    C.efwCloseFilter()
}

func SetFilterPosition(position int) {
    C.efwSetFilterPosition(C.int(position))
}
