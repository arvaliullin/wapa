import argparse
import requests
import sys
import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd

API_DEFAULT = "http://localhost:8080"

def fetch_data(api, endpoint, params=None):
    url = api.rstrip("/") + endpoint
    resp = requests.get(url, params=params)
    if not resp.ok:
        print(f"Ошибка API {url}: {resp.status_code} {resp.text}")
        sys.exit(1)
    try:
        data = resp.json()
    except Exception as e:
        print(f"Ошибка разбора JSON: {e}")
        sys.exit(1)
    return data

def flatten_results(results):
    """Преобразует json в DataFrame. results может быть dict или list dicts."""
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

def plot_bar(df, x, y, hue=None, title="", ylabel="", xlabel="", show=True):
    """Barplot"""
    plt.figure(figsize=(10,6))
    sns.barplot(data=df, x=x, y=y, hue=hue)
    plt.title(title)
    plt.ylabel(ylabel)
    plt.xlabel(xlabel)
    plt.xticks(rotation=30)
    plt.tight_layout()
    if show:
        plt.show()

def plot_box(df, x, y, hue=None, title="", ylabel="", xlabel="", show=True):
    """Boxplot (пригодится если несколько запусков для каждой функции)"""
    plt.figure(figsize=(10,6))
    sns.boxplot(data=df, x=x, y=y, hue=hue)
    plt.title(title)
    plt.ylabel(ylabel)
    plt.xlabel(xlabel)
    plt.xticks(rotation=30)
    plt.tight_layout()
    if show:
        plt.show()

def plot_benchmark_diff(api, arch, metric, compare="lang", plot_type="bar"):
    """
    Строит графики для разницы: d_func = func - funcMock
    compare: 'lang' — сравнивать языки для каждой функции
    """
    endpoint = "/api/benchmark-diff"
    data = fetch_data(api, endpoint, params={"metric": metric, "arch": arch})
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
                title=f"Распределение дельты времени по языкам ({metric}, {arch})",
                ylabel="Δ Время, с", xlabel="Язык"
            )
        else:
            plot_bar(
                dfm, x="name", y="time", hue="lang",
                title=f"Сравнение дельты по языкам для функций ({metric}, {arch})",
                ylabel="Δ Время, с", xlabel="Функция"
            )
    elif compare == "func":
        dfm = df.melt(
            id_vars=["arch", "metric", "name"],
            value_vars=["go", "cpp", "rust", "javascript"],
            var_name="lang", value_name="time"
        )
        if plot_type == "box":
            plot_box(
                dfm, x="name", y="time", hue="lang",
                title=f"Boxplot дельты по функциям для языков ({metric}, {arch})",
                ylabel="Δ Время, с", xlabel="Функция"
            )
        else:
            plot_bar(
                dfm, x="lang", y="time", hue="name",
                title=f"Δ времени: функции в каждом языке ({metric}, {arch})",
                ylabel="Δ Время, с", xlabel="Язык"
            )
    else:
        print("compare поддерживает только lang, func")
        sys.exit(2)

def main():
    parser = argparse.ArgumentParser(
        description="Построение графиков сравнения разницы между функцией и её Mock-аналогом (/api/benchmark-diff)"
    )
    parser.add_argument("--api", default=API_DEFAULT)
    parser.add_argument("--arch", type=str, required=True, help="Архитектура (amd64, arm64)")
    parser.add_argument("--metric", type=str, required=True, help="Метрика (mean, median, min, max, stddev)")
    parser.add_argument("--compare", type=str, default="lang",
        choices=["lang", "func"], help="'lang': сравнить языки для каждой функции; 'func': сравнить функции для языка")
    parser.add_argument("--plot-type", type=str, default="bar",
        choices=["bar", "box"], help="Тип графика: bar/box")

    args = parser.parse_args()
    plot_benchmark_diff(
        api=args.api,
        arch=args.arch,
        metric=args.metric,
        compare=args.compare,
        plot_type=args.plot_type
    )

if __name__ == "__main__":
    main()
