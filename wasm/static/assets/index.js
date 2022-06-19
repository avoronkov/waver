function initPage() {
    loadDefaultCode();
}

function logMessage(msg) {
    const msgArea = document.getElementById("inst-story");
    const value = msgArea.value.trim();
    msgArea.value = (value ? value + "\n" : value) + `${msg}`;
}

const init = () => {
    initCodeMirror();
    initGo();
};

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

const initCodeMirror = () => {
    CodeMirror.defineSimpleMode('waver', {
        start: [
            { regex: /(?:tempo|sample|inst)\b/, token: 'comment' }, // keyword
            { regex: /(?:min|maj|min7|maj7|min9|maj9)\b/, token: 'variable-2' },
            { regex: /(?:seq|rand|repeat|_dur|_)\b/, token: 'variable-2' },
            { regex: /".*"/, token: 'string' },
            { regex: /'.*'/, token: 'string' },
            { regex: /#.*/, sol: true, token: 'meta' },
            { regex: /\b(?:[\d]+(\.[\d]*)?)\b/, token: 'number' },
            { regex: /->/, token: 'comment' },
            { regex: /(% |%%)/, sol: true, token: 'comment' },
            { regex: /[+=\-:<>]/, token: 'atom' },
            // notes
            { regex: /\b[ABCDEFG][sb]?\d\b/, token: 'keyword' },
            // filters tokens (variable-3)
            { regex: /(?:adsr|delay|dist|distortion|vibrato|am|timeshift|harmonizer|harm|flanger|exp|movexp|ration|swingexp)\b/, token: 'variable-3' },
            { regex: /(?:attackLevel|decayLevel|attackLen|decayLen|sustainLen|releaseLen)\b/, token: 'variable-3' },
            { regex: /(?:int|interval|times|fade|value)\b/, token: 'variable-3'},
            { regex: /(?:freq|frequency|shift|amp|amplitude|wave)\b/, token: 'variable-3'},
        ],
    });

    const ta = document.getElementById('code-story')
    const codeMirror = CodeMirror.fromTextArea(ta, {
                    mode:  "waver",
                    lineNumbers: true,
    });
    codeMirror.setSize(null, 590);
    
    window.codeMirror = codeMirror;
    return codeMirror;
};

const updateCode = () => {
    try {
        logMessage('Updating code...');
        goPause(false);
        const doc = window.codeMirror.getDoc();
        const input = doc.getValue();
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
    const doc = window.codeMirror.getDoc();
    doc.setValue(code);
    document.getElementById('update-code').disabled = false;
};

const clearCode = () => {
    const doc = window.codeMirror.getDoc();
    doc.setValue('');
};

const makeSharedLink = () => {
    const doc = window.codeMirror.getDoc();
    const code = doc.getValue();
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
