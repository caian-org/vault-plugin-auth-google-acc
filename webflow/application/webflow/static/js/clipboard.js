const copyButton = document.getElementById('copy-btn');

copyButton.addEventListener('click', event => {
    const vaultTokenTextbox = document.getElementById('token-txtbox');
    vaultTokenTextbox.focus();
    vaultTokenTextbox.select();

    document.execCommand('copy');
});
