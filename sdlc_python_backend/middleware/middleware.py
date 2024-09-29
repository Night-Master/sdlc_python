from flask import request, jsonify
from utils import parse_token

def auth_middleware(view_func):
    def wrapper(*args, **kwargs):
        token = request.headers.get('Authorization')
        if not token:
            return jsonify({"status": 0, "message": "您未登录，请勿未授权访问！"}), 200

        claims, err = parse_token(token)
        if err:
            return jsonify({"status": 0, "message": "登录凭证过期"}), 200

        request.username = claims['username']
        return view_func(*args, **kwargs)
    return wrapper