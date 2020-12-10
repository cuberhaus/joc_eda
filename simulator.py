# Exemple per executar: 
#   python simulate.py Stonks Stonks Stonks2 Stonks2 100

import subprocess
import random
import os
import sys

from multiprocessing import Pool
from collections import defaultdict

N_THREADS = 6 # Recomanació: 2*(nombre de nuclis) - 1 o 2*(nombre de nuclis) - 2

def player(x):
    pi1,pi2,pi3,pi4,seed = x
    command = f"./Game -s {seed} {pi1} {pi2} {pi3} {pi4} --input default.cnf --output /dev/null"
    res = subprocess.run(list(command.split(" ")), capture_output=True)
    for line in res.stderr.decode().splitlines():
        if "got top score" in line:
            for p in (pi1,pi2,pi3,pi4):
                if p in line.split():
                  return (p, seed)

if __name__ == '__main__':
    assert(len(sys.argv) >= 6)
    
    p1,p2,p3,p4, N_PARTIDES = sys.argv[1:]
    
    N_PARTIDES = int(N_PARTIDES)
    
    with Pool(processes=N_THREADS) as pool:
      result = pool.map(player, ((p1,p2,p3,p4,i) for i in random.choices(range(99999), k=N_PARTIDES)))
      res = defaultdict(list) 
      for i, j in result: 
          res[i].append(j)
      for p in res:
        print(f"Player {p} won with seeds {res[p]}\n")
        
      print("\n \nAND:")
      
      for p in res:
        print(f"Player {p:12} got {len(res[p])} wins, ie {100.0*len(res[p])/N_PARTIDES:2.2f}%")
      
      
