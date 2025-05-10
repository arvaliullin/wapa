.PHONY: plot-not-mock install-plot-deps

install-plot-deps:
	pip install matplotlib seaborn pandas requests

plot-not-mock-amd64-mean:
	python scripts/plot.py --arch amd64 --metric mean --endpoint notmock --compare lang

plot-not-mock-arm64-mean:
	python scripts/plot.py --arch arm64 --metric mean --endpoint notmock --compare lang

plot-not-mock-amd64-median:
	python scripts/plot.py --arch amd64 --metric median --endpoint notmock --compare lang

plot-not-mock-arm64-median:
	python scripts/plot.py --arch arm64 --metric median --endpoint notmock --compare lang

plot-not-mock-box-amd64-mean:
	python scripts/plot.py --arch amd64 --metric mean --endpoint notmock --compare lang --plot-type box

plot-not-mock-box-arm64-mean:
	python scripts/plot.py --arch arm64 --metric mean --endpoint notmock --compare lang --plot-type box
