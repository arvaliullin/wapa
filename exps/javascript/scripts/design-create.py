#!/usr/bin/env python3

import argparse
import json
import requests

def upload_design(api_endpoint, design_file, js_file, wasm_file):
    with open(design_file, 'r') as f:
        design_data = json.load(f)

    data = {
        'name': design_data.get('name'),
        'lang': design_data.get('lang'),
        'functions': json.dumps(design_data.get('functions')),
    }
    files = {
        'js': (design_data.get('js'), open(js_file, 'rb')) if js_file else None,
        'wasm': (design_data.get('wasm'), open(wasm_file, 'rb')) if wasm_file else None,
    }

    files = {k: v for k, v in files.items() if v is not None}

    response = requests.post(api_endpoint, data=data, files=files)

    if response.status_code == 201:
        print(f"Эксперимент создан успешно! ID: {response.text}")
    else:
        print(f"Ошибка: {response.status_code}, {response.text}")

def main():
    parser = argparse.ArgumentParser(description="CLI для создания эксперимента через API")
    parser.add_argument('--endpoint', required=True, help="URL эндпоинта API для создания эксперимента")
    parser.add_argument('--design', required=True, help="Путь к файлу design.json")
    parser.add_argument('--js', required=True, help="Путь к JS-файлу")
    args = parser.parse_args()
    upload_design(args.endpoint, args.design, args.js, None)

if __name__ == '__main__':
    main()
