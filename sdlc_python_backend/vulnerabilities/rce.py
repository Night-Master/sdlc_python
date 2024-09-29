from flask import jsonify, request
import subprocess
import platform

def contains(slice, item):
    return item in slice

def execute_command():
#对用户输入的命令参数检测策略不严格，只要包含了指定命令即可，可以使用连接符绕过
    data = request.get_json()
    if not data or 'command' not in data:
        return jsonify({"status": 0, "message": "无效请求"}), 400

    command = data['command']

    # 检查命令是否符合指定命令
    valid_commands = ["dir", "ls", "ipconfig", "ifconfig"]
    input_command = command.split(" ")[0]
    if input_command not in valid_commands:
        return jsonify({"status": 0, "message": "你输入的不是指定的命令"}), 400

    # 根据操作系统执行命令
    try:
        if platform.system() == "Windows":
            output = subprocess.check_output(["cmd", "/C", command], stderr=subprocess.STDOUT, text=True)
        else:
            output = subprocess.check_output(["sh", "-c", command], stderr=subprocess.STDOUT, text=True)
    except subprocess.CalledProcessError as e:
        return jsonify({"status": 0, "message": f"命令执行失败: {e.output.strip()}"}), 400

    # 返回命令执行结果
    return jsonify({"status": 1, "message": output.strip()})

def execute_command_safe():
# // 白名单策略：使用 valid_commands 映射来定义允许的命令及其参数。例如，dir、ls、ipconfig 和 ifconfig 命令没有额外的参数。
# // 命令参数验证：检查命令参数是否在白名单中。如果命令参数不在白名单中，则返回错误信息。
# // 通过这种方式，我们可以确保只有白名单中的命令和参数可以被执行，从而大大降低命令注入的风险。
    data = request.get_json()
    if not data or 'command' not in data:
        return jsonify({"status": 0, "message": "无效请求"}), 400

    command = data['command']

    # 检查命令是否符合指定命令
    valid_commands = {
        "dir": [],
        "ls": [],
        "ipconfig": [],
        "ifconfig": [],
    }

    cmd_parts = command.split(" ")
    input_command = cmd_parts[0]
    if input_command not in valid_commands:
        return jsonify({"status": 0, "message": "你输入的不是指定的命令"}), 200

    # 检查命令参数是否在白名单中
    if len(cmd_parts) > 1:
        for arg in cmd_parts[1:]:
            if arg not in valid_commands[input_command]:
                return jsonify({"status": 0, "message": "命令参数不在白名单中"}), 200

    # 根据操作系统执行命令
    try:
        if platform.system() == "Windows":
            output = subprocess.check_output(["cmd", "/C", command], stderr=subprocess.STDOUT, text=True)
        else:
            output = subprocess.check_output(["sh", "-c", command], stderr=subprocess.STDOUT, text=True)
    except subprocess.CalledProcessError as e:
        return jsonify({"status": 0, "message": f"命令执行失败: {e.output.strip()}"}), 400

    # 返回命令执行结果
    return jsonify({"status": 1, "message": output.strip()})