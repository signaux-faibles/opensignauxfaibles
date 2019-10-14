#!/usr/bin/env bash
cat "$@" |
  head |
  awk 'BEGIN{OFS=FS=""; srand();
  new[1] = "a"
  new[2] = "b"
  new[3] = "c"
  new[4] = "d"
  new[5] = "e"
  new[6] = "f"
  new[7] = "g"
  new[8] = "h"
  new[9] = "i"
  new[10] = "j"
  new[11] = "k"
  new[12] = "l"
  new[13] = "m"
  new[14] = "n"
  new[15] = "o"
  new[16] = "p"
  new[17] = "q"
  new[18] = "r"
  new[19] = "s"
  new[20] = "t"
  new[21] = "u"
  new[22] = "v"
  new[23] = "w"
  new[24] = "x"
  new[25] = "y"
  new[26] = "z"
}
  NR>1{for (i=1; i <= NF; i++) {sub(/[0-9]/, int(10*rand()), $i);
    sub(/[[:alpha:]]/, new[int(26*rand() + 1)], $i);}}1'

