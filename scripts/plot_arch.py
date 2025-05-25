import argparse
import requests
import sys
import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd

API_DEFAULT = "http://localhost:8080"
ENDPOINT = "/api/benchmark-diff/all"

def fetch_all_diffs(api):
    url = api.rstrip("/") + ENDPOINT
    try:
        resp = requests.get(url)
    except Exception as e:
        print(f"Ошибка при выполнении запроса: {e}")
        sys.exit(1)
    if not resp.ok:
        print(f"Ошибка API {url}: {resp.status_code} {resp.text}")
        sys.exit(1)
    try:
        data = resp.json()
    except Exception as e:
        print(f"Ошибка декодирования JSON: {e}")
        sys.exit(1)
    return data

def flatten_benchmark_diff(data, metric_filter=None, bench_filter=None):
    """Плоская таблица: arch, metric, bench, go, cpp, rust, javascript"""
    rows = []
    for group in data:
        arch = group.get("arch")
        metric = group.get("metric")
        if metric_filter and metric != metric_filter:
            continue
        for res in group.get("results", []):
            name = res.get("name")
            if bench_filter and name != bench_filter:
                continue
            rows.append({
                "arch": arch,
                "metric": metric,
                "name": name,
                "go": res.get("go"),
                "cpp": res.get("cpp"),
                "rust": res.get("rust"),
                "javascript": res.get("javascript")
            })
    return pd.DataFrame(rows)

def plot_bar_arch(df, bench, metric):
    """Plot сравнение архитектур по языкам"""
    dfm = df.melt(
        id_vars=["arch", "metric", "name"],
        value_vars=["go", "cpp", "rust", "javascript"],
        var_name="lang", value_name="time"
    )

    plt.figure(figsize=(8,6))
    sns.barplot(
        data=dfm, x="arch", y="time", hue="lang"
    )
    plt.title(f"Архитектуры по языкам для {bench} ({metric})")
    plt.ylabel("Время, с")
    plt.xlabel("Архитектура")
    plt.tight_layout()
    plt.show()

def main():
    parser = argparse.ArgumentParser(
        description="Построение графика сравнения архитектур для одной функции и метрики (benchmark-diff/all)"
    )
    parser.add_argument("--api", default=API_DEFAULT)
    parser.add_argument("--metric", required=True, help="Метрика (mean, median, min, max, stddev)")
    parser.add_argument("--bench", required=True, help="Имя функции (например d_x2Integrate)")
    args = parser.parse_args()

    data = fetch_all_diffs(args.api)
    df = flatten_benchmark_diff(data, metric_filter=args.metric, bench_filter=args.bench)
    if df.empty:
        print(f"Нет данных для bench={args.bench} и metric={args.metric} на /api/benchmark-diff/all")
        sys.exit(2)
    plot_bar_arch(df, bench=args.bench, metric=args.metric)

if __name__ == "__main__":
    main()
