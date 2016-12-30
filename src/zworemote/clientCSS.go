package zworemote

    const ClientCSS = `
  body {
    background-color: #202020;
  }
  input {
    background: red;
  }
  .imgBox {
    position: relative;+
    left: 0;
    top: 0;
  }
  .brcontrols {
    position:fixed;
    bottom:10px;
    right:10px;
  }
  .bldisplay {
    color: red;
    font-size:20px;
    position:fixed;
    bottom:10px;
    left:10px;
  }
  .trcontrols {
    position:fixed;
    top:10px;
    right:10px;
    display: flex;
    flex-direction: column;
  }
  .brcontrols a, .tcontrols a, .trcontrols a {
    display:block;
    padding:10px;
    margin:10px;
    font-size:20px;
    border-radius:8px;
    background:red;
  }
  .tcontrols a {
    width: 40px;
  }
  .trinner {
    align-self: center;
  }
  .trinner input {
    max-width: 70px;
  }
  .rcinner {
    display: flex;
    flex-direction: column;
  }
  .rcinner input {
    align-self: center;
    max-width: 50px;
  }
  .rcinner svg {
    padding: 2px
  }
  .zoomstats {
    color: red;
    font-size:20px;
  }
  @media (max-width: 640px) {
      .bcontrols {
        position:fixed;
        bottom:100px;
        left:60px;
      }
      .tcontrols {
        position:fixed;
        top:100px;
        left:60px;
      }
      .bcinner {
      }
      .rcontrols {
        position:fixed;
        top:100px;
        right:10px;
      }
      .rcinner {
      }
      .bcontrols a {
        height:40px;
        padding:10px;
        margin:10px;
        font-size:40px;
        border-radius:8px;
        background:red;
      }
      .rcontrols a {
        display:block;
        padding:10px;
        margin:10px;
        font-size:40px;
        border-radius:8px;
        background:red;
      }
      .brcontrols {
        position:fixed;
        bottom:100px;
        right:10px;
      }
      .trcontrols {
        position:fixed;
        top:10px;
        right:10px;
      }
  }
  @media (min-width: 641px) {
      .tcontrols {
        position:fixed;
        top:20px;
        left:50%%;
      }
      .bcontrols {
        position:fixed;
        bottom:20px;
        left:50%%;
      }
      .bcinner {
        margin-left:-50%%;
      }
      .tcinner {
        margin-left:-50%%;
      }
      .tcinner input {
        background: red;
        color: black;
      }
      .rcontrols {
        position:fixed;
        top:50%%;
        right:0px;
      }
      .rcinner {
        margin-top: -50%%;
      }
      .bcontrols a {
        height:40px;
        padding:10px;
        margin:20px;
        font-size:20px;
        border-radius:8px;
        background:red;
      }
      .rcontrols a {
        display:block;
        padding:10px;
        margin:10px;
        font-size:20px;
        border-radius:8px;
        background:red;
      }
  }
`
