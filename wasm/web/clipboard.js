function copyCodeToClipboard() {
    const code = getCodeMirrorCode();
    navigator.clipboard.writeText(code)
        .then(() => console.log('copied to clipboard'))
        .catch((e) => console.error('failed to copy:', e));
}
