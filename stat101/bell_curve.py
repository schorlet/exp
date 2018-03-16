#-*- coding: utf-8 -*-
import itertools as it
from collections import defaultdict
from math import fsum

if __name__ == '__main__':
	import sys
	n = int(sys.argv[1])

	if 8 < n:
		raise ValueError('unable to compute more than 8')

	p = {
		1: 1/6.0,
		2: 1/6.0,
		3: 1/6.0,
		4: 1/6.0,
		5: 1/6.0,
		6: 1/6.0
	}
	assert fsum(p.values()) == 1

	bell = defaultdict(int)
	for event in it.product(p.keys(), repeat=n):
		se = 0
		ze = 1
		for e in event:
			se += e
			ze *= p[e]
		bell[se] += ze

	# assert fsum(bell.values()) == 1
	assert abs(1-fsum(bell.values())) < 1e-09

	r = 0.6 / max(bell.values())
	for k, v in sorted(bell.items()):
		print '%2d' % k, '%5.2f' % (v*100), '-'*int(round(v*100*r))
