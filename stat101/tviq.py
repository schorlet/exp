#!/usr/bin/env python3
#-*- coding: utf-8 -*-
import csv
from pprint import pprint
import math

from matplotlib import pyplot


def mean(data):
	return float(sum(data)) / len(data)


def variance(data, mu=None):
	if mu is None:
		mu = mean(data)
	return sum(pow(x-mu, 2) for x in data) / len(data)


def stdev(data, mu=None):
	return math.sqrt(variance(data, mu))


def spearman(data):
	# rank_x
	data = [tuple(row)+(i+1,)
		for i, row in enumerate(sorted(data, key=lambda r:r[0]))]
	# rank_y
	data = [tuple(row)+(i+1,)
		for i, row in enumerate(sorted(data, key=lambda r:r[1]))]
	# pprint(sorted(data))

	rank = range(1, len(data) + 1)
	mu = (1 + rank[-1]) / 2.0
	dev = stdev(rank, mu)
	# print rank, mu, dev

	cov = sum((row[2]-mu) * (row[3]-mu) for row in data) / len(data)
	# print(cov)
	p = cov / (dev * dev)
	print('spearman: {}'.format(p))


def pearson(data):
	xs = [row[0] for row in data]
	mean_x = mean(xs)
	dev_x = stdev(xs, mean_x)
	# print(mean_x, dev_x)

	ys = [row[1] for row in data]
	mean_y = mean(ys)
	dev_y = stdev(ys, mean_y)
	# print(mean_y, dev_y)

	cov = sum((row[0]-mean_x) * (row[1]-mean_y) for row in data) / len(data)
	# print(cov)
	p = cov / (dev_x * dev_y)
	print('pearson: {}'.format(p))


def scatter(path, data):
	xs = [row[0] for row in data]
	ys = [row[1] for row in data]

	pyplot.scatter(xs, ys, c='red', marker='+')
	pyplot.grid()
	# pyplot.savefig(path)
	pyplot.show()


if __name__ == '__main__':
	with open('tviq.txt', newline='') as lines:
		rows = csv.reader(lines, delimiter='\t', strict=True)
		data = [tuple(map(int, row)) for row in rows]

	spearman(data)
	pearson(data)

	scatter('tviq.png', data)
