function getPayloadData () {
  const queryParams = new URLSearchParams(window.location.search);

  return {
    code: decodeURIComponent(queryParams.get('code')),
    role: document.getElementById('role-list').value
  }
}

function makeRequest () {
  const form = document.createElement('form');
  form.method = 'POST';
  form.action = '/login';

  const data = getPayloadData();
  Object.keys(data).forEach((key) => {
    const input = document.createElement('input');
    input.type = 'hidden';
    input.name = key;
    input.value = data[key];

    form.appendChild(input);
  })

  document.body.appendChild(form);
  form.submit();
}
