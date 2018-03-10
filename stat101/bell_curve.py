#-*- coding: utf-8 -*-
import itertools as it
from collections import defaultdict

if __name__ == '__main__':
	import sys
	n = int(sys.argv[1])

	# coins = (1, 2)
	coins = (1, 2, 3, 4, 5, 6)

	bell = defaultdict(int)
	for event in it.product(coins, repeat=n):
		bell[sum(event)]+=1

	for k, v in sorted(bell.items()):
		print '%2d' % k, '-'*v

