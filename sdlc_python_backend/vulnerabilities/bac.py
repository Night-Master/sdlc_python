from flask import jsonify, request
import sqlite3

def get_profile_unauthorized():
# //对用户输入的购买数量不进行校验，导致购买负数商品，导致余额反而增加
    data = request.get_json()
    if not data or 'username' not in data:
        return jsonify({"status": 0, "message": "无效请求"}), 400

    username = data['username']

    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    cursor.execute("SELECT id, username, email, signature, avatar, birthdate, balance FROM users WHERE username = ?", (username,))
    user_row = cursor.fetchone()
    if not user_row:
        return jsonify({"status": 0, "message": "用户未找到"}), 404

    user = {
        "id": user_row[0],
        "username": user_row[1],
        "email": user_row[2],
        "signature": user_row[3],
        "avatar": user_row[4],
        "birthdate": user_row[5],
        "balance": user_row[6]
    }

    cursor.execute("SELECT o.order_id, o.amount, o.quantity, p.name FROM orders o JOIN products p ON o.product_id = p.id WHERE o.user_id = ?", (user['id'],))
    order_rows = cursor.fetchall()
    orders = []
    for order_row in order_rows:
        order = {
            "order_id": order_row[0],
            "amount": order_row[1],
            "quantity": order_row[2],
            "name": order_row[3]
        }
        orders.append(order)

    db.close()

    return jsonify({"status": 1, "message": "Profile fetched successfully", "data": {"user": user, "orders": orders}})

def get_profile():
# //对用户输入的购买数量不进行校验，若购买数量小于0，则提示错误
    data = request.get_json()
    if not data or 'username' not in data:
        return jsonify({"status": 0, "message": "无效请求"}), 400

    username = data['username']

    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    cursor.execute("SELECT id, username, email, signature, avatar, birthdate, balance FROM users WHERE username = ?", (username,))
    user_row = cursor.fetchone()
    if not user_row:
        return jsonify({"status": 0, "message": "用户未找到"}), 404

    user = {
        "id": user_row[0],
        "username": user_row[1],
        "email": user_row[2],
        "signature": user_row[3],
        "avatar": user_row[4],
        "birthdate": user_row[5],
        "balance": user_row[6]
    }

    cursor.execute("SELECT o.order_id, o.amount, o.quantity, p.name FROM orders o JOIN products p ON o.product_id = p.id WHERE o.user_id = ?", (user['id'],))
    order_rows = cursor.fetchall()
    orders = []
    for order_row in order_rows:
        order = {
            "order_id": order_row[0],
            "amount": order_row[1],
            "quantity": order_row[2],
            "name": order_row[3]
        }
        orders.append(order)

    db.close()

    return jsonify({"status": 1, "message": "Profile fetched successfully", "data": {"user": user, "orders": orders}})

def get_profile_safe():
    data = request.get_json()
    if not data or 'username' not in data:
        return jsonify({"status": 0, "message": "无效请求"}), 400

    username = data['username']
    token_username = request.username

    if token_username != username:
        return jsonify({"status": 0, "message": "您只能查看自己的账号信息"}), 403

    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    cursor.execute("SELECT id, username, email, signature, avatar, birthdate, balance FROM users WHERE username = ?", (username,))
    user_row = cursor.fetchone()
    if not user_row:
        return jsonify({"status": 0, "message": "用户未找到"}), 404

    user = {
        "id": user_row[0],
        "username": user_row[1],
        "email": user_row[2],
        "signature": user_row[3],
        "avatar": user_row[4],
        "birthdate": user_row[5],
        "balance": user_row[6]
    }

    cursor.execute("SELECT o.order_id, o.amount, o.quantity, p.name FROM orders o JOIN products p ON o.product_id = p.id WHERE o.user_id = ?", (user['id'],))
    order_rows = cursor.fetchall()
    orders = []
    for order_row in order_rows:
        order = {
            "order_id": order_row[0],
            "amount": order_row[1],
            "quantity": order_row[2],
            "name": order_row[3]
        }
        orders.append(order)

    db.close()

    return jsonify({"status": 1, "message": "Profile fetched successfully", "data": {"user": user, "orders": orders}})