
# gosumemory

Yet another memory reader for osu! Supports both Linux and Windows. (**requires sudo on Linux** though since only root can read /proc)

  

Build custom pp counters with ease!\
You can contact us here: https://discord.gg/8enr4qD
  
 

# Real-World examples:

[FlyingTuna](https://www.twitch.tv/flyingtuna/clip/TransparentObliviousHawkAMPEnergyCherry)\
[Alumetri](https://www.twitch.tv/alumetri/clip/WonderfulVenomousCougarGrammarKing)\
[Sotarks](https://youtu.be/cRlSIOYkZbM?t=26)\
[Mathi](https://www.youtube.com/watch?v=rtmKxbnCQtA)

# Usage

First, download the latest [Release](https://github.com/l3lackShark/gosumemory/releases/latest). Then, in the root folder of the program, you can find the **static** directory. It contains all of the available counters. Those are getting streamed via HTTP-File server. To access one of them, simply specify its name in the URL bar of your browser (Ex: http://localhost:24050/Classic). If using built-in counters covers all of your needs, then you are done here. **Please note that auto-updates only cover the executable itself, however, if a new counter gets released, we will mention it in the Release Notes.** If you want to make your own, just create a new directory in the *static* folder.  

# Included counters:
<details>
  <summary>Click ME</summary>
  
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

</details>

# How does it work?

gosumemory streams WebSocket data to **ws://localhost:24050/ws** that you can use in any programming language to develop a frontend. We recommend JavaScript though, as it's much easier to make something pretty with the Web framework. All of the included counters are good starting points. There is also http://localhost:24050/json that you can open in a web browser to see the available data. We strongly recommend against sending GET requests to that address, please **use WebSocket instead**.
 
  
  

# How do I submit a pp counter?

Head over to [static](https://github.com/l3lackShark/static) and create a pull request there. If it's good quality, then it will get approved and included in the next release.

  

# Linux

You have two options. Either run native, but with sudo privileges, or through WINE. If you choose the latter **AND using a 64-Bit Release**, then please start the program with the `-wine=true` flag to get proper leaderboard data.
Please note that 32-Bit builds are experimental and don't have leaderboard data. They could contain crashes. Report them if you encounter any.

  

# This project depends on:

* [cast](https://github.com/spf13/cast)
* [gorilla-websocket](https://github.com/gorilla/websocket)
* [kiwi](https://github.com/l3lackShark/kiwi)
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
* [KotRik](https://github.com/KotRikD) for porting legacy counters
