import argparse
import requests
import sys
import pandas as pd

API_DEFAULT = "http://localhost:8080"
ALL_BENCHMARKS_ENDPOINT = "/api/benchmark/all"
LANGS = ["go", "cpp", "rust", "javascript"]

def fetch_all(api):
    url = api.rstrip("/") + ALL_BENCHMARKS_ENDPOINT
    resp = requests.get(url)
    if not resp.ok:
        print(f"Ошибка API {url}: {resp.status_code} {resp.text}")
        sys.exit(1)
    try:
        data = resp.json()
    except Exception as e:
        print(f"Невозможно разобрать JSON: {e}")
        sys.exit(1)
    return data

def build_tables(data):
    rows = []
    for entry in data:
        arch = entry.get("arch")
        metric = entry.get("metric")
        for r in entry["results"]:
            name = r["name"]
            for lang in LANGS:
                if lang in r and r[lang] is not None:
                    rows.append({
                        "arch": arch,
                        "name": name,
                        "lang": lang,
                        "metric": metric,
                        "value": r[lang]
                    })
    return pd.DataFrame(rows)

def pick_reliable(df, arch, cv_threshold=0.2, min_mean=1e-12):
    """
    df: DataFrame с колонками arch, name, lang, metric, value
    arch: выбранная архитектура (строка)
    возвращает: список функций, которые можно доверять сравнивать по всем языкам
    """
    df = df[df['arch'] == arch]
    means = df[df['metric'] == 'mean'].set_index(['name','lang'])
    stds  = df[df['metric'] == 'stddev'].set_index(['name','lang'])

    joined = means[['value']].join(stds[['value']], how='inner', lsuffix='_mean', rsuffix='_std')
    joined = joined.reset_index()
    joined = joined[joined['value_mean'] > min_mean]
    joined['cv'] = joined['value_std'] / joined['value_mean']

    reliable = joined[joined['cv'] <= cv_threshold]

    grouped = reliable.groupby('name').lang.nunique()
    reliable_names = grouped[grouped == len(LANGS)].index.tolist()
    return sorted(reliable_names)

def main():
    parser = argparse.ArgumentParser(
        description="Отбор надёжных бенчмарков по коэффициенту вариации."
    )
    parser.add_argument("--api", default=API_DEFAULT)
    parser.add_argument("--arch", required=True, help="Архитектура (например, arm64, amd64)")
    parser.add_argument("--cv-threshold", type=float, default=0.15, help="Максимальное значение Cv (stddev/mean), по умолчанию 0.15")
    args = parser.parse_args()

    data = fetch_all(args.api)
    df = build_tables(data)
    reliable_names = pick_reliable(df, args.arch, args.cv_threshold)

    if reliable_names:
        print(f"Для архитектуры {args.arch} можно использовать ({len(reliable_names)}):\n" +
              ", ".join(reliable_names))
    else:
        print(f"Нет надёжных функций для архитектуры {args.arch}")

if __name__ == "__main__":
    main()
