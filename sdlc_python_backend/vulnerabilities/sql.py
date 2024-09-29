import sqlite3
from flask import request, jsonify
from utils import generate_token
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.backends import default_backend
from Crypto.PublicKey import RSA
import base64

# 生成RSA密钥对
private_key = rsa.generate_private_key(
    public_exponent=65537,
    key_size=2048,
)

# 将公钥导出为 PKCS#1 格式
public_key_pem = private_key.public_key().public_bytes(
    encoding=serialization.Encoding.PEM,
    format=serialization.PublicFormat.PKCS1  # 使用PKCS1格式
)

def encrypt(data):
    ciphertext = private_key.public_key().encrypt(
        data.encode('utf-8'),  # 确保数据是字节串
        padding.PKCS1v15()
    )
    return ciphertext

def decrypt(ciphertext):
    # 首先将Base64编码的字符串解码为字节串
    ciphertext_bytes = base64.b64decode(ciphertext)
    
    # 解密字节串
    plaintext = private_key.decrypt(
        ciphertext_bytes,
        padding.PKCS1v15()
    )
    return plaintext.decode('utf-8')  # 解密后返回字符串

def get_public_key():
    return public_key_pem.decode('utf-8')

def sql_injection_sqlite3():
#拼接sql导致sql注入
    data = request.json
    username = data.get('username')
    password = data.get('password')
    
    db = sqlite3.connect('test.db')
    cursor = db.cursor()
    query = f"SELECT * FROM users WHERE username='{username}' AND password='{password}'"
    cursor.execute(query)
    rows = cursor.fetchall()
    db.close()

    if rows:
        return jsonify({"status": 1, "message": "Login successful!"})
    else:
        return jsonify({"status": 0, "message": "Login failed!"})

def sql_injection_sqlite3_safe():
#参数化查询，它通过将用户输入的数据作为参数传递给查询，而不是直接将用户输入嵌入到查询字符串中，从而有效地防止SQL注入攻击
    data = request.get_json()
    
    # 假设请求中的用户名和密码是加密的，需要解密
    encrypted_username = data.get('username')
    encrypted_password = data.get('password')

    # 解密用户名和密码
    try:
        username = decrypt(encrypted_username)
        password = decrypt(encrypted_password)
    except Exception as e:
        return jsonify({"status": 0, "message": "解密失败"}), 400

    # 使用解密后的用户名和密码进行数据库操作
    conn = sqlite3.connect('test.db')
    cursor = conn.cursor()
    query = "SELECT * FROM users WHERE username=? AND password=?"
    cursor.execute(query, (username, password))
    result = cursor.fetchone()
    
    if result:
        token = generate_token(username)
        return jsonify({"status": 1, "message": "Login successful!", "token": token})
    else:
        return jsonify({"status": 0, "message": "Login failed!"})

def change_password_safe():
#服务端先生成一对rsa公私钥，把公钥发给前端加密，后端使用私钥解密
    data = request.get_json()
    if not data or 'currentPassword' not in data or 'newPassword' not in data:
        return jsonify({"status": 0, "message": "无效请求"}), 400

    encrypted_current_password = data['currentPassword']
    encrypted_new_password = data['newPassword']

    try:
        current_password = decrypt(encrypted_current_password)
        new_password = decrypt(encrypted_new_password)
    except Exception as e:
        return jsonify({"status": 0, "message": "解密失败"}), 400

    username = request.username
    if not username:
        return jsonify({"status": 0, "message": "未授权"}), 401

    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    try:
        cursor.execute("SELECT * FROM users WHERE username=? AND password=?", (username, current_password))
        if not cursor.fetchone():
            return jsonify({"status": 0, "message": "当前密码不正确"}), 400

        cursor.execute("UPDATE users SET password=? WHERE username=? AND password=?", (new_password, username, current_password))
        db.commit()

        return jsonify({"status": 1, "message": "密码更改成功"})
    except Exception as e:
        db.rollback()
        return jsonify({"status": 0, "message": f"内部服务器错误: {str(e)}"}), 500
    finally:
        db.close()

def change_password_unsafe():
#如果没有使用https保护，且敏感信息不使用任何加密，则构成明文传输漏洞
    data = request.get_json()
    if not data or 'currentPassword' not in data or 'newPassword' not in data:
        return jsonify({"status": 0, "message": "无效请求"}), 400

    current_password = data['currentPassword']
    new_password = data['newPassword']

    username = request.username
    if not username:
        return jsonify({"status": 0, "message": "未授权"}), 401

    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    try:
        cursor.execute("SELECT * FROM users WHERE username=? AND password=?", (username, current_password))
        if not cursor.fetchone():
            return jsonify({"status": 0, "message": "当前密码不正确"}), 400

        cursor.execute("UPDATE users SET password=? WHERE username=? AND password=?", (new_password, username, current_password))
        db.commit()

        return jsonify({"status": 1, "message": "密码更改成功"})
    except Exception as e:
        db.rollback()
        return jsonify({"status": 0, "message": f"内部服务器错误: {str(e)}"}), 500
    finally:
        db.close()