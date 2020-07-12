# gosumemory
Yet another memory reader for osu!

Build custom pp counters with ease!


# Real-World examples:
[FlyingTuna](https://www.twitch.tv/flyingtuna/clip/TransparentObliviousHawkAMPEnergyCherry)\
[Alumetri](https://www.twitch.tv/alumetri/clip/WonderfulVenomousCougarGrammarKing)\
[Sotarks](https://youtu.be/cRlSIOYkZbM?t=26)\
[Mathi](https://www.youtube.com/watch?v=rtmKxbnCQtA)

# Included counters:
As of now, we have these:
TBA


# How does it work?
gosumemory streams WebSocket data to **ws://localhost:24050/ws** that you can use in any programming language to develop a frontend. We recommend JavaScript though, as it's much easier to make something pretty with the Web framework. All of the included counters are good starting points. There is also a http://localhost:24050/json that you can open in a web browser to see the available data. We strongly recommend against sending GET requests to that address, please **use WebSocket** instead.


# How do I submit a pp counter?
Head over to [static](https://github.com/l3lackShark/static) and create a pull request there. If it's good quality, then it will get approved and will be included in the next release.

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
* [tdeo](https://github.com/tadeokondrak) for the [Memory Signature Scanner](https://github.com/l3lackShark/gosumemory/tree/master/mem) package  
* [omkelderman](https://github.com/omkelderman) for helping out with the [db](https://github.com/l3lackShark/gosumemory/tree/master/db) package
* [jamuwu](https://github.com/jamuwu/osu-strain) and his [osu-strain](https://github.com/jamuwu/osu-strain) for difficulty strain logic
* [cyperdark](https://github.com/cyperdark) and [Dartandr](https://github.com/Dartandr) for frontend designs
