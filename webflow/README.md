# Mithril Lambda

> Mithril! All folk desired it. It could be beaten like copper, and polished like
> glass; and the Dwarves could make of it a metal, light and yet harder than
> tempered steel. Its beauty was like to that of common silver, but the beauty of
> mithril did not tarnish or grow dim.

`mithril-lambda` é uma função lambda com o objetivo de facilitar o processo de
autenticação do Vault com a conta corporativa do G Suite.


## Funcionamento

A lambda expõe três rotas: `/` (root route), `/role` e `/write`. A *root
route*, quando acessada via `GET`, se conecta ao Vault e o consulta pela URL de
autenticação OAuth2 do Google, redirecionando o usuário a ela.

![](docs/oauth.png)

Finalizado o fluxo de autenticação com o Google, o usuário será enviado a URL
de callback configurada no [plugin](https://github.com/erozario/vault-auth-google).
Este redirecionamento à URL de callback, realizado via `GET`, carrega o token
de autenticação OAuth2 do Google. A URL de callback será a lambda na rota
`/role`. Esta rota exibirá ao usuário a lista das **roles** disponíveis no
Vault. Cada role é uma associação entre um grupo de usuários do G Suite e um
set de policies no Vault. A role é, em geral, nomeada a partir do nome do
time.

![](docs/role.png)

A lambda utilizará este token para "escrever" o login no Vault, gerando então o
token de autenticação do Vault que será impresso na tela.

![](docs/token.png)

A validação das permissões (cruzamento de grupos do G Suite e policies do Vault)
é realizado no processo de escrita do login no Vault.


## Requerimentos

A função requer o Python na versão `3.6`. As seguintes bibliotecas 3rd-party
foram utilizadas:

- [hvac](https://github.com/hvac/hvac): API client para o Vault;
- [zappa](https://github.com/Miserlou/Zappa): Framework de deployment de
    funções lambda com Python;
- [flask](http://flask.pocoo.org/): Microframework para aplicações web;
- [flask-restful](https://flask-restful.readthedocs.io/en/latest): Extenção do
    `flask` para construção de APIs RESTful.


## Permissões (IAM)

A documentação do `zappa` não explicita um set mínimo de permissões necessárias
para a utilização da ferramenta e [não há](https://github.com/Miserlou/Zappa/issues/244),
até o momento, um consenso a respeito.

O `deploy`, `undeploy` e `update` (atualização da lambda) foram testados
utilizando [estas permissões](docs/policy.md).


## Deployment

O deployment é realizado através do `zappa`.


```sh

# Crie o virtualenv para o deployment
virtualenv env

# Copie o arquivo de requirements para o ambiente
cp ./requirements.txt env/

# Instale as dependências
./env/bin/pip3 install -r env/requirements.txt

# Copie os arquivos da lambda e a configuração do zappa
cp -r src/* env
cp zappa.yml env

# Acesse e ative o venv
cd env
source ./bin/activate

# Caso o deploy seja local, altere as variáveis de ambiente em "zappa.yml"
sed -i -e 's@__token__@'"<TOKEN-DE-AUTENTICACAO-DO-VAULT>"'@g' zappa.yml
sed -i -e 's@__address__@'"<URL-DO-VAULT>"'@g' zappa.yml
sed -i -e 's@__auth_path__@'"<IDENTIFICADOR-DO-AUTH-METHOD>"'@g' zappa.yml

# Faça o deploy
zappa deploy prod -s zappa.yml


```


### Variáveis de ambiente

A lambda requer as seguintes variáveis de ambiente:

- `VAULT_ADDRESS`: URL do Vault.
- `VAULT_TOKEN`: Token de autenticação do Vault. Necessário para a consulta da
    URL Google e escrita do registro de login.
- `VAULT_AUTH_PATH`: Rota do Vault onde o método de autenticação é montado.
