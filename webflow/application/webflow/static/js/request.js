const makeRequest = () => {
    const form = document.createElement('form');
    form.method = 'POST';
    form.action = getHostname(window.location.href) + '/write';

    document.body.appendChild(form);

    data = {
        code: decodeURIComponent(getParameter('code')),
        role: document.getElementById('role-list').value
    }

    for(let key in data) {
        let input = document.createElement('input');
        input.type = 'hidden';
        input.name = key;
        input.value = data[key];

        form.appendChild(input);
    }

    form.submit();
}
