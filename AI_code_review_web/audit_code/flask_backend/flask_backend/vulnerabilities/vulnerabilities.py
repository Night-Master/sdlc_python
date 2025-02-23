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

print(public_key_pem.decode('utf-8'))


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
    # print(request.username)
    data = request.json
    username = data.get('username')
    password = data.get('password')
        # username = request.username
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




