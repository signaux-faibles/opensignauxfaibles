BEGIN {
       FS = ";"
       OFS = ";"
}
NR==1 {
    printf "%s", "\"Annee\""
    for (i=1;i<=NF;i++) {
        if ((tgt == "") || ($i !~ "201") || ($i ~ tgt)) {
            f[++nf] = i 
	    printf "%s%s", OFS, gensub(" "tgt,"","g",$i)
        }
    }
    printf "%s", ORS
}
NR > 1{
    printf "%i%s", tgt, OFS
    for (i=1; i<=nf; i++) {
        printf "%s%s", $(f[i]), (i<nf?OFS:ORS)
    }
}
