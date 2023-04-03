// Basic calls to interact with the "car" RESTCONF API

let car = null;

// Reading config AND metrics together for the whole car
function update() {
    fetch("/restconf/data/car:").then((resp) => {
        return resp.json();
    }).then((data) => {
        car = data;
        let e = document.getElementById("car");
        e.innerHTML = `
            <div>Miles: ${data.miles} miles</div>
            <div>Running: ${data.running}</div>
            <div>Speed: ${data.speed} mps</div>
            ${data.tire.map((tire, index) =>
                `<div>Tire ${index}: Flat: ${tire.flat}, Wear: ${tire.wear}</div>`
            ).join('')}
        `;
    });
};

// Call one of the RPCs. ours are simple, no input or expected output but
// you could easily add them here
function run(action) {
    fetch(`/restconf/data/car:${action}`, {
        method: 'POST',
    }).then(() => {
        update();
    });
}       

// Update config edit.  
//   PATCH is like an upsert. 
//   PUT is like a delete and replace with this.
//   POST is create new
function speed(s) {
    fetch('/restconf/data/car:', {
        method: 'PATCH',
        body : JSON.stringify({
                speed: car.speed + s
            })
    }).then(() => {
        update();
    });                
}

// Subscribe to notifications. EventSource is part of every browser and works 
// on HTTP (i.e. http) but more efficient on HTTP/2 (i.e. https)
function subscribeToUpdateEventStream() {
    // this will appear as a pending request.
    const events = new EventSource("/restconf/data/car:update?simplified");
	events.addEventListener("message", e => {
        // this decodes SSE events for you to give you just the messages
        const log = document.getElementById("log");
        log.appendChild(document.createTextNode(e.data + '\n'));
	});
    // to unsubscribe just close the stream.
    //  events.close();
}

// Metrics change too often for subscription so just get the latest
// every 2s.
const pollInterval = 2000;
function pollforUpdates()  {
    update();
    setTimeout(pollforUpdates, pollInterval);
}

// initial read of data
pollforUpdates();
// watch for update events as they happen, no polling here
subscribeToUpdateEventStream();
