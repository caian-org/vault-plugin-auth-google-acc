# standard
import logging

# 3rd-party
from flask import redirect
from flask import make_response
from flask import render_template
from flask import Response
from flask_restful import Resource

# modules
from webflow.resource._shared.service import Service

from .error import ErrorResponse


class Controller(Resource):
    def __init__(self):
        self.logger = logging.getLogger(__name__)

        self.error = ErrorResponse()
        self.service = Service()

    @staticmethod
    def redirect_to(url: str) -> Response:
        return redirect(url)

    @staticmethod
    def render_page(page: str, **kwargs) -> Response:
        return make_response(render_template(page, **kwargs))
