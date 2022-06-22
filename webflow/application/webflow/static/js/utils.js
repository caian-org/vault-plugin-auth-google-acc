const getHostname = url => {
    const c = url.split('/');

    if(url.indexOf('//') > -1) {
        const protocol = c[0],
            domain = c[2];

        return `${protocol}//${domain}`;
    }
    return c[0];
}

const getParameter = val => {
    let params = {};
    location.search.replace(/[?&]+([^=&]+)=([^&]*)/gi,
        (_, key, value) => {
            params[key] = value;
        })
    return val ? params[val] : params;
}
