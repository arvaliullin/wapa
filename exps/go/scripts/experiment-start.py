#!/usr/bin/env python3

import argparse
import json
import requests
import sys

def get_design_uid(api_endpoint, name):
    """Получить UID дизайна по его имени через API."""
    try:
        response = requests.get(f"{api_endpoint}/designs")
        if response.status_code != 200:
            print(f"Ошибка получения списка дизайнов: {response.status_code}, {response.text}")
            sys.exit(1)

        designs = response.json()
        for design in designs:
            if design.get('name') == name:
                return design.get('id')
        
        print(f"Дизайн с именем '{name}' не найден.")
        sys.exit(1)

    except requests.RequestException as e:
        print(f"Ошибка при выполнении запроса: {e}")
        sys.exit(1)

def start_experiment(api_endpoint, uid, repeats, warmup):
    """Запуск эксперимента с указанным UID и параметрами."""
    url = f"{api_endpoint}/experiment/{uid}/start"
    payload = {
        "repeats": repeats,
        "warmup": warmup
    }

    try:
        response = requests.post(url, json=payload)
        
        if response.status_code == 200:
            print(f"Эксперимент успешно запущен: {response.json()}")
        elif response.status_code == 404:
            print("Дизайн не найден.")
        elif response.status_code == 400:
            print(f"Некорректный запрос: {response.text}")
        else:
            print(f"Ошибка сервера: {response.status_code}, {response.text}")
    except requests.RequestException as e:
        print(f"Ошибка при отправке запроса на запуск эксперимента: {e}")
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description="CLI для запуска эксперимента через API")
    parser.add_argument('--endpoint', required=True, help="URL эндпоинта API")
    parser.add_argument('--name', required=True, help="Имя дизайна для получения UID")
    parser.add_argument('--repeats', type=int, default=1, help="Количество повторений для эксперимента")
    parser.add_argument('--warmup', action='store_true', help="Флаг разогрева (warmup) перед запуском эксперимента")
    args = parser.parse_args()

    api_endpoint = args.endpoint.rstrip('/')
    design_name = args.name
    repeats = args.repeats
    warmup = args.warmup

    print(f"Получение UID для дизайна с именем '{design_name}'...")
    uid = get_design_uid(api_endpoint, design_name)

    if uid:
        print(f"UID найден: {uid}. Запуск эксперимента...")
        start_experiment(api_endpoint, uid, repeats, warmup)

if __name__ == '__main__':
    main()
