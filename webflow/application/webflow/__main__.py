# 3rd-party
from flask import Flask
from flask_restful import Api as RestfulApi

# modules
from webflow.router import init_router

app = Flask(__name__, static_folder='static', template_folder='static/templates')
api = RestfulApi(app)

init_router(api)


def main() -> None:
    app.run(host='0.0.0.0', port=29747)


if __name__ == '__main__':
    main()
