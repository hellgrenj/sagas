import { messages } from "../stores/messages";

let msgs = [];
export const connectAndConsume = () => {
    const ws = new WebSocket("ws://localhost:8080/ws");
    let prevEv = null;
    let currentColor = getRandomColor();

    ws.onopen = () => {
        console.log('successfully connected...');
    }
    ws.onclose = (e) => {
        console.log('connection closed, reconnecting...');
        connectAndConsume();
    };
    ws.onmessage = function (e) {
        const ev = JSON.parse(e.data);

        if (prevEv == null) {
            addNewMsgs({
                Header: true,
                CorrelationId: ev.CorrelationId,
                Color: currentColor,
            });
        }
        if (prevEv && prevEv.CorrelationId !== ev.CorrelationId) {
            currentColor = getRandomColor();
            addNewMsgs({
                Header: true,
                CorrelationId: ev.CorrelationId,
                Color: currentColor,
            });
        }
        ev.Color = currentColor;
        prevEv = ev;
        addNewMsgs(ev);

        msgs = msgs;
    };
};

function getRandomColor() {
    const letters = "0123456789ABCDEF";
    let color = "#";
    for (let i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}
const uniqueCorrelationIds = [];
function addNewMsgs(ev) {
    if (uniqueCorrelationIds.length == 6) {
        const corrIdToRemove = uniqueCorrelationIds.shift();
        const newMsgs = msgs.filter((m) => m.CorrelationId != corrIdToRemove);
        newMsgs.push(ev);
        msgs = newMsgs;
        messages.set(msgs)
    } else {
        msgs.push(ev);
        msgs = msgs;
        messages.set(msgs)
    }
    if (!uniqueCorrelationIds.includes(ev.CorrelationId)) {
        uniqueCorrelationIds.push(ev.CorrelationId);
    }
}