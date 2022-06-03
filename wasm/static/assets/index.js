function initPage() {
    loadDefaultCode();
    loadDefaultInstruments();
}

const initGo = async () => {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(
        fetch("demo.wasm"),
        go.importObject
    );
    go.run(result.instance);
};
// initGo();

const updateCode = () => {
    goPause(false);
    const input = document.getElementById("code-story").value;
    goPlay(input);
};

const loadDefaultCode = () => {
    const code = goGetDefaultCode();
    document.getElementById("code-story").value = code;
};

const updateInstruments = () => {
    goPause(false);
    const input = document.getElementById("inst-story").value;
    goUpdateInstruments(input);
};

const loadDefaultInstruments = () => {
    const inst = goGetDefaultInstruments();
    document.getElementById("inst-story").value = inst;
};
