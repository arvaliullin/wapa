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

def pick_reliable(df, cv_threshold=0.2, min_mean=1e-12):
    """
    На входе DF с колонками: arch, name, lang, metric, value
    Возвращает: только те (arch, name, lang), где stddev/mean < threshold
    """
    means = df[df['metric'] == 'mean'].set_index(['arch','name','lang'])
    stds  = df[df['metric'] == 'stddev'].set_index(['arch','name','lang'])

    joined = means[['value']].join(stds[['value']], how='inner', lsuffix='_mean', rsuffix='_std')
    joined = joined.reset_index()
    joined = joined[joined['value_mean'] > min_mean]
    joined['cv'] = joined['value_std'] / joined['value_mean']

    reliable = joined[joined['cv'] <= cv_threshold]
    return reliable

def main():
    parser = argparse.ArgumentParser(
        description="Отбор надёжных бенчмарков по коэффициенту вариации."
    )
    parser.add_argument("--api", default=API_DEFAULT)
    parser.add_argument("--cv-threshold", type=float, default=0.2, help="Максимальное значение Cv (stddev/mean), по умолчанию 0.2")
    parser.add_argument("--output", type=str, help="Путь к файлу для сохранения результата (csv/json)")
    args = parser.parse_args()

    data = fetch_all(args.api)
    df = build_tables(data)
    reliable = pick_reliable(df, args.cv_threshold)

    print(f"Найдено {len(reliable)} надёжных точек")
    if args.output:
        if args.output.endswith('.csv'):
            reliable.to_csv(args.output, index=False)
        elif args.output.endswith('.json'):
            reliable.to_json(args.output, orient="records", indent=2, force_ascii=False)
        else:
            print("Неизвестный формат вывода!")
    else:
        print(reliable.to_string(index=False))

if __name__ == "__main__":
    main()
