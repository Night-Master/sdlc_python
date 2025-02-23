from flask import Flask, request, jsonify
from flask_cors import CORS
from middleware import auth_middleware
from vulnerabilities import (
    sql_injection_sqlite3, sql_injection_sqlite3_safe,get_public_key
    # , reflect_xss, reflect_xss_safe,
    # get_profile_unauthorized, get_products, purchase_product, purchase_product_safe,
    # get_comments, create_comments, create_comments_safe, clear_comments, execute_command,
    # execute_command_safe, get_profile, get_profile_safe, upload_file, upload_file_safe,
    # list_images, download_file, download_file_safe, xml_parse
)

def setup_routes(app: Flask):
    CORS(app, resources={r"/*": {"origins": "*", "headers": ["Origin", "Content-Length", "Content-Type", "Authorization"]}})

    # app.add_url_rule('/init_sqllite3', 'init_sqllite3', init_sqllite3, methods=['GET'])
    app.add_url_rule('/sql_injection_sqlite3_safe', 'sql_injection_sqlite3_safe', sql_injection_sqlite3_safe, methods=['POST'])
    app.add_url_rule('/sql_injection_sqlite3', 'sql_injection_sqlite3', sql_injection_sqlite3, methods=['POST'])
    app.add_url_rule('/getPublicKey', 'getPublicKey', get_public_key, methods=['GET'])
    # app.add_url_rule('/reflect_xss', 'reflect_xss', reflect_xss, methods=['POST'])
    # app.add_url_rule('/reflect_xss_safe', 'reflect_xss_safe', reflect_xss_safe, methods=['POST'])
    # app.add_url_rule('/get_profile_unauthorized', 'get_profile_unauthorized', get_profile_unauthorized, methods=['POST'])

    # 保护需要验证的API

    # auth_group = ['/products', '/purchase', '/purchase_safe', '/get_comments', '/create_comments', '/create_comments_safe', '/clear_comments', '/execute_command', '/execute_command_safe', '/get_profile', '/get_profile_safe', '/upload_file', '/upload_file_safe', '/list_images', '/download', '/download_safe', '/xxe']
    # for route in auth_group:
    #     app.add_url_rule(route, route[1:], auth_middleware(globals()[route[1:]]), methods=['GET', 'POST'])