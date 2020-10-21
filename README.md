# gosumemory

Yet another memory reader for osu! Supports both Linux and Windows. (**requires sudo on Linux** though since only root can read /proc)

Build custom pp counters with ease!\
You can contact us here: https://discord.gg/8enr4qD

# Real-World examples:

[![cpol](https://img.shields.io/badge/cpol%20v5.0---?style=for-the-badge&color=DFD895)](https://youtu.be/neHcOLycieE)
[![FlyingTuna](https://img.shields.io/badge/FlyingTuna%20v1.5---?style=for-the-badge&color=527FD5)](https://www.twitch.tv/flyingtuna/clip/TransparentObliviousHawkAMPEnergyCherry)
[![Alumetri](https://img.shields.io/badge/Alumetri%20v1.2---?style=for-the-badge&color=FF94B6)](https://mega.nz/file/QV1gTKoI#j1QRjDkrjnFvIhyb9JuGi3g_0XZCFzXEXz9PKWcxgmI)
[![Sotarks](https://img.shields.io/badge/Sotarks%20v1.0---?style=for-the-badge&color=C63F55)](https://mega.nz/file/oAlmlQoY#8ABeJPGboMLgCiaY5vR21HX2Km--_jiwqRHOmUJvVmg)
[![Mrekk](https://img.shields.io/badge/Mrekk%20v1.0---?style=for-the-badge&color=72a0d4)](https://mega.nz/file/UZsEUKDK#Ji3JAUr8_04Q7u0RG1BAJFGzZ2-CRhRZkEQqdXVrv60)
[![Mathi](https://img.shields.io/badge/Mathi%20v1.5---?style=for-the-badge&color=4981CE)](https://mega.nz/file/5dsk1QJD#noUKykU5qJYv53I2DPZ7PY2CIQOftS1ufqzOh4rqOb8)
[![RyuK](https://img.shields.io/badge/RyuK%20v1.0---?style=for-the-badge&color=f72f4d)](https://mega.nz/file/dY8k1YyZ#1Phdta1CzxXDotjtllUKsZunnCdliYlQ1VrZ_BNaNIs)

# Usage
     
1. [Download the latest Release](https://github.com/l3lackShark/gosumemory/releases/latest)
    * Unzip files anywhere
    > In the root folder of the program, you can find the **static** directory. It contains all of the available counters. Those are getting streamed via HTTP-File server

2. Run gosumemory & osu!
    * Visit this link in your browser: http://localhost:24050 and choose the counter that you like.
    * Add a browser source in OBS (Width and Height could be found in the **Included counters** spoiler)
4. If using built-in counters covers all of your needs, then you are done here.
> **Please note that auto-updates only cover the executable itself, however, if a new counter gets released, we will mention it in the Release Notes.**\
> If you want to make your own, just create a new directory in the *static* folder.  

# Included counters:
<details>
  <summary>Click ME</summary>

### MonokaiPane

> Size: 512x150\
> *Song Selection*\
<img src="https://i.imgur.com/T8p0R29.png" width="500">\
>*Gameplay 1*\
<img src="https://i.imgur.com/TAmHvFM.png" width="500">\
>*Gameplay 2*\
<img src="https://i.imgur.com/FpHkdLg.png" width="500">\
By: Xynogen<br>

### Classic

> Size: 550x300\
<img  src="https://cdn.discordapp.com/attachments/641255341245333514/731838930340544573/unknown.png"  width="500">\
By: [Dartandr][1]<br>

### OldClassic

> Size: 550x300\
<img  src="https://cdn.discordapp.com/attachments/530940222771560452/732545954468593664/unknown.png"  width="500">\
By: [Dartandr][1]<br>
  
### DarkAndWhite

> Size: 840x140\
<img  src="https://i.imgur.com/mBN375B.jpg"  width="500">\
By: [cyperdark][2]<br>

### Kerli1 & Kerli2

> Size (1)(2): 794x124 | 353x190\
<img  src="https://i.imgur.com/n2w260o.jpg"  width="500">\
By: [Dartandr][1]<br>

### Luscent

> Size: 1920x1080\
Open-Source Implementation of [Luscent's][3] overlay. No elements were stolen. This is a remake. Please [consider buying](https://gumroad.com/l/Luscent) his version!\
<img  src="https://media.discordapp.net/attachments/641255341245333514/731843129833160704/unknown.png"  width="500">\
Remake by: [Dartandr][1]

### VictimCrasherCompact

> Size: 550x132\
> *Song Selection*\
<img  src="https://i.imgur.com/1F1GK3Z.png" width="500">\
>
> *Gameplay*\
<img  src="https://i.imgur.com/epx6dij.png" width="500">\
By: [VictimCrasher][4]<br>
  
### VictimCrasherOverlay

> Size: 1920x1080\
<img  src="https://i.imgur.com/Wo6wI1B.png"  width="500">\
By: [VictimCrasher][4]<br>

### UnstableRate

> Size: 300x100\
Just a plain number that shows current UnstableRate, could be useful if you want to put it above your UR Bar.\
By: [Dartandr][1]

### MaximalLime

> Size: 800x306\
<img  src="https://cdn.discordapp.com/attachments/641255341245333514/731841741715669002/unknown.png"  width="500">\
By: [cyperdark][2]<br>

### MinimalLime

> Size: 640x130\
<img  src="https://cdn.discordapp.com/attachments/641255341245333514/731840161612300358/unknown.png"  width="500">\
By: [cyperdark][2]<br>

### TrafficLight

> Size: 458x380\
<img  src="https://cdn.discordapp.com/attachments/641255341245333514/731842011514011698/unknown.png">\
By: [cyperdark][2]<br>

[1]: https://github.com/Dartandr

[2]: https://github.com/cyperdark

[3]: https://github.com/inix1257

[4]: https://github.com/VictimCrasher

</details>

# How does it work?

gosumemory streams WebSocket data to **ws://localhost:24050/ws** that you can use in any programming language to develop a frontend. We recommend JavaScript though, as it's much easier to make something pretty with the Web framework. All of the included counters are good starting points. There is also http://localhost:24050/json that you can open in a web browser to see the available data. We strongly recommend against sending GET requests to that address in production, please **use WebSocket instead**.

**[Example JSON and a little wiki of it's values](https://github.com/l3lackShark/gosumemory/wiki/JSON-values)**

# What if I don't know any programming languages but still want to output data to OBS?
https://www.youtube.com/watch?v=8ApXBEO5bes 

# How do I submit a pp counter?

Head over to [static](https://github.com/l3lackShark/static) and create a pull request there. If it's good quality, then it will get approved and included in the next release.

# Tournament Client

When operating in tourney mode, real-time pp counters for each client, leaderboard and grades don't work. Each gameplay object is sorted by tournament client "id", "menu" object is a tournament manager (state 22).    

# Linux

You have two options. Either run native, but with sudo privileges, or through WINE.
Please note that Linux builds are not well tested and could contain crashes. Report them if you encounter any.

# Consider supporting
<a href="https://www.buymeacoffee.com/BlackShark" target="_blank"><img src="https://cdn.discordapp.com/attachments/530940222771560452/750051751025049741/default-blue_1.png" alt="Buy Me A Coffee" style="height: 51px !important;width: 217px !important;" ></a>

# This project depends on:

* [cast](https://github.com/spf13/cast)
* [gorilla-websocket](https://github.com/gorilla/websocket)
* [open](https://github.com/skratchdot/open-golang)
* [semver](https://github.com/blang/semver)
* [selfupdate](https://github.com/rhysd/go-github-selfupdate)
* [pretty-print](https://github.com/k0kubun/pp)
* [mp3](https://github.com/tcolgate/mp3)
* [go-windows](https://github.com/elastic/go-windows)

# Special Thanks to:

* [Piotrekol](https://github.com/Piotrekol/) and his [ProcessMemoryDataFinder](https://github.com/Piotrekol/ProcessMemoryDataFinder) for most of the memory signatures
* [Francesco149](https://github.com/Francesco149) and his [oppai-ng](https://github.com/Francesco149/oppai-ng) for the pp calculator that we use
* [tdeo](https://github.com/tadeokondrak) for the [Memory Signature Scanner](https://github.com/l3lackShark/gosumemory/tree/master/mem) package
* [omkelderman](https://github.com/omkelderman) for helping out with the [db](https://github.com/l3lackShark/gosumemory/tree/master/db) package
* [jamuwu](https://github.com/jamuwu/osu-strain) and his [osu-strain](https://github.com/jamuwu/osu-strain) for difficulty strain logic
* [cyperdark](https://github.com/cyperdark) and [Dartandr](https://github.com/Dartandr) for frontend designs
* [KotRik](https://github.com/KotRikD) for making an [OBS Script](https://github.com/l3lackShark/gosumemory-helpers/blob/master/gosumemory-reader.py) and porting legacy counters
