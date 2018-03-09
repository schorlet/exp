#-*- coding: utf-8 -*-
import coin_change as coin
import itertools as it
from collections import defaultdict

if __name__ == '__main__':
	import sys
	n = int(sys.argv[1])

	# coins = [1, 2]
	coins = [1, 2, 3, 4, 5, 6]
	nb_coins = len(coins)

	bell = defaultdict(int)
	total = 0

	for x in range(n, n*nb_coins+1):
		xcount = 0
		for chg in coin.change(x, coins):
			if sum(chg.values()) != n:
				continue

			vals = []
			for k, v in chg.items():
				for vi in range(v):
					vals.append(k)

			count = len(set(it.permutations(vals, n)))
			xcount += count
			total += count

			print x, vals, count
		bell[x] += xcount
	assert total == pow(nb_coins, n)

	for k, v in sorted(bell.items()):
		print '%2d' % k, '-'*v
