package zworemote

    const ClientHTML = `
<html>
    <head>
        <meta name="viewport" content="initial-scale=0.5, width=640, user-scalable=no">
        <meta name="apple-mobile-web-app-status-bar-style" content="black">
        <meta name="apple-mobile-web-app-capable" content="yes">
        <style>` +
            ClientCSS + `
        </style>
        <script>` +
            ClientJavaScript + `
        </script>
    </head>
    <body>
    <div class="imgBox" onclick="imageDispatch(event)">
        <img id="cam" onclick="imageClick(event)" onload="adjustSizes()"
            style="-webkit-filter:brightness(140%%)contrast(300%%);position: relative; top: 0; left: 0;height:100%%;">
        <img id="solvedfield" onload="adjustSizes()" onclick="solvedClick(event)"
            onerror="this.style.display='none';"
            style="position: absolute; top: 0; left: 0;">
        <img id="zoom" class="zoom"
            style="display: none; -webkit-filter:brightness(140%%)contrast(300%%);position: absolute; top: 0; left: 0;width:320px;height:240px;">
        <svg id="bull" width="100%%" height="100%%" style="opacity:0; position: absolute; top: 0; left: 0;">
            <g >
                <line x1="0px" y1="50%%" x2="100%%" y2="50%%" stroke="red" stroke-width="1" />
                <line x1="50%%" y1="0px" x2="50%%" y2="100%%" stroke="red" stroke-width="1" />
                <circle cx="50%%" cy="50%%" r="10%%" stroke="red" stroke-width="1" fill="none" />
                <circle cx="50%%" cy="50%%" r="4%%" stroke="red" stroke-width="1" fill="none" />
                <circle cx="50%%" cy="50%%" r="2%%" stroke="red" stroke-width="1" fill="none" />
            </g>
        </svg>
        <svg id="arrow" width="40" height="40" style="opacity: 0; position: absolute; top: 0; left: 0;" >
            <polygon points="20,0 25,20 22,20 23,40 17,40 18,20 15,20" stroke="red" stroke-width="1.0" fill="firebrick"/>
        </svg>
        <svg id="galIcon" width="40" height="40" style="opacity: 0; position: absolute; top: 0; left: 0;" >
            <circle cx="50%%" cy="50%%" r="48%%" stroke="red" stroke-width="1.0" fill="none"/>
        </svg>
        <svg id="starIcon" width="10" height="10" style="opacity: 0; position: absolute; top: 0; left: 0;" >
            <circle cx="50%%" cy="50%%" r="48%%" stroke="red" stroke-width="1.0" stroke-dasharray="2 2" fill="none"/>
        </svg>
        <svg id="destIcon" width="10" height="10" style="opacity: 0; position: absolute; top: 0; left: 0;">
            <circle cx="50%%" cy="50%%" r="48%%" stroke="red" stroke-width="1.0" fill="none"/>
        </svg>
        <svg id="marker" width="20" height="20" style="position: absolute; top: 0; left: 0;">
            <g id="m-select" style="opacity:0">
                <rect x="-4" y="-4" width="10" height="10" stroke="white" stroke-width="2" fill="none" />
                <rect x="14" y="-4" width="10" height="10" stroke="white" stroke-width="2" fill="none" />
                <rect x="-4" y="14" width="10" height="10" stroke="white" stroke-width="2" fill="none" />
                <rect x="14" y="14" width="10" height="10" stroke="white" stroke-width="2" fill="none" />
            </g>
            <g id="m-calib" style="opacity:0">
                <rect x="0" y="0" width="20" height="20" stroke="yellow" stroke-width="4" stroke-dasharray="2 2" fill="none" />
            </g>
            <g id="m-guide" style="opacity:0">
                <line x1="10" y1="0" x2="10" y2="20" stroke="green" stroke-width="2" />
                <line x1="0" y1="10" x2="20" y2="10" stroke="green" stroke-width="2" />
                <rect x="4" y="4" width="12" height="12" stroke="green" stroke-width="2" fill="none" />
            </g>
            <g id="m-lost"  style="opacity:0">
                <line x1="0" y1="0" x2="20" y2="20" stroke="red" stroke-width="2" />
                <line x1="20" y1="0" x2="0" y2="20" stroke="red" stroke-width="2" />
                <rect x="0" y="0" width="20" height="20" stroke="red" stroke-width="4" fill="none" />
            </g>
        </svg>
        <div id="zoomstats" class="zoomstats" style="display: none; position: absolute; top: 0; left: 0;">42</div>
    </div>
    <div class="rcontrols" >
      <div class="rcinner" >
        <a onclick="expose(100)">0.1s</a>
        <a onclick="expose(1000)">1.0s</a>
        <a onclick="expose(10000)">10.0s</a>
        <a onclick="runSeries()">
          <svg width="40px" height="30px">
            <g stroke-width="1.5" stroke="black" fill="none">
                <rect x="3" y="0" width="20" height="20" />
                <rect x="10" y="5" width="20" height="20" />
                <rect x="17" y="10" width="20" height="20" />
            </g>
          </svg>
        </a>
        <input id="countField" onblur="updateSeriesCount(this)" onchange="updateSeriesCount(this)">
      </div>
    </div>
    <div class="bcontrols" style="display: none;">
      <div class="bcinner" >
        <a onclick="guide()">GUIDE</a>
        <a onclick="stop()">STOP</a>
        <a onclick="loop()">LOOP</a>
      </div>
    </div>
    <div class="tcontrols" >
      <div class="tcinner" >
        <input style="display:none;" onchange="findField(this.value);" >
        <a onclick="toggleLooping()">
            <svg width="40px" height="40px"><g>
                <circle cx="20" cy="20" r="10" stroke="black" stroke-width="1" fill="none" stroke-dasharray="1, 1"/>
                <line x1="6" y1="13" x2="13" y2="14" stroke="black" stroke-width="1"/>
                <line x1="13" y1="21" x2="13" y2="14" stroke="black" stroke-width="1"/>
            </g></svg>
        </a>
      </div>
    </div>
    <div class="trcontrols" >
      <div class="trinner" >
        <a draggable="true"
            ontouchstart="adjustStart(event)" ondragstart="adjustStart(event)"
            ondrag="adjustImage(event)" ontouchmove="adjustImage(event)">
          <svg width="60px" height="60px">
            <g >
                <path d="M30,10 L30,50 A20,20 0 0,1 30,10 z" fill="black" />
                <path d="M30,50 L30,10 A20,20 0 0,1 30,50 z" fill="firebrick" />
            </g>
          </svg>
        </a>
      </div>
      <div class="trinner" >
        <a draggable="true"
            ontouchstart="adjustExposureStart(event)" ondragstart="adjustExposureStart(event)"
            ondrag="adjustExposure(event)" ontouchmove="adjustExposure(event)"
            ontouchend="adjustExposureEnd(event)" ondragend="adjustExposureEnd(event)">
          <svg width="60px" height="60px">
            <g >
                <rect x="10" y="10" width="20" height="40" fill="black" />
                <rect x="30" y="10" width="20" height="40" fill="firebrick" />
            </g>
          </svg>
        </a>
      </div>
      <div class="trinner" >
        <input id="gainField" onblur="adjustExposureWithFields()" onchange="adjustExposureWithFields()" autocomplete="off">
      </div>
      <div class="trinner">
        <input id="exposureField" onblur="adjustExposureWithFields()" onchange="adjustExposureWithFields()" autocomplete="off">
      </div>
    </div>
    <div class="brcontrols" >
      <div class="brinner" >
        <a onclick="toggleSolved()" style="display:none;">
            <svg width="40px" height="40px"><g >
                <animateTransform id="solvedspinner"
                    attributeName="transform"
                    attributeType="XML"
                    type="rotate"
                    from="0 20 20"
                    to="360 20 20"
                    dur="10s"
                    begin="indefinite"
                    repeatCount="indefinite"/>
                <line x1="60%%" y1="30%%" x2="20%%" y2="60%%" stroke="black" stroke-width="1" />
                <line x1="20%%" y1="60%%" x2="80%%" y2="80%%" stroke="black" stroke-width="1" />
                <line x1="80%%" y1="80%%" x2="60%%" y2="30%%" stroke="black" stroke-width="1" />
                <circle cx="60%%" cy="30%%" r="8%%" stroke="black" stroke-width="1" fill="firebrick" />
                <circle cx="20%%" cy="60%%" r="8%%" stroke="black" stroke-width="1" fill="firebrick" />
                <circle cx="80%%" cy="80%%" r="8%%" stroke="black" stroke-width="1" fill="firebrick" />
            </g></svg>
        </a>
        <a onclick="toggleBullseye()">
            <svg width="40px" height="40px"><g >
                <line x1="0px" y1="50%%" x2="100%%" y2="50%%" stroke="black" stroke-width="1" />
                <line x1="50%%" y1="0px" x2="50%%" y2="100%%" stroke="black" stroke-width="1" />
                <circle cx="50%%" cy="50%%" r="20%%" stroke="black" stroke-width="1" fill="none" />
                <circle cx="50%%" cy="50%%" r="10%%" stroke="black" stroke-width="1" fill="none" />
            </g></svg>
        </a>
        <a onclick="toggleGraphs()">
            <svg width="40px" height="40px"><g>
                <path d="M 10, 30 C 20, 30, 15, 10, 20, 10" stroke="black" stroke-width="2" fill="none" stroke-dasharray="1, 1"/>
                <path d="M 20, 10 C 25, 10, 20, 30, 30, 30" stroke="black" stroke-width="2" fill="none" stroke-dasharray="1, 1"/>
            </g></svg>
        </a>
      </div>
    </div>
    <div id="bldisplay" class="bldisplay" >
    </div>
    </body>
</html>
`
