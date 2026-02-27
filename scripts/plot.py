import argparse
import requests
import sys
import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd

API_DEFAULT = "http://localhost:8080"

BENCHMARK_ENDPOINTS = {
    "notmock": "/api/benchmark/not-mock",
    "mock": "/api/benchmark/mock",
}

def fetch_data(api, endpoint, params=None):
    url = api.rstrip("/") + endpoint
    resp = requests.get(url, params=params)
    if not resp.ok:
        print(f"Ошибка API {url}: {resp.status_code} {resp.text}")
        sys.exit(1)
    try:
        data = resp.json()
    except Exception as e:
        print(f"Невозможно разобрать JSON: {e}")
        sys.exit(1)
    return data

def flatten_results(results):
    """Перевести JSON данные в pandas DataFrame"""
    if isinstance(results, dict):
        results = [results]
    rows = []
    for group in results:
        arch = group.get("arch")
        metric = group.get("metric")
        for r in group["results"]:
            rows.append({
                "arch": arch,
                "metric": metric,
                "name": r["name"],
                "go": r.get("go"),
                "cpp": r.get("cpp"),
                "rust": r.get("rust"),
                "javascript": r.get("javascript")
            })
    return pd.DataFrame(rows)

def plot_bar(df, x, y, hue=None, title="", ylabel="", xlabel=""):
    plt.figure(figsize=(10,6))
    sns.barplot(data=df, x=x, y=y, hue=hue)
    plt.title(title)
    plt.ylabel(ylabel)
    plt.xlabel(xlabel)
    plt.xticks(rotation=30)
    plt.tight_layout()
    plt.show()

def plot_box(df, x, y, hue=None, title="", ylabel="", xlabel=""):
    plt.figure(figsize=(10,6))
    sns.boxplot(data=df, x=x, y=y, hue=hue)
    plt.title(title)
    plt.ylabel(ylabel)
    plt.xlabel(xlabel)
    plt.xticks(rotation=30)
    plt.tight_layout()
    plt.show()

def plot_benchmarks(
        api, 
        arch, 
        metric,
        compare="lang",
        endpoint="notmock",
        plot_type="bar"
    ):
    """
    compare: 'lang' - сравнить языки для каждой функции (default)
    """
    if endpoint not in BENCHMARK_ENDPOINTS:
        print("Неизвестный endpoint!")
        sys.exit(2)

    data = fetch_data(api, BENCHMARK_ENDPOINTS[endpoint], params={"metric": metric, "arch": arch})

    df = flatten_results(data)

    if compare == "lang":
        dfm = df.melt(
            id_vars=["arch", "metric", "name"], 
            value_vars=["go", "cpp", "rust", "javascript"],
            var_name="lang", value_name="time"
        )
        if plot_type == "box":
            plot_box(
                dfm, x="lang", y="time", hue=None,
                title=f"Распределение времени по языкам ({metric}, {arch})",
                ylabel="Время, с", xlabel="Язык"
            )
        else:
            plot_bar(
                dfm, x="name", y="time", hue="lang",
                title=f"Сравнение языков по функциям ({metric}, {arch})",
                ylabel="Время, с", xlabel="Функция"
            )
    else:
        print("Поддерживается только compare=lang для not-mock.")
        sys.exit(3)

def main():
    parser = argparse.ArgumentParser(
        description="Построение графиков сравнения языков по real-функциям (api /api/benchmark/not-mock)"
    )
    parser.add_argument("--api", default=API_DEFAULT)
    parser.add_argument("--arch", type=str, required=True, help="Архитектура (например amd64/arm64)")
    parser.add_argument("--metric", type=str, required=True, help="Метрика (mean/median/min/max/stddev)")
    parser.add_argument("--compare", type=str, default="lang", 
        choices=["lang"], 
        help="Вариант сравнения (только 'lang' поддерживается для not-mock)")
    parser.add_argument("--endpoint", type=str, default="notmock",
        choices=["notmock", "mock"], 
        help="API endpoint (по умолчанию only real-funcs)")
    parser.add_argument("--plot-type", type=str, default="bar",
        choices=["bar", "box"], help="Тип графика: bar или box")

    args = parser.parse_args()
    plot_benchmarks(
        api=args.api,
        arch=args.arch,
        metric=args.metric,
        compare=args.compare,
        endpoint=args.endpoint,
        plot_type=args.plot_type
    )

if __name__ == "__main__":
    main()
