import argparse
import requests
import sys
import json
import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd

API_DEFAULT = "http://localhost:8080"

BENCHMARK_ENDPOINTS = {
    "all": "/api/benchmark/all",
    "diff": "/api/benchmark-diff/all",
    "bench": "/api/benchmark",
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
    rows = []
    if isinstance(results, dict):
        results = [results]
    for group in results:
        arch = group.get("arch")
        metric = group.get("metric")
        for r in group["results"]:
            rows.append({
                "arch": arch,
                "metric": metric,
                "name": r["name"],
                "go": r["go"],
                "cpp": r["cpp"],
                "rust": r["rust"],
                "javascript": r["javascript"]
            })
    return pd.DataFrame(rows)

def plot_bar(df, x, y, hue=None, title="", ylabel="", xlabel=""):
    plt.figure(figsize=(10,6))
    sns.barplot(data=df, x=x, y=y, hue=hue)
    plt.title(title)
    plt.ylabel(ylabel)
    plt.xlabel(xlabel)
    plt.tight_layout()
    plt.show()

def plot_box(df, x, y, hue=None, title="", ylabel="", xlabel=""):
    plt.figure(figsize=(12,6))
    sns.boxplot(data=df, x=x, y=y, hue=hue)
    plt.title(title)
    plt.ylabel(ylabel)
    plt.xlabel(xlabel)
    plt.tight_layout()
    plt.show()

def plot_benchmarks(
    api, 
    arch=None, 
    metric=None, 
    function=None,
    compare="lang", 
    endpoint="all"
):
    """
    compare: могут быть 'lang', 'function', 'arch', 'metric', 'mock'
    """
    if endpoint == "all":
        data = fetch_data(api, BENCHMARK_ENDPOINTS["all"])
    elif endpoint == "diff":
        data = fetch_data(api, BENCHMARK_ENDPOINTS["diff"])
    elif endpoint == "bench":
        data = fetch_data(api, BENCHMARK_ENDPOINTS["bench"], 
            params={"metric": metric, "arch": arch})
    elif endpoint == "notmock":
        data = fetch_data(api, BENCHMARK_ENDPOINTS["notmock"], 
            params={"metric": metric, "arch": arch})
    elif endpoint == "mock":
        data = fetch_data(api, BENCHMARK_ENDPOINTS["mock"], 
            params={"metric": metric, "arch": arch})
    else:
        print("Неизвестный endpoint!")
        sys.exit(2)

    df = flatten_results(data)

    if arch:
        df = df[df.arch == arch]
    if metric:
        df = df[df.metric == metric]
    if function:
        df = df[df.name.str.contains(function)]
    
    if compare == "lang":
        dfm = df.melt(id_vars=["arch","metric","name"], 
                      value_vars=["go","cpp","rust","javascript"],
                      var_name="lang", value_name="time")
        plot_bar(dfm, x="name", y="time", hue="lang",
                 title=f"Время выполнения задач ({metric}, {arch})",
                 ylabel="Время, с", xlabel="Функция")
    elif compare == "arch":
        dfm = df.melt(id_vars=["arch","metric","name"],
                      value_vars=["go","cpp","rust","javascript"],
                      var_name="lang", value_name="time")
        plot_bar(dfm, x="arch", y="time", hue="lang",
                title=f"Сравнение архитектур ({metric})",
                ylabel="Время, с", xlabel="Архитектура")
    elif compare == "function":
        dfm = df.melt(id_vars=["arch","metric","name"],
                      value_vars=["go","cpp","rust","javascript"],
                      var_name="lang", value_name="time")
        plot_bar(dfm, x="name", y="time", hue="arch",
            title=f"Сравнение функций по архитектурам ({metric})",
            ylabel="Время, с", xlabel="Функция")
    elif compare == "metric":
        dfm = df.melt(id_vars=["arch","metric","name"],
                      value_vars=["go","cpp","rust","javascript"],
                      var_name="lang", value_name="time")
        plot_bar(dfm, x="metric", y="time", hue="lang",
                 title=f"Сравнение метрик по языкам",
                 ylabel="Время, с", xlabel="Метрика")
    elif compare == "mock":
        dfmock = df[df.name.str.endswith("Mock")]
        dfnm = df[~df.name.str.endswith("Mock")]
        for lang in ["go", "cpp", "rust", "javascript"]:
            plt.figure(figsize=(10,6))
            plt.bar(dfnm.name, dfnm[lang], label='real')
            plt.bar(dfmock.name, dfmock[lang], label='mock', alpha=0.7)
            plt.title(f"Mock vs Not-Mock {lang} ({metric}, {arch})")
            plt.xlabel("Функция")
            plt.ylabel("Время, с")
            plt.legend()
            plt.tight_layout()
            plt.show()
    else:
        print("Неизвестный режим сравнения.")
        sys.exit(3)

def main():
    parser = argparse.ArgumentParser(
        description="Графики для анализа бенчмарков (API wapa)"
    )
    parser.add_argument("--api", default=API_DEFAULT)
    parser.add_argument("--arch", type=str, help="Архитектура (amd64/arm64)")
    parser.add_argument("--metric", type=str, help="Метрика (mean/median/min/max/stddev)")
    parser.add_argument("--function", type=str, help="Функция, фильтр по имени")
    parser.add_argument("--compare", type=str, default="lang",
        choices=["lang", "arch", "function", "metric", "mock"], 
        help="Тип сравнения")
    parser.add_argument("--endpoint", type=str, default="all",
        choices=["all", "diff", "bench", "notmock", "mock"], 
        help="API endpoint")
    
    args = parser.parse_args()
    plot_benchmarks(
        api=args.api,
        arch=args.arch,
        metric=args.metric,
        function=args.function,
        compare=args.compare,
        endpoint=args.endpoint
    )

if __name__ == "__main__":
    main()
