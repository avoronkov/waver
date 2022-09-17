function initCodeMirror() {
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
}

function getCodeMirrorCode() {
    const doc = window.codeMirror.getDoc();
    return doc.getValue();
}

function setCodeMirrorCode(text) {
    const doc = window.codeMirror.getDoc();
    doc.setValue(text);
}
