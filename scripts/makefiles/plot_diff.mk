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


plot-diff-amd64-mean-d_x2Integrate:
	python scripts/plot_arch.py --metric mean --bench d_x2Integrate

plot-diff-amd64-median-d_x2Integrate:
	python scripts/plot_arch.py --metric median --bench d_x2Integrate

plot-diff-amd64-min-d_x2Integrate:
	python scripts/plot_arch.py --metric min --bench d_x2Integrate

plot-diff-amd64-max-d_x2Integrate:
	python scripts/plot_arch.py --metric max --bench d_x2Integrate

plot-diff-amd64-stddev-d_x2Integrate:
	python scripts/plot_arch.py --metric stddev --bench d_x2Integrate

plot-diff-amd64-mean-d_multiply:
	python scripts/plot_arch.py --metric mean --bench d_multiply

plot-diff-amd64-median-d_factorize:
	python scripts/plot_arch.py --metric median --bench d_factorize

plot-diff-amd64-min-d_factorize:
	python scripts/plot_arch.py --metric min --bench d_factorize

plot-diff-amd64-max-d_multiply:
	python scripts/plot_arch.py --metric max --bench d_multiply

plot-diff-amd64-stddev-d_factorize:
	python scripts/plot_arch.py --metric stddev --bench d_factorize

plot-diff-arm64-mean-d_x2Integrate:
	python scripts/plot_arch.py --metric mean --bench d_x2Integrate

plot-diff-arm64-median-d_x2Integrate:
	python scripts/plot_arch.py --metric median --bench d_x2Integrate

plot-diff-arm64-min-d_x2Integrate:
	python scripts/plot_arch.py --metric min --bench d_x2Integrate

plot-diff-arm64-max-d_x2Integrate:
	python scripts/plot_arch.py --metric max --bench d_x2Integrate

plot-diff-arm64-stddev-d_x2Integrate:
	python scripts/plot_arch.py --metric stddev --bench d_x2Integrate

plot-diff-arm64-mean-d_multiply:
	python scripts/plot_arch.py --metric mean --bench d_multiply

plot-diff-arm64-median-d_factorize:
	python scripts/plot_arch.py --metric median --bench d_factorize

plot-diff-arm64-min-d_factorize:
	python scripts/plot_arch.py --metric min --bench d_factorize

plot-diff-arm64-max-d_multiply:
	python scripts/plot_arch.py --metric max --bench d_multiply

plot-diff-arm64-stddev-d_factorize:
	python scripts/plot_arch.py --metric stddev --bench d_factorize
