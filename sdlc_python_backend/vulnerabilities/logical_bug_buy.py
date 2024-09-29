from flask import jsonify, request
import sqlite3
import uuid

def products():

    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    cursor.execute("SELECT id, name, price FROM products")
    rows = cursor.fetchall()
    db.close()

    products = [{"id": row[0], "name": row[1], "price": row[2]} for row in rows]
    return jsonify({"status": 1, "message": "Products fetched successfully", "data": products})

def purchase():

	# // 对用户输入的购买数量不进行校验，导致购买负数商品，导致余额反而增加

# //对用户输入的购买数量不进行校验，导致购买负数商品，导致余额反而增加

    data = request.get_json()
    if not data or 'productId' not in data or 'quantity' not in data:
        return jsonify({"status": 0, "message": "无效请求"}), 400

    product_id = data['productId']
    quantity = data['quantity']

    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    username = request.username
    if not username:
        return jsonify({"status": 0, "message": "未授权"}), 401

    try:
        cursor.execute("SELECT id FROM users WHERE username = ?", (username,))
        user_id = cursor.fetchone()[0]

        cursor.execute("SELECT price FROM products WHERE id = ?", (product_id,))
        product_price = cursor.fetchone()[0]

        total_amount = product_price * quantity

        cursor.execute("SELECT balance FROM users WHERE id = ?", (user_id,))
        user_balance = cursor.fetchone()[0]

        if user_balance < total_amount:
            return jsonify({"status": 0, "message": "余额不足", "balance": user_balance})

        cursor.execute("UPDATE users SET balance = balance - ? WHERE id = ?", (total_amount, user_id))

        order_id = str(uuid.uuid4())
        cursor.execute("INSERT INTO orders (user_id, product_id, order_id, amount, quantity) VALUES (?, ?, ?, ?, ?)",
                       (user_id, product_id, order_id, total_amount, quantity))

        db.commit()

        cursor.execute("SELECT balance FROM users WHERE id = ?", (user_id,))
        updated_balance = cursor.fetchone()[0]

        return jsonify({"status": 1, "message": f"成功购买 {quantity} 个商品 ID {product_id}", "balance": updated_balance})
    except Exception as e:
        db.rollback()
        return jsonify({"status": 0, "message": f"内部服务器错误: {str(e)}"}), 500
    finally:
        db.close()

def purchase_safe():

# //对用户输入的购买数量不进行校验，若购买数量小于0，则提示错误

    data = request.get_json()
    if not data or 'productId' not in data or 'quantity' not in data:
        return jsonify({"status": 0, "message": "无效请求"}), 400

    product_id = data['productId']
    quantity = data['quantity']

    if quantity <= 0:
        return jsonify({"status": 0, "message": "购买数量必须大于0"})

    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    username = request.username
    if not username:
        return jsonify({"status": 0, "message": "未授权"}), 401

    try:
        cursor.execute("SELECT id FROM users WHERE username = ?", (username,))
        user_id = cursor.fetchone()[0]

        cursor.execute("SELECT price FROM products WHERE id = ?", (product_id,))
        product_price = cursor.fetchone()[0]

        total_amount = product_price * quantity

        cursor.execute("SELECT balance FROM users WHERE id = ?", (user_id,))
        user_balance = cursor.fetchone()[0]

        if user_balance < total_amount:
            return jsonify({"status": 0, "message": "余额不足", "balance": user_balance})

        cursor.execute("UPDATE users SET balance = balance - ? WHERE id = ?", (total_amount, user_id))

        order_id = str(uuid.uuid4())
        cursor.execute("INSERT INTO orders (user_id, product_id, order_id, amount, quantity) VALUES (?, ?, ?, ?, ?)",
                       (user_id, product_id, order_id, total_amount, quantity))

        db.commit()

        cursor.execute("SELECT balance FROM users WHERE id = ?", (user_id,))
        updated_balance = cursor.fetchone()[0]

        return jsonify({"status": 1, "message": f"成功购买 {quantity} 个商品 ID {product_id}", "balance": updated_balance})
    except Exception as e:
        db.rollback()
        return jsonify({"status": 0, "message": f"内部服务器错误: {str(e)}"}), 500
    finally:
        db.close()