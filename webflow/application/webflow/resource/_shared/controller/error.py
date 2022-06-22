# 3rd-party
from flask import Response
from flask import make_response
from flask import render_template


class ErrorResponse:
    def _renderr(self, msg: str) -> Response:
        return make_response(render_template('error.html', message=msg))

    @property
    def connection_error(self):
        return self._renderr('Could not connect to Vault server')

    @property
    def forbidden(self):
        return self._renderr('Vault refused connection')

    @property
    def invalid_request(self):
        return self._renderr('Authenticated Google account is not authorized to use the selected role')
