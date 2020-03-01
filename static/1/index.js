let socket = new WebSocket("ws://127.0.0.1:8085/ws");

let bg = document.getElementById("bg");
let star = document.getElementById("star");
let pp = document.getElementById("pp");
let hun = document.getElementById("h100");
let fifty = document.getElementById("h50");
let miss = document.getElementById("h0");

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

let tempState;
let tempImg;
socket.onmessage = event => {
	let data = JSON.parse(event.data);
	if (tempState !== data.menuContainer.innerBG) {
		tempState = data.menuContainer.innerBG
		bg.setAttribute('src', `.././Songs/${data.menuContainer.innerBG}?a=${Math.random(10000)}`)
	}
	if (data.gameplayContainer.pp != '') {
		let ppData = JSON.parse(data.gameplayContainer.pp);
		pp.innerHTML = Math.round(ppData.pp) + "pp"
	} else {
		pp.innerHTML = 0 + "pp"
	}
	if (data.menuContainer.pp99 != '') {
		let SS = JSON.parse(data.menuContainer.pp99);
		star.innerHTML = SS.stars.toFixed(2)
	} else {
		star.innerHTML = 0
	}
	if (data.gameplayContainer[100] > 0) {
		hun.innerHTML = data.gameplayContainer[100];
	} else {
		hun.innerHTML = 0
	}
	if (data.gameplayContainer[50] > 0) {
		fifty.innerHTML = data.gameplayContainer[50];
	} else {
		fifty.innerHTML = 0
	}
	if (data.gameplayContainer.miss > 0) {
		miss.innerHTML = data.gameplayContainer.miss;
	} else {
		miss.innerHTML = 0
	}
}



//Received: '{"menuContainer":{"osuState":2,"bmID":1219126,"bmSetID":575767,"CS":4,"AR":9.5,"OD":8,"HP":6,"bmInfo":"BTS - Not Today [Tomorrow]","bmFolder":"575767 BTS - Not Today","pathToBM":"BTS - Not Today (DeRandom Otaku) [Tomorrow].osu","bmCurrentTime":8861,"bmMinBPM":0,"bmMaxBPM":0},"gameplayContainer":{"300":21,"100":0,"50":0,"miss":0,"accuracy":100,"score":24612,"combo":36,"gameMode":0,"appliedMods":2048,"maxCombo":36,"pp":""}}'