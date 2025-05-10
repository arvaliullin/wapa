.PHONY: plot install-plot-deps

# Установка python-библиотек
install-plot-deps:
	pip install matplotlib seaborn pandas requests

# График: сравнение языков для amd64/mean
plot-lang-amd64:
	python scripts/plot.py --arch amd64 --metric mean --compare lang --endpoint all

# Сравнение mock vs не-mock
plot-mock-vs-real:
	python scripts/plot.py --arch amd64 --metric mean --compare mock --endpoint all

# Дифф между функциями (diff endpoint)
plot-diff-functions:
	python scripts/plot.py --metric mean --compare function --endpoint diff

# Сравнение языков между архитектурами
plot-arch:
	python scripts/plot.py --metric mean --compare arch --endpoint all

# Коробочные графики распределения по метрикам
plot-metric-box:
	python scripts/plot.py --compare metric --endpoint all
