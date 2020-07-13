
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
  
* [**Classic**](https://github.com/l3lackShark/static/tree/master/Classic)

![](https://cdn.discordapp.com/attachments/641255341245333514/731838930340544573/unknown.png)\
Designer: [Dartandr](https://github.com/Dartandr)

* [**MinimalLime**](https://github.com/l3lackShark/static/tree/master/MinimalLime) 

![](https://cdn.discordapp.com/attachments/641255341245333514/731840161612300358/unknown.png)\
Designer: [cyperdark](https://github.com/cyperdark)

* [**MaximalLime**](https://github.com/l3lackShark/static/tree/master/MaximalLime)

![](https://cdn.discordapp.com/attachments/641255341245333514/731841741715669002/unknown.png)\
Designer: [cyperdark](https://github.com/cyperdark)

* [**TrafficLight**](https://github.com/l3lackShark/static/tree/master/TrafficLight)

![](https://cdn.discordapp.com/attachments/641255341245333514/731842011514011698/unknown.png)\
Designer: [cyperdark](https://github.com/cyperdark)

* [**Luscent**](https://github.com/l3lackShark/static/tree/master/Luscent) - Open-Source Implementation of [Luscent's](https://gumroad.com/l/Luscent) overlay. No elements were stolen. This is a remake. Please consider buying his version!

![](https://media.discordapp.net/attachments/641255341245333514/731843129833160704/unknown.png)
Initial Design by [Luscent](https://github.com/inix1257), Remake by [Dartandr](https://github.com/Dartandr)

* **Kerli package** [**Kerli1**](https://github.com/l3lackShark/static/tree/master/Kerli1) [**Kerli2**](https://github.com/l3lackShark/static/tree/master/Kerli2)

![](https://cdn.discordapp.com/attachments/530940222771560452/732038445266108486/Kerli_hud.png)
Designer: [Dartandr](https://github.com/Dartandr)
</details>

# How does it work?

gosumemory streams WebSocket data to **ws://localhost:24050/ws** that you can use in any programming language to develop a frontend. We recommend JavaScript though, as it's much easier to make something pretty with the Web framework. All of the included counters are good starting points. There is also http://localhost:24050/json that you can open in a web browser to see the available data. We strongly recommend against sending GET requests to that address, please **use WebSocket instead**.
 
  
  

# How do I submit a pp counter?

Head over to [static](https://github.com/l3lackShark/static) and create a pull request there. If it's good quality, then it will get approved and included in the next release.

  

# Linux

You have two options. Either run native, but with sudo privileges, or through WINE. If you choose the latter, then please start the program with the `-wine=true` flag.
Please note that we currently don't support 32-Bit builds. You would need a 64-Bit WINEPREFIX in order for it to work.

  

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
