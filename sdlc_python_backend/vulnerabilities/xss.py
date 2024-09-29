from flask import jsonify, request
import html
import sqlite3

def reflect_xss():
#未对接受的参数进行处理
    data = request.get_json()
    if not data or 'input' not in data:
        return jsonify({"status": 0, "message": "Invalid request"}), 400

    input_data = data['input']
    return jsonify({"status": 1, "message": input_data})

def reflect_xss_safe():
#对接受的参数使用 html.escape转义
    data = request.get_json()
    if not data or 'input' not in data:
        return jsonify({"status": 0, "message": "Invalid request"}), 400

    input_data = html.escape(data['input'])
    return jsonify({"status": 1, "message": input_data})

def get_comments():
    db = sqlite3.connect('test.db')
    cursor = db.cursor()
    cursor.execute("SELECT username, content FROM comments")
    rows = cursor.fetchall()
    db.close()

    comments = [{"username": row[0], "content": row[1]} for row in rows]
    return jsonify({"status": 1, "message": "Comments fetched successfully", "data": comments})

def create_comments():
#未对接收上来的Content字段进行转义处理
    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    data = request.get_json()
    if not data or 'content' not in data:
        return jsonify({"status": 0, "message": "Invalid request"}), 400

    username = request.username
    if not username:
        return jsonify({"status": 0, "message": "Unauthorized"}), 401

    cursor.execute("INSERT INTO comments (username, content) VALUES (?, ?)", (username, data['content']))
    db.commit()
    db.close()

    return jsonify({"status": 1, "message": "Comment added successfully"})

def create_comments_safe():
#对接收上来的Content字段进行html.escape转义处理
    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    data = request.get_json()
    if not data or 'content' not in data:
        return jsonify({"status": 0, "message": "Invalid request"}), 400

    username = request.username
    if not username:
        return jsonify({"status": 0, "message": "Unauthorized"}), 401

    escaped_content = html.escape(data['content'])
    cursor.execute("INSERT INTO comments (username, content) VALUES (?, ?)", (username, escaped_content))
    db.commit()
    db.close()

    return jsonify({"status": 1, "message": "Comment added successfully"})

def clear_comments():
    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    username = request.username
    if not username or username != 'admin':
        return jsonify({"status": 0, "message": "您不是管理员"}), 200

    cursor.execute("DELETE FROM comments")
    db.commit()
    db.close()

    return jsonify({"status": 1, "message": "All comments cleared successfully"})