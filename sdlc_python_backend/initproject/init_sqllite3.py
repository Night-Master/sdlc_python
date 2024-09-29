import ast
import os
import sqlite3
import logging
from typing import Dict, List, Optional
from flask import jsonify

logging.basicConfig(level=logging.INFO)

def initialize_database():
    """
    初始化 SQLite3 数据库，创建必要的表，解析 Flask 路由文件，
    提取 API 信息和函数代码，并将其存储到数据库中。同时插入示例数据。
    """

    # 定义数据库路径
    parent_dir = os.path.dirname(os.path.dirname(__file__))
    db_path = os.path.join(parent_dir, "test.db")
    
    # 定义路由文件路径和 vulnerabilities/__init__.py 路径
    route_file_path = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "routes", "routes.py"))
    vuln_init_file = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "vulnerabilities", "__init__.py"))

    # 连接到 SQLite3 数据库
    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        logging.info("成功连接到 SQLite3 数据库。")
    except sqlite3.Error as e:
        logging.error(f"无法连接到数据库: {e}")
        return jsonify({"status": "error", "message": f"无法连接到数据库: {e}"}), 500

    # 删除原有表
    try:
        cursor.execute("DROP TABLE IF EXISTS users")
        cursor.execute("DROP TABLE IF EXISTS products")
        cursor.execute("DROP TABLE IF EXISTS orders")
        cursor.execute("DROP TABLE IF EXISTS comments")
        cursor.execute("DROP TABLE IF EXISTS api_info")
        logging.info("已删除原有表。")
    except sqlite3.Error as e:
        logging.error(f"删除表时出错: {e}")
        conn.close()
        return jsonify({"status": "error", "message": f"删除表时出错: {e}"}), 500

    # 创建表
    try:
        cursor.execute("""
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL,
            password TEXT NOT NULL,
            email TEXT NOT NULL,
            signature TEXT,
            avatar TEXT,
            birthdate TEXT,
            balance REAL
        )
        """)
        logging.info("已创建 'users' 表。")

        cursor.execute("""
        CREATE TABLE products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            price REAL NOT NULL
        )
        """)
        logging.info("已创建 'products' 表。")

        cursor.execute("""
        CREATE TABLE orders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            product_id INTEGER NOT NULL,
            order_id TEXT NOT NULL,
            amount REAL NOT NULL,
            quantity INTEGER NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users (id),
            FOREIGN KEY (product_id) REFERENCES products (id)
        )
        """)
        logging.info("已创建 'orders' 表。")

        cursor.execute("""
        CREATE TABLE comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL,
            content TEXT NOT NULL
        )
        """)
        logging.info("已创建 'comments' 表。")

        cursor.execute("""
        CREATE TABLE api_info (
            api_name TEXT NOT NULL,
            function_name TEXT NOT NULL,
            code TEXT,
            PRIMARY KEY (api_name, function_name)
        )
        """)
        logging.info("已创建 'api_info' 表。")

        cursor.execute("CREATE INDEX idx_api_name ON api_info (api_name)")
        cursor.execute("CREATE INDEX idx_function_name ON api_info (function_name)")
        logging.info("已创建索引 'idx_api_name' 和 'idx_function_name'。")

        conn.commit()
    except sqlite3.Error as e:
        logging.error(f"创建表时出错: {e}")
        conn.close()
        return jsonify({"status": "error", "message": f"创建表时出错: {e}"}), 500

    # 插入示例数据
    try:
        cursor.executemany("""
        INSERT INTO users (username, password, email, signature, avatar, birthdate, balance) 
        VALUES (?, ?, ?, ?, ?, ?, ?)
        """, [
            ("user1", "hello", "123@qq.com", "Hello, I'm user1!", "avatar1.jpg", "1990-01-01", 100.0),
            ("user2", "password1", "1223@qq.com", "Hello, I'm user2!", "avatar2.jpg", "1995-05-05", 200.0)
        ])
        logging.info("已插入示例数据到 'users' 表。")

        cursor.executemany("""
        INSERT INTO products (name, price) VALUES (?, ?)
        """, [
            ("Product A", 50.0),
            ("Product B", 75.0)
        ])
        logging.info("已插入示例数据到 'products' 表。")

        def generate_order_id():
            import random
            return ''.join(random.choices('0123456789', k=16))

        order_ids = [generate_order_id() for _ in range(2)]
        cursor.executemany("""
        INSERT INTO orders (user_id, product_id, order_id, amount, quantity) VALUES (?, ?, ?, ?, ?)
        """, [
            (1, 1, order_ids[0], 50.0, 1),
            (2, 2, order_ids[1], 75.0, 1)
        ])
        logging.info("已插入示例数据到 'orders' 表。")

        cursor.executemany("""
        INSERT INTO comments (username, content) VALUES (?, ?)
        """, [
            ("user1", "hello"),
            ("user2", "我要学渗透")
        ])
        logging.info("已插入示例数据到 'comments' 表。")

        conn.commit()
        logging.info("示例数据已成功插入到各个表中。")
    except sqlite3.Error as e:
        logging.error(f"插入示例数据时出错: {e}")
        conn.rollback()

    # 解析路由文件
    if not os.path.isfile(route_file_path):
        logging.error(f"路由文件不存在: {route_file_path}")
        conn.close()
        return jsonify({"status": "error", "message": f"路由文件不存在: {route_file_path}"}), 404

    with open(route_file_path, "r", encoding="utf-8") as f:
        route_source = f.read()

    try:
        route_tree = ast.parse(route_source, filename=route_file_path)
        logging.info(f"成功解析路由文件: {route_file_path}")
    except SyntaxError as e:
        logging.error(f"解析路由文件失败: {e}")
        conn.close()
        return jsonify({"status": "error", "message": f"解析路由文件失败: {e}"}), 500

    # 解析 vulnerabilities/__init__.py 以建立函数到实际模块的映射
    def parse_vulnerabilities_init(vuln_init_file: str) -> Dict[str, str]:
        """
        解析 vulnerabilities/__init__.py 文件，建立函数名到模块路径的映射。
        返回格式: {function_name: 'module_path'}
        """
        function_map = {}
        try:
            with open(vuln_init_file, "r", encoding="utf-8") as f:
                source = f.read()
            tree = ast.parse(source, filename=vuln_init_file)
            for node in ast.walk(tree):
                if isinstance(node, ast.ImportFrom):
                    module = node.module  # e.g., 'sql'
                    if not module:
                        continue
                    for alias in node.names:
                        if alias.name == '*':
                            # 处理 from .sql import *
                            sub_module = module.lstrip('.')
                            sub_module_path = os.path.join(os.path.dirname(vuln_init_file), sub_module.replace('.', os.sep) + ".py")
                            if os.path.isfile(sub_module_path):
                                try:
                                    with open(sub_module_path, "r", encoding="utf-8") as smf:
                                        sm_source = smf.read()
                                    sm_tree = ast.parse(sm_source, filename=sub_module_path)
                                    for sm_node in ast.walk(sm_tree):
                                        if isinstance(sm_node, ast.FunctionDef):
                                            function_map[sm_node.name] = f"{module}.{sm_node.name}"
                                except Exception as e:
                                    logging.error(f"解析子模块 '{sub_module_path}' 时出错: {e}")
                        else:
                            # 处理 from .sql import encrypt
                            function_map[alias.name] = f"{module}.{alias.name}"
        except Exception as e:
            logging.error(f"解析 vulnerabilities/__init__.py 时出错: {e}")
        return function_map

    function_to_module_map = parse_vulnerabilities_init(vuln_init_file)

    logging.info("函数到模块的映射 (function name -> module path):")
    for func, mod in function_to_module_map.items():
        logging.info(f"  {func} -> {mod}")

    # 解析导入映射
    imports_map: Dict[str, str] = {}
    for node in ast.walk(route_tree):
        if isinstance(node, ast.ImportFrom):
            module = node.module
            if module:
                for alias in node.names:
                    name = alias.asname if alias.asname else alias.name
                    imports_map[name] = module

    logging.info("导入映射 (alias -> module path):")
    for alias, module in imports_map.items():
        logging.info(f"  {alias} -> {module}")

    http_methods = {"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}

    class APIInfo:
        def __init__(self, api_name: str, function_name: str, code: str):
            self.api_name = api_name
            self.function_name = function_name
            self.code = code

    api_infos: List[APIInfo] = []

    # 解析 auth_group 动态获取 auth_group 列表
    auth_group = []
    for node in ast.walk(route_tree):
        if isinstance(node, ast.Assign):
            if isinstance(node.targets[0], ast.Name) and node.targets[0].id == 'auth_group':
                if isinstance(node.value, ast.List):
                    for elt in node.value.elts:
                        if isinstance(elt, ast.Constant) and isinstance(elt.value, str):
                            auth_group.append(elt.value)
                        elif isinstance(elt, ast.Str):
                            auth_group.append(elt.s)

    logging.info("从 routes.py 中读取到的 auth_group 路由：")
    for route in auth_group:
        logging.info(f"  {route}")

    # 获取函数代码直接从 AST
    def get_function_code_from_module(module_path: str, function_name: str) -> Optional[str]:
        """
        根据模块路径和函数名从 AST 中提取函数代码
        """
        # 构建模块文件路径
        # 假设模块路径如 'sql.encrypt' 对应 'vulnerabilities/sql.py'
        module_parts = module_path.split('.')
        if len(module_parts) < 1:
            return None
        # 假设所有模块都在 'vulnerabilities' 目录下
        module_file = os.path.join(parent_dir, "vulnerabilities", *module_parts[:-1], f"{module_parts[-1]}.py")
        if not os.path.isfile(module_file):
            logging.error(f"模块文件不存在: {module_file}")
            return None

        with open(module_file, "r", encoding="utf-8") as f:
            module_source = f.read()

        try:
            module_tree = ast.parse(module_source, filename=module_file)
        except SyntaxError as e:
            logging.error(f"解析模块文件 {module_file} 时出错: {e}")
            return None

        for node in ast.walk(module_tree):
            if isinstance(node, ast.FunctionDef) and node.name == function_name:
                # 提取函数的起始和结束行
                start_line = node.lineno
                end_line = node.end_lineno if hasattr(node, 'end_lineno') else node.lineno + 10
                # 提取源代码
                source_lines = module_source.splitlines()
                function_code = "\n".join(source_lines[start_line - 1:end_line])
                return function_code
        return None

    # 遍历 AST，查找所有 add_url_rule 调用
    for node in ast.walk(route_tree):
        if isinstance(node, ast.Call):
            if isinstance(node.func, ast.Attribute) and node.func.attr == "add_url_rule":
                args = node.args
                if len(args) < 3:
                    continue

                # 提取 API 路径
                api_path_node = args[0]
                if isinstance(api_path_node, ast.Constant) and isinstance(api_path_node.value, str):
                    api_path = api_path_node.value
                elif isinstance(api_path_node, ast.Str):
                    api_path = api_path_node.s
                else:
                    continue

                # 提取 endpoint
                endpoint_node = args[1]
                if isinstance(endpoint_node, ast.Constant) and isinstance(endpoint_node.value, str):
                    func_name = endpoint_node.value
                else:
                    logging.warning(f"无法提取 endpoint，跳过路径: {api_path}")
                    continue

                # 解析函数全名
                function_full_name = function_to_module_map.get(func_name, func_name)
                logging.info(f"解析出的函数名: {func_name}, 全名: {function_full_name}")

                # 提取模块路径和函数名
                if '.' in function_full_name:
                    module_path, function_name = function_full_name.rsplit('.', 1)
                else:
                    module_path = 'vulnerabilities'  # 默认模块路径
                    function_name = function_full_name

                # 提取函数代码
                code = get_function_code_from_module(module_path, function_name)
                if not code:
                    logging.warning(f"无法从模块 {module_path} 中提取函数 {function_name} 的代码。")
                else:
                    logging.info(f"成功提取函数 {function_name} 的代码。")

                # 将 API 信息加入到列表中
                api_infos.append(APIInfo(api_path, function_full_name, code if code else ""))

    # 处理 auth_group 中的 API
    for route in auth_group:
        func_name = route.lstrip('/')  # 去掉路径前的 '/'
        function_full_name = function_to_module_map.get(func_name, func_name)
        logging.info(f"处理 auth_group 路由: Path='{route}', Function='{function_full_name}'")

        # 提取模块路径和函数名
        if '.' in function_full_name:
            module_path, function_name = function_full_name.rsplit('.', 1)
        else:
            module_path = 'vulnerabilities'  # 默认模块路径
            function_name = function_full_name

        # 提取函数代码
        code = get_function_code_from_module(module_path, function_name)
        if not code:
            logging.warning(f"无法从模块 {module_path} 中提取函数 {function_name} 的代码。")
        else:
            logging.info(f"成功提取函数 {function_name} 的代码。")

        logging.info(f"auth_group 路由解析: Path='{route}', Function='{function_full_name}', Code='{code[:30]}...'" if code else f"auth_group 路由解析: Path='{route}', Function='{function_full_name}', Code=空")
        api_infos.append(APIInfo(route, function_full_name, code if code else ""))

    logging.info(f"总共找到 {len(api_infos)} 个 API 路由。")

    # 插入 API 信息到数据库
    try:
        for api in api_infos:
            cursor.execute("""
            INSERT OR IGNORE INTO api_info (api_name, function_name, code) VALUES (?, ?, ?)
            """, (api.api_name.strip('/'), api.function_name, api.code))
        conn.commit()
        logging.info(f"API 信息已成功扫描并存储到数据库。共插入 {len(api_infos)} 条记录。")
    except sqlite3.Error as e:
        logging.error(f"插入 API 信息时出错: {e}")
        conn.rollback()

    conn.close()
    logging.info("数据库连接已关闭。")

    return jsonify({"status": "success", "message": "数据库初始化完成，并插入了示例数据和 API 信息。"}), 200