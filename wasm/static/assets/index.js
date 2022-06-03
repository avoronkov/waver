function initPage() {
    console.log("initPage...");
    console.log("initPage continue");
    loadDefaultCode();
    console.log("initPage OK");
}

const initGo = async () => {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(
        fetch("demo.wasm"),
        go.importObject
    );
    go.run(result.instance);
};
initGo();

const updateCode = () => {
    const input = document.getElementById("story").value;
    goPlay(input);
};

const loadDefaultCode = () => {
    const code = goGetDefaultCode();
    document.getElementById("story").value = code;
};

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
