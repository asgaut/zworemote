package zworemote

    const ClientJavaScript = `

var currentInfo = { };
var currentTarget = "";
var targetRA = 0.0;
var targetDEC = 0.0;
var ws = new WebSocket("ws://" + location.host + "/echo/");
ws.onmessage = function(msg) {console.log(msg.data);
    var msgJSON = JSON.parse(msg.data);
    console.log(msgJSON.Event);
    var marker = document.getElementById("marker");
    if ("LoopingExposures" == msgJSON.Event)  {
        updateCam();
    };
    if ("StartCalibration" == msgJSON.Event)  {
       showMarker("calib");
    };
    if ("GuideStep" == msgJSON.Event)  {
       updateCam();
       showMarker("guide");
    };
    if ("StarLost" == msgJSON.Event)  {
       showMarker("lost");
    };
};

function updateCam() {
    var camImg = document.getElementById("cam");
    camImg.src = "cam.jpg?" + serializeCurrentParams();
    updateURL();
}
function updateSeriesCount(field) {
    currentParams["n"] = parseInt(field.value);
    updateURL();
}
function serializeCurrentParams() {
    var e = currentParams["e"];
    var eClause = e ? "&e=" + e : "";
    var g = currentParams["g"];
    var gClause = g ? "&g=" + g : "";
    var n = currentParams["n"];
    var nClause = n ? "&n=" + n : "&n=3";
    var query = eClause + gClause + nClause + "&r=" + rand();
    return query;
}
function updateURL() {
    window.history.pushState(currentParams, 'ZWO Remote', '?' + serializeCurrentParams());
}
function runSeries() {
    clearTimeout(zoomTimer);
    var seriesURL = "series?" + serializeCurrentParams() + "&d=16";
    httpGet(seriesURL, function(data) { console.log(data)});
}
function showMarker(name) {
    clearMarkers();
    document.getElementById("m-" + name).style["opacity"] = 1.0;
}
function clearMarkers() {
    var marker = document.getElementById("marker");
    for (i = 0; i < marker.childNodes.length; i++)  {
       if (!marker.childNodes[i].style) { continue; };
       marker.childNodes[i].style["opacity"] = 0;
    }
}
function getScaledCoordinates(img, coords) {
    return {
        x: Math.round(coords.x * img.naturalWidth / img.width),
        y: Math.round(coords.y * img.naturalHeight / img.height)
    };
}
function getClickPosition(e) {
    var parentPosition = getPosition(e.currentTarget);
    return {
        x: e.clientX - parentPosition.x,
        y: e.clientY - parentPosition.y
    }
}
function getPosition(element) {
    var x = 0;
    var y = 0;
    while (element) {
        x += (element.offsetLeft - element.scrollLeft +
            element.clientLeft);
        y += (element.offsetTop - element.scrollTop +
            element.clientTop);
        element = element.offsetParent;
    }
    return { x: x, y: y };
}
function rand() {
    return Math.random().toString(16).substring(2);
}
var currentParams = { };
function getQueryParams() {
    if (Object.keys(currentParams).length > 0) {
        return currentParams;
    }
    var qString = document.location.search.substring(1);
    var pairs = qString.split("&");
    var params = {};
    for (i = 0; i < pairs.length; i++)  {
        var pair = pairs[i].split("=");
        params[pair[0]] = pair[1];
    }
    currentParams = params;
    currentParams["e"] = parseFloat(currentParams["e"]);
    currentParams["g"] = parseFloat(currentParams["g"]);
    currentParams["n"] = parseInt(currentParams["n"]);
    return params;
}
var startX = 0;
var startY = 0;
var newX = 0;
var newY = 0;
var camContrast = 3.0;
var camBrightness = 1.4;
var startContrast = 3.0;
var startBrightness = 1.4;
function adjustStart(event) {
    startX = event.pageX;
    startY = event.pageY;
    startContrast = camContrast;
    startBrightness = camBrightness;
}
function adjustImage(event) {
    var deltaX = event.pageX - startX;
    var deltaY = event.pageY - startY;
    camContrast = startContrast + deltaX / 100.0;
    camBrightness = startBrightness + deltaY / 100.0;
    var camElement = document.getElementById("cam");
    camElement.style.webkitFilter =
        "brightness(" + camBrightness + ") contrast(" + camContrast + ")";
}
var camGain = 1.4;
var startExposure = 200;
var startGain = 1;
function adjustExposureStart(event) {
    startX = event.pageX;
    startY = event.pageY;
    startExposure = parseFloat(currentParams["e"]);
    startGain = parseFloat(currentParams["g"]);
}
function adjustExposure(event) {
    var deltaX = event.pageX - startX;
    var deltaY = event.pageY - startY;
    var exposure = Math.max(0, (startExposure + 1.0)  * (1.0 + deltaY / 100.0));
    var gain = Math.max(0, (startGain + 1.0) * (1.0 + deltaX / 100.0));
    adjustExposureWithGainExposure(gain, exposure);
}
function adjustExposureWithFields() {
    var gain = parseFloat(document.getElementById("gainField").value);
    var exposure = parseFloat(document.getElementById("exposureField").value);
    adjustExposureWithGainExposure(gain, exposure);
    updateCam();
}
function adjustExposureWithGainExposure(gain, exposure) {
    currentParams["e"] = exposure;
    currentParams["g"] = gain;
    var display = document.getElementById("bldisplay");
    display.innerHTML = currentParams["e"].toFixed(2) + " x " + currentParams["g"].toFixed(0);
    updateInputFields();
}
function updateInputFields() {
    var gainField = document.getElementById("gainField");
    var exposureField = document.getElementById("exposureField");
    var countField = document.getElementById("countField");
    exposureField.value = currentParams["e"].toFixed(2);
    gainField.value = currentParams["g"].toFixed(0);
    countField.value = currentParams["n"];
}
function adjustExposureEnd(event) {
    updateCam();
}
function imageClick(event) {
    clearTimeout(zoomTimer);
    var imgClick = getClickPosition(event);
//            ws.send(JSON.stringify({method: "set_lock_position",
//                params: [imgClick.x, imgClick.y], id: 42}));

    updateZoom(imgClick.x, imgClick.y);
    var marker = document.getElementById("marker");
    marker.style.top = imgClick.y - 10;
    marker.style.left = imgClick.x - 10;
    showMarker("select");
};
var zoomTimer;
function updateZoom(x, y) {
    var e = currentParams["e"];
    var g = currentParams["g"];
    var graphs = currentParams["graphs"];
    var graphsClause = graphs ? "&graphs=" + graphs : "";
    var zoomElement = document.getElementById("zoom");
    var camElement = document.getElementById("cam");
    var coords = getScaledCoordinates(camElement, {x: x, y: y});
    if (coords.x >= 320) {
        coords.x -= 320;
    }
    if (coords.y >= 240) {
        coords.y -= 240;
    }
    zoomElement.src = src="cam.jpg?" + serializeCurrentParams() + graphsClause + "&w=640&h=480&x=" + coords.x + "&y=" + coords.y + "&r=" + rand();
    zoomElement.style.top = y - 120 + "px";
    zoomElement.style.left = x - 160 + "px";
    zoomTimer = setTimeout(function() { updateZoom(x, y); }, parseFloat(e) + 200.0);
}
function guide() {
    console.log("guide");
    ws.send(JSON.stringify({method:"guide",
        params:[{pixels:1.5, time:8, timeout:40}, false], id:1}));
};
function stop() {
    console.log("stop");
    ws.send(JSON.stringify({"method":"set_paused","params":[true,"full"],"id":2}));
};
function loop() {
    console.log("loop");
    ws.send(JSON.stringify({method:"loop", id:3}));
};
function expose(t) {
    console.log("expose" + t);
//            ws.send(JSON.stringify({method:"set_exposure", params:[t], id:4}));

    var e = t;
    document.location = "?" + serializeCurrentParams();


};
function toggleBullseye() {
    var bullseyeElement = document.getElementById("bull");
    bullseyeElement.style["opacity"] = 1.0 - bullseyeElement.style["opacity"];
}
function toggleSolved() {
    var solvedElement = document.getElementById("solvedfield");
    var solvedSpinner = document.getElementById("solvedspinner");
    var newOpacity = 0.5 - solvedElement.style["opacity"];
    if (newOpacity > 0) {
        solvedSpinner.beginElement();
        solvedElement.src = "solved.jpg?" + new Date().getTime();
        solvedElement.onload = function() {
            solvedElement.style["opacity"] = newOpacity;
            solvedSpinner.endElement();
            getAndDisplayInfo();
       }
    } else {
        solvedElement.style["opacity"] = newOpacity;
    }
}
function toggleGraphs() {
    if (currentParams["graphs"] == "all") {
        delete currentParams["graphs"];
    } else {
        currentParams["graphs"] = "all";
    }
}
function processInfo(data) {
    var newInfo = data.split("\n");
    for (index in newInfo) {
        var entry = newInfo[index].split(" ");
        currentInfo[entry[0]] = entry[1];
    }
    if (currentTarget) {
        findField(currentTarget);
    }
}
function getAndDisplayInfo() {
    httpGet("solved.wcsinfo?" + new Date().getTime(), processInfo);
}
function testMarkers() {
    displayMarker("Medusa", 112.48, 13.20694);
    displayMarker("HD 61199", 114.5745875, 4.94234722);
    displayMarker("HD 61112", 114.45529167, 4.53981389);
    displayMarker("HD 60803", 114.1445625, 5.86169722);
    displayMarker("HD 61696", 115.195175, 5.00018611);
    displayMarker("Procyon", 115.028, 5.187222);
    displayMarker("HD 61664", 115.17254583, 6.62442778);
}
function processLookup(data) {
    console.log("lookup response " + data);
    var coords = JSON.parse(data);
    displayMarker("*", coords.ra, coords.dec);

}
function findField(targetText) {
    targetText = targetText.toLowerCase();
    currentTarget = targetText;
    var targetCoords = targetText.split(",");
    if (targetCoords.length > 1) {
        targetRA = parseFloat(targetCoords[0]);
        targetDEC = parseFloat(targetCoords[1]);
        displayMarker("+", targetRA, targetDEC);
    } else {
        httpGet("lookup?o=" + targetText, processLookup);
    }

}
function displayMarker(label, targetRA, targetDEC) {
    var arrowElement = document.getElementById("arrow");
    var camElement = document.getElementById("cam");
    arrowElement.style.opacity = 1.0;
    var deltaRA = targetRA - currentInfo.ra_center;
    var deltaDEC = targetDEC - currentInfo.dec_center;
    var rads =  currentInfo.orientation / 180.0 * Math.PI;
    var pointRA = deltaRA * Math.cos(rads) - deltaDEC * Math.sin(rads);
    var pointDEC = deltaRA * Math.sin(rads) + deltaDEC * Math.cos(rads);
    var norm = Math.sqrt(pointRA * pointRA + pointDEC * pointDEC);
    var scaledRA = pointRA / norm * camElement.height * 0.4;
    var scaledDEC = pointDEC / norm * camElement.height * 0.4;
    console.log("scaledRA " + scaledRA + " scaledDEC " + scaledDEC);
    var rotation = (90 + 180 * Math.atan(scaledDEC / scaledRA) / Math.PI);
    if (scaledRA < 0) {
        rotation = rotation + 180;
    }
    arrowElement.style.top = ((camElement.height / 2 + scaledDEC)) + "px";
    arrowElement.style.left = (camElement.width / 2 + scaledRA) + "px";
    arrowElement.firstElementChild.setAttribute("transform",
            "rotate(" + rotation + " 20 20)");
}
function depositMarker(label, scaledRA, scaledDEC) {
    var camElement = document.getElementById("cam");
    var marker = document.createElement("div");
    marker.innerHTML = label;
    marker.style.top = ((camElement.height / 2 + scaledDEC)) + "px";
    marker.style.left = (camElement.width / 2 + scaledRA) + "px";
    marker.style.position = "absolute";
    document.body.appendChild(marker);
}
function imageDispatch(event) {
    var solvedElement = document.getElementById("solvedfield");
    if (solvedElement.style["opacity"] > 0) {
        solvedClick(event);
    } else {
        imageClick(event);
    }
}
var solvedClickGal = true;
function solvedClick(event) {
    var camElement = document.getElementById("cam");
    var galIcon = document.getElementById("galIcon");
    var starIcon = document.getElementById("starIcon");
    var destIcon = document.getElementById("destIcon");
    var theIcon = galIcon;

    destIcon.style["opacity"] = 0.0;
    if (!solvedClickGal) {
        theIcon = starIcon;
        destIcon.style["opacity"] = 1.0;
    }

    var pos = getClickPosition(event);
    var iconSize = svgSize(theIcon);
    var destIconSize = svgSize(destIcon);

    theIcon.style.top = pos.y - (iconSize / 2);
    theIcon.style.left = pos.x - (iconSize / 2);
    theIcon.style["opacity"] = 1.0;
    solvedClickGal = !solvedClickGal;

    var totalExtra = svgSize(galIcon) / 2;
    destIcon.style.top = svgTop(starIcon) +
            ((camElement.height / 2) - svgTop(galIcon)) - totalExtra;
    destIcon.style.left = svgLeft(starIcon) -
            (svgLeft(galIcon) - (camElement.width / 2)) - totalExtra;

};
function httpGet(url, callback) {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if ((xhr.readyState) == 4 && (xhr.status == 200)) {
           callback(xhr.responseText);
        }
    }
    xhr.open("GET", url, true);
    xhr.send();
}
function svgTop(elm) {
    return parseFloat(elm.style.top);
}
function svgLeft(elm) {
    return parseFloat(elm.style.left);
}
function svgSize(elm) {
    return parseFloat(elm.getAttribute("height"));
}
function adjustSizes() {
    var bullseyeElement = document.getElementById("bull");
    var camElement = document.getElementById("cam");
    bullseyeElement.style.width = camElement.width;
    bullseyeElement.style.height = camElement.height;
    var solvedElement = document.getElementById("solvedfield");
    solvedElement.style.width = camElement.width;
    solvedElement.style.height = camElement.height;
}
window.onresize = function(event)  {
    adjustSizes();
}
window.onload = function() {
    getQueryParams();
    updateCam();
    httpGet("cam.json", function(text) {
        var stats = JSON.parse(text);
        var display = document.getElementById("bldisplay");
        display.innerHTML = parseFloat(stats["temperature"]) + " &deg;C";
    });
    updateInputFields();
}


`
