.PHONY: plot-diff install-plot-deps

install-plot-deps:
	pip install matplotlib seaborn pandas requests

plot-diff-amd64-mean:
	python scripts/plot_diff.py --arch amd64 --metric mean --compare lang --plot-type bar

plot-diff-arm64-mean:
	python scripts/plot_diff.py --arch arm64 --metric mean --compare lang --plot-type bar

plot-diff-amd64-median:
	python scripts/plot_diff.py --arch amd64 --metric median --compare lang --plot-type bar

plot-diff-arm64-median:
	python scripts/plot_diff.py --arch arm64 --metric median --compare lang --plot-type bar

plot-diff-amd64-min:
	python scripts/plot_diff.py --arch amd64 --metric min --compare lang --plot-type bar

plot-diff-arm64-min:
	python scripts/plot_diff.py --arch arm64 --metric min --compare lang --plot-type bar

plot-diff-amd64-max:
	python scripts/plot_diff.py --arch amd64 --metric max --compare lang --plot-type bar

plot-diff-arm64-max:
	python scripts/plot_diff.py --arch arm64 --metric max --compare lang --plot-type bar

plot-diff-amd64-stddev:
	python scripts/plot_diff.py --arch amd64 --metric stddev --compare lang --plot-type bar

plot-diff-arm64-stddev:
	python scripts/plot_diff.py --arch arm64 --metric stddev --compare lang --plot-type bar

plot-diff-box-amd64-mean:
	python scripts/plot_diff.py --arch amd64 --metric mean --compare func --plot-type box

plot-diff-box-arm64-mean:
	python scripts/plot_diff.py --arch arm64 --metric mean --compare func --plot-type box

choice:
	python scripts/choice_data.py --arch arm64

