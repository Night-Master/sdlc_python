from flask import Flask, request, jsonify
from flask_cors import CORS
from middleware import auth_middleware
from initproject import initialize_database,get_project_name
from utils import get_api_code
from vulnerabilities import (
    sql_injection_sqlite3, sql_injection_sqlite3_safe,get_public_key
    , reflect_xss, reflect_xss_safe,list_images, download, download_safe, upload_file, upload_file_safe,get_comments, create_comments, create_comments_safe, clear_comments,execute_command_safe,execute_command
    , products, purchase, purchase_safe
    , get_profile, get_profile_safe,get_profile_unauthorized, change_password_unsafe,change_password_safe
    
    # , xml_parse
)

def setup_routes(app: Flask):
    CORS(app, resources={r"/*": {"origins": "*", "headers": ["Origin", "Content-Length", "Content-Type", "Authorization"]}})

    # app.add_url_rule('/init_sqllite3', 'init_sqllite3', init_sqllite3, methods=['GET'])
    app.add_url_rule('/sql_injection_sqlite3_safe', 'sql_injection_sqlite3_safe', sql_injection_sqlite3_safe, methods=['POST'])
    app.add_url_rule('/sql_injection_sqlite3', 'sql_injection_sqlite3', sql_injection_sqlite3, methods=['POST'])
    app.add_url_rule('/getPublicKey', 'getPublicKey', get_public_key, methods=['GET'])
    app.add_url_rule('/reflect_xss', 'reflect_xss', reflect_xss, methods=['POST'])
    app.add_url_rule('/reflect_xss_safe', 'reflect_xss_safe', reflect_xss_safe, methods=['POST'])
    app.add_url_rule('/get_profile_unauthorized', 'get_profile_unauthorized', get_profile_unauthorized, methods=['POST'])
    app.add_url_rule('/init_sqllite3', 'initialize_database', initialize_database, methods=['GET'])
    app.add_url_rule('/get_api_code', 'get_api_code', get_api_code, methods=['POST'])
    app.add_url_rule('/get_project_name', 'get_project_name', get_project_name, methods=['GET'])
    
    # 保护需要验证的API
    auth_group = ['/change_password_unsafe','/change_password_safe','/products', '/purchase', '/purchase_safe', '/get_profile', '/get_profile_safe', '/list_images', '/download', '/download_safe', '/upload_file', '/upload_file_safe', '/get_comments', '/create_comments', '/create_comments_safe', '/clear_comments','/execute_command_safe','/execute_command']

    # auth_group = [  '/xxe']
    for route in auth_group:
        app.add_url_rule(route, route[1:], auth_middleware(globals()[route[1:]]), methods=['GET', 'POST'])