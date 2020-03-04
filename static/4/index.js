let socket = new WebSocket("ws://127.0.0.1:8085/ws");
let mapid = document.getElementById('mapid');

let bg = document.getElementById("bg");
let title = document.getElementById("title");
let currentPP = document.getElementById("currentPP");
let ifFC = document.getElementById("ifFC");
let state = document.getElementById("state");
let hun = document.getElementById("100");
let fifty = document.getElementById("50");
let miss = document.getElementById("miss");
let cs = document.getElementById("cs");
let ar = document.getElementById("ar");
let od = document.getElementById("od");
let hp = document.getElementById("hp");
let mods = document.getElementById("mods");
const modsImgs = {
    'ez': './static/easy.png',
    'nf': './static/nofail.png',
    'ht': './static/halftime.png',
    'hr': './static/hardrock.png',
    'sd': './static/suddendeath.png',
    'pf': './static/perfect.png',
    'dt': './static/doubletime.png',
    'nc': './static/nightcore.png',
    'hd': './static/hidden.png',
    'fl': './static/flashlight.png',
    'rx': './static/relax.png',
    'ap': './static/autopilot.png',
    'so': './static/spunout.png',
    'at': './static/autoplay.png',
    'cn': './static/cinema.png',
    'v2': './static/v2.png',
}

let aaaaa ={
    'aaaa': {
        'dasda': 'dasda'
    }
}

socket.onopen = () => {
    console.log("Successfully Connected");
};

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
    socket.send("Client Closed!")
};

socket.onerror = error => {
    console.log("Socket Error: ", error);
};
let tempImg;
let tempCs;
let tempAr;
let tempOd;
let tempHp;
let tempTitle;
let tempMods;
let gameState;

socket.onmessage = event => {
    let data = JSON.parse(event.data);
    if(tempImg !== data.menuContainer.innerBG){
        tempImg = data.menuContainer.innerBG
        data.menuContainer.innerBG = data.menuContainer.innerBG.replace(/#/g,'%23').replace(/%/g,'%25')
        bg.setAttribute('src',`http://127.0.0.1:24050/Songs/${data.menuContainer.innerBG}?a=${Math.random(10000)}`)
    }
    if(gameState !== data.menuContainer.osuState){
        gameState = data.menuContainer.osuState
        if(gameState === 2 || gameState === 7 || gameState === 14){
            state.style.transform = "translateY(0)"
        }else{
            state.style.transform = "translateY(-50px)"
        }
    }
    if(tempTitle !== data.menuContainer.bmInfo.split('[')[0]){
        tempTitle = data.menuContainer.bmInfo.split('[')[0];
        title.innerHTML = tempTitle
    }
    if(data.menuContainer.ppSS != tempCs){
        tempCs = JSON.parse(data.menuContainer.ppSS).cs
        cs.innerHTML= `CS: ${Math.round(tempCs * 10) / 10} <hr>`
    }
    if(data.menuContainer.ppSS != tempAr){
        tempAr = JSON.parse(data.menuContainer.ppSS).ar
        ar.innerHTML= `AR: ${Math.round(tempAr * 10) / 10} <hr>`
    }
    if(data.menuContainer.ppSS != tempOd){
        tempOd = JSON.parse(data.menuContainer.ppSS).od
        od.innerHTML= `OD: ${Math.round(tempOd * 10) / 10} <hr>`
    }
    if(data.menuContainer.ppSS != tempHp){
        tempHp = JSON.parse(data.menuContainer.ppSS).hp
        console.log(tempHp);
        hp.innerHTML= `HP: ${Math.round(tempHp * 10) / 10} <hr>`
    }
    if(data.gameplayContainer.pp != ''){
        let ppData = JSON.parse(data.gameplayContainer.pp) 
        currentPP.innerHTML = Math.round(ppData.pp)
    }else{
        currentPP.innerHTML = 0
    }
    if(data.gameplayContainer.ppIfFC != ''){
        let ppData = JSON.parse(data.gameplayContainer.ppIfFC) 
        ifFC.innerHTML = Math.round(ppData.pp)
    }else{
        ifFC.innerHTML = 0
    }
    if(data.gameplayContainer[100] > 0){
        hun.innerHTML = data.gameplayContainer[100]
    }else{
        hun.innerHTML = 0
    }
    if(data.gameplayContainer[50] > 0){
        fifty.innerHTML = data.gameplayContainer[50]
    }else{
        fifty.innerHTML = 0
    }
    if(data.gameplayContainer.miss > 0){
        miss.innerHTML = data.gameplayContainer.miss
    }else{
        miss.innerHTML = 0
    }
    if(tempMods != data.menuContainer.appliedModsString){
        tempMods = data.menuContainer.appliedModsString
        if (tempMods == ""){
            mods.innerHTML = '';
        }
        else{
            mods.innerHTML = '';
            let modsApplied = tempMods.toLowerCase();
            
            if(modsApplied.indexOf('nc') != -1){
                modsApplied = modsApplied.replace('dt','')
            }
            if(modsApplied.indexOf('pf') != -1){
                modsApplied = modsApplied.replace('sd','')
            }
            console.log(modsApplied);
            let modsArr = modsApplied.match(/.{1,2}/g);
            for(let i = 0; i < modsArr.length; i++){
                let mod = document.createElement('div');
                mod.setAttribute('class','mod');
                let modImg = document.createElement('img');
                modImg.setAttribute('src', modsImgs[modsArr[i]]);
                mod.appendChild(modImg);
                mods.appendChild(mod);
            }
        }
    }
}
