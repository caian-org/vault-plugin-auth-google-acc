# 3rd-party
from flask import Response
from flask import make_response
from flask import render_template


class ErrorResponse:
    def __init__(self, logger):
        self.logger = logger

    def _rd(self, code: str, msg: str, ex = None) -> Response:
        if ex is not None:
            self.logger.warn(ex)

        return make_response(render_template('error.html', code=code, message=msg.capitalize()))

    def connection_error(self, code: str, ex):
        return self._rd(code, 'unreacheable server', ex)

    def incorrectly_configured(self, code: str, ex):
        return self._rd(code, 'incorrectly configured server', ex)

    def forbidden(self, code: str, ex):
        return self._rd(code, 'incorrectly configured server', ex)

    def invalid_request(self, code: str, ex):
        return self._rd(code, 'authenticated Google account is not authorized to use the selected role', ex)
