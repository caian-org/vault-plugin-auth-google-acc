# 3rd-party
from flask import Response
from flask import make_response
from flask import render_template


class ErrorResponse:
    def _rd(self, code: str, msg: str) -> Response:
        return make_response(render_template('error.html', code=code, message=msg.capitalize()))

    def connection_error(self, code: str):
        return self._rd(code, 'unreacheable server')

    def incorrectly_configured(self, code: str):
        return self._rd(code, 'incorrectly configured server')

    def forbidden(self, code: str):
        return self._rd(code, 'incorrectly configured server')

    def invalid_request(self, code: str):
        return self._rd(code, 'authenticated Google account is not authorized to use the selected role')
