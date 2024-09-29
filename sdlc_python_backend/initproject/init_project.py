import ast
import os
import sqlite3
from typing import Dict, List, Optional
from flask import jsonify
import sys
def get_project_name():
    return jsonify({"project_name":"Python"}),200