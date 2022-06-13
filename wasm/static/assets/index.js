function initPage() {
    loadDefaultCode();
    // loadDefaultInstruments();
}

function logMessage(msg) {
    const msgArea = document.getElementById("inst-story");
    msgArea.value = msgArea.value + "\n" + `${msg}`;
}

const initGo = async () => {
    try {
        const buffer = pako.ungzip(
            await (await fetch("demo.wasm.gz")).arrayBuffer()
        );
        const go = new Go();
        const result = await WebAssembly.instantiate(buffer, go.importObject);
        await go.run(result.instance);
    } catch (e) {
        logMessage(`FAILED running Go: ${e}`);
    }
};
// initGo();

const updateCode = () => {
    try {
        logMessage('Updating code...');
        goPause(false);
        const input = document.getElementById("code-story").value;
        const message = goPlay(input);
        logMessage(message);
    } catch (e) {
        logMessage(`FAILED: ${e}`);
    }
};

const loadDefaultCode = () => {
    const code = goGetDefaultCode();
    document.getElementById("code-story").value = code;
    document.getElementById("update-code").disabled = false;
};
/*
const updateInstruments = () => {
    goPause(false);
    const input = document.getElementById("inst-story").value;
    goUpdateInstruments(input);
};

const loadDefaultInstruments = () => {
    const inst = goGetDefaultInstruments();
    document.getElementById("inst-story").value = inst;
};
*/
