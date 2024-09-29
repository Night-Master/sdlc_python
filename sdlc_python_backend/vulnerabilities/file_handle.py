from flask import jsonify, request, send_file
import os
import re
from werkzeug.utils import secure_filename

UPLOAD_FOLDER = './uploads'
ALLOWED_EXTENSIONS = {'.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.tiff', '.svg'}

def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

def upload_file():
#未对文件路径和文件类型进行检查
    if 'file' not in request.files:
        return jsonify({"status": 0, "message": "文件获取失败"}), 400

    file = request.files['file']
    if file.filename == '':
        return jsonify({"status": 0, "message": "文件名不能为空"}), 400

    filename = secure_filename(file.filename)
    file_path = os.path.join(UPLOAD_FOLDER, filename)
    file.save(file_path)

    return jsonify({"status": 1, "message": "文件上传成功", "filePath": file_path})

def upload_file_safe():
#对上传的文件进行上传路径和文件类型进行检查
    if 'file' not in request.files:
        return jsonify({"status": 0, "message": "文件获取失败"}), 400

    file = request.files['file']
    if file.filename == '':
        return jsonify({"status": 0, "message": "文件名不能为空"}), 400

    ext = os.path.splitext(file.filename)[1].lower()
    if ext not in ALLOWED_EXTENSIONS:
        return jsonify({"status": 0, "message": f"不允许的文件类型: {ext}"}), 400

    filename = secure_filename(file.filename)
    filename = re.sub(r'[\\/]', '', filename)  # 移除所有非法字符
    file_path = os.path.join(UPLOAD_FOLDER, filename)
    file.save(file_path)

    return jsonify({"status": 1, "message": "文件上传成功", "filePath": file_path})

def list_images():
    files = os.listdir(UPLOAD_FOLDER)
    images = [file for file in files if os.path.isfile(os.path.join(UPLOAD_FOLDER, file))]
    return jsonify({"status": 1, "message": "获取图片列表成功", "images": images})

def download():
#未对下载的文件的路径和文件类型进行检查
    data = request.get_json()
    if not data or 'fileName' not in data:
        return jsonify({"status": 0, "message": "无效的请求"}), 400

    file_name = data['fileName']
    file_path = os.path.join(UPLOAD_FOLDER, file_name)

    if not os.path.exists(file_path):
        return jsonify({"status": 0, "message": "文件不存在"}), 404

    return send_file(file_path, as_attachment=True)

def download_safe():
#对下载的文件的路径和文件类型进行检查
    data = request.get_json()
    if not data or 'fileName' not in data:
        return jsonify({"status": 0, "message": "无效的请求"}), 400

    file_name = data['fileName']
    if re.search(r'[\\/]', file_name):
        return jsonify({"status": 0, "message": "文件名包含非法字符"}), 400

    ext = os.path.splitext(file_name)[1].lower()
    if ext not in ALLOWED_EXTENSIONS:
        return jsonify({"status": 0, "message": "不允许的文件类型"}), 400

    file_path = os.path.join(UPLOAD_FOLDER, file_name)
    abs_path = os.path.abspath(file_path)
    expected_dir = os.path.abspath(UPLOAD_FOLDER)

    if not abs_path.startswith(expected_dir):
        return jsonify({"status": 0, "message": "文件路径非法"}), 400

    if not os.path.exists(abs_path):
        return jsonify({"status": 0, "message": "文件不存在"}), 404

    return send_file(abs_path, as_attachment=True)