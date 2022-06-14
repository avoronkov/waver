function initPage() {
    loadDefaultCode();
    // loadDefaultInstruments();
}

function logMessage(msg) {
    const msgArea = document.getElementById("inst-story");
    const value = msgArea.value.trim();
    msgArea.value = (value ? value + "\n" : value) + `${msg}`;
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
    const params = new URLSearchParams(window.location.search);
    let code = undefined;
    if (params.has('code')) {
        const encoded = params.get('code');
        const { data, error } = goDecode(encoded);
        if (error) {
            logMessage(`FAILED: ${error}`);
            return;
        }
        code = data
    } else {
        code = goGetDefaultCode();
    }
    document.getElementById('code-story').value = code;
    document.getElementById('update-code').disabled = false;
};

const clearCode = () => {
    document.getElementById("code-story").value = '';
};

const makeSharedLink = () => {
    const code = document.getElementById('code-story').value;
    const { data, error } = goEncode(code);
    if (error) {
        logMessage(`FAILED: ${error}`);
        return;
    }
    const { protocol, host, pathname} = window.location;
    const baseUrl = `${protocol}//${host}${pathname}`;
    const link = `${baseUrl}?code=${data}`;
    window.prompt("Copy link to clipboard: Ctrl+C, Enter", link);
}
