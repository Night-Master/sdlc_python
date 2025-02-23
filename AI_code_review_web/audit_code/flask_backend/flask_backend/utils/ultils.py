import jwt
import time

jwt_secret = "sldc_python"

def generate_token(username):
    now_time = time.time()
    expire_time = now_time + 3 * 60 * 60
    claims = {
        "username": username,
        "exp": expire_time,
        "iss": "your_app_name"
    }
    token = jwt.encode(claims, jwt_secret, algorithm="HS256")
    return token

def parse_token(token):
    try:
        claims = jwt.decode(token, jwt_secret, algorithms=["HS256"])
        return claims, None
    except jwt.ExpiredSignatureError:
        return None, "Token expired"
    except jwt.InvalidTokenError:
        return None, "Invalid token"