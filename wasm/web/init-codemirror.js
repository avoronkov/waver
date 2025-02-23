// This file was generated by syntaxgen.

function initCodeMirror() {
    CodeMirror.defineSimpleMode('waver', {
        start: [
            { regex: /(?:filter|form|inst|lagrange|sample|scale|srand|stop|tempo|wave)\b/, token: 'comment' }, // keyword
            { regex: /(?:noise|saw|semisine|sine|square|triangle)\b/, token: 'comment' }, // keyword
            { regex: /(?:maj|maj7|maj9|min|min7|min9)\b/, token: 'variable-2' },
            { regex: /(?:concat|down|loop|rand|repeat|seq|up)\b/, token: 'variable-2' },
            { regex: /(?:bits|eucl|eucl')\b/, token: 'variable-3' },
            { regex: /(?:_|_dur|true|false)\b/, token: 'variable-3' },
            { regex: /".*"/, token: 'string' },
            { regex: /'.*'/, token: 'string' },
            { regex: /#.*/, sol: true, token: 'meta' },
            { regex: /\b(?:[\d]+(\.[\d]*)?)\b/, token: 'number' },
            { regex: /->/, token: 'comment' },
            { regex: /(% |%%)/, sol: true, token: 'comment' },
            // TODO generate atoms with syntaxgen
            { regex: /[+=\-:<>]/, token: 'atom' },
            // notes
            { regex: /\b[ABCDEFG][sb]?\d\b/, token: 'keyword' },
            // filters tokens (variable-3)
            { regex: /(?:8bit|adsr|am|delay|dist|distortion|exp|flanger|movexp|pan|ratio|swingexp|swingpan|timeshift|vibrato)\b/, token: 'variable-3' },
            { regex: /(?:abs|abstime|amp|amplitude|attackLen|attackLevel|attacklen|attacklevel|bits|carrier|carrierctx|decayLen|decayLevel|decaylen|decaylevel|fade|fadeout|feedback|freq|frequency|initialValue|initialvalue|int|interval|inverse|l|left|maxShift|maxshift|r|releaseLen|releaselen|right|shifter|shifterctx|speed|sustainLen|sustainlen|tempo|times|value)\b/, token: 'variable-3' },
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
