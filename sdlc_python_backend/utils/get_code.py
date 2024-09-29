import sqlite3
from flask import request, jsonify

def get_api_code():
    data = request.get_json()
    if not data or 'api_name' not in data:
        return jsonify({"status": 0, "message": "无效请求"}), 400

    api_name = data['api_name']

    # 连接到 SQLite 数据库
    db = sqlite3.connect('test.db')
    cursor = db.cursor()

    # 查询 API 的 code
    cursor.execute("SELECT function_name, code FROM api_info WHERE api_name = ?", (api_name,))
    api_row = cursor.fetchone()

    if not api_row:
        db.close()
        return jsonify({"status": 0, "message": "API 未找到"}), 404

    api_info = {
        "function_name": api_row[0],
        "code": api_row[1]
    }

    db.close()

    return jsonify({"status": 1, "message": "API 代码查询成功", "data": api_info})
