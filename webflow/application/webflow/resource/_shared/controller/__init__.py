# standard
import logging

# 3rd-party
from flask import make_response
from flask import request
from flask import redirect
from flask import render_template
from flask_restful import Resource

# modules
from webflow.resource._shared.service import Service

from .error import ErrorResponse


class Controller(Resource):
    def __init__(self):
        self.logger = logging.getLogger(__name__)

        self.error = ErrorResponse(self.logger)
        self.service = Service()

    @staticmethod
    def redirect_to(url: str):
        return redirect(url)

    @staticmethod
    def render_page(page: str, **kwargs):
        return make_response(render_template(page, **kwargs))

    @staticmethod
    def fetch_form_val(name: str) -> str:
        return request.form.get(name)
