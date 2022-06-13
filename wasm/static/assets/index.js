function initPage() {
    loadDefaultCode();
    // loadDefaultInstruments();
}

const initGo = async () => {
    const buffer = pako.ungzip(
        await (await fetch("demo.wasm.gz")).arrayBuffer()
    );
    const go = new Go();
    const result = await WebAssembly.instantiate(buffer, go.importObject);
    go.run(result.instance);
};
// initGo();

const updateCode = () => {
    goPause(false);
    const input = document.getElementById("code-story").value;
    const message = goPlay(input);
    document.getElementById("inst-story").value = message;
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
