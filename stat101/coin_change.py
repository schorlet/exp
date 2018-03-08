#-*- coding: utf-8 -*-
from collections import Hashable

def memoize(f):
	cache = {}
	def wrapper(total, coins):
		if not isinstance(coins, Hashable):
			return f(total, coins)
		k = (total, coins)
		if k not in cache:
			cache[k] = f(total, coins)
		return cache[k]
	return wrapper

@memoize
def change(total, coins):
	"""
	>>> change(2, [1, 2])
	[{2: 1}, {1: 2}]

	>>> change(3, [1, 2])
	[{1: 1, 2: 1}, {1: 3}]

	>>> change(4, [1, 2])
	[{2: 2}, {1: 2, 2: 1}, {1: 4}]

	>>> change(6, [1, 2])
	[{2: 3}, {1: 2, 2: 2}, {1: 4, 2: 1}, {1: 6}]

	>>> change(5, [1, 2, 3, 4, 5, 6])
	[{5: 1}, {1: 1, 4: 1}, {2: 1, 3: 1}, {1: 2, 3: 1}, {1: 1, 2: 2}, {1: 3, 2: 1}, {1: 5}]

	>>> change(6, [1, 2, 3, 4, 5, 6])
	[{6: 1}, {1: 1, 5: 1}, {2: 1, 4: 1}, {1: 2, 4: 1}, {3: 2}, {1: 1, 2: 1, 3: 1}, {1: 3, 3: 1}, {2: 3}, {1: 2, 2: 2}, {1: 4, 2: 1}, {1: 6}]

	>>> change(7, [1, 2, 3, 4, 5, 6])
	[{1: 1, 6: 1}, {2: 1, 5: 1}, {1: 2, 5: 1}, {3: 1, 4: 1}, {1: 1, 2: 1, 4: 1}, {1: 3, 4: 1}, {1: 1, 3: 2}, {2: 2, 3: 1}, {1: 2, 2: 1, 3: 1}, {1: 4, 3: 1}, {1: 1, 2: 3}, {1: 3, 2: 2}, {1: 5, 2: 1}, {1: 7}]
	"""
	solutions = []

	if 0 == total:
		return solutions

	if 100 < total:
		raise ValueError('unable to compute more than 100')

	coins = tuple(sorted(coins, reverse=True))

	for coin in coins:
		if coin > total:
			continue

		if 1 == coin:
			solutions.append({coin: total})
			break

		nb_coins = total // coin
		for count in range(nb_coins, 0, -1):
			subt = total - (coin * count)

			if 0 == subt:
				solutions.append({coin: count})
				continue

			i = coins.index(coin)
			for suite in change(subt, coins[i+1:]):
				next = {k:v for k,v in suite.items()}
				next[coin] = count
				solutions.append(next)

	return solutions

if __name__ == '__main__':
	import sys
	total = int(sys.argv[1])
	# coins = [1, 2]
	coins = [1, 2, 3, 4, 5, 6]
	for chg in change(total, coins):
		val = sum(k*v for k, v in chg.items())
		print sum(chg.values()), val, chg
		assert val == total
