<template>
  <div>
    <v-card
      @click="showEtablissement()"
      style="height: 80px; background: linear-gradient(#fff, #eee 45%, #ccc);"
      class="elevation-2 ma-2 pointer"
    >
      <div style="height: 100%; width: 100%; overflow: hidden;">
        <div class="entete pointer">
          <PredictionWidgetScore id="widget" :prob="prediction.prob" :diff="prediction.diff"/>
        </div>
        <div class="corps">
          <div style="left: 250px; position: absolute;" :id="'marge_' + prediction._id.siret"></div>
          <div style="white-space: nowrap; overflow: hidden; max-width: 400px; max-height:40px">
            <span style="font-size: 28px; color: #333; line-height: 40px; display: inline-block; font-family: 'Oswald'; max-width: '100px'">
              {{ prediction.etablissement.sirene.raison_sociale }}
            </span>
          </div>
          <span style="font-family: 'Abel'; font-size: 17px; font-weight: 400;">{{ (naf.n5 || {})[(sirene.ape || '')] }}</span>

          <v-img
            style="position: absolute; left: 550px; top: 10px;"
            width="70"
            :src="'/static/' + (urssaf?'red':'gray') + '_urssaf.svg'"
          ></v-img>

          <v-img
            style="position: absolute; left: 550px; bottom: 10px;"
            width="18"
            height="24"
            :src="'/static/' + apart + '_apart.svg'"
          ></v-img>

          <div style="position: absolute; left: 584px; bottom: 3px; color: #333">
            <span
              style="font-size: 24px"
            >{{ prediction.etablissement.dernier_effectif.effectif || 'n/c' }}</span>
            <v-icon style="position:relative; top:-4px" small v-if="variationEffectif">{{ variationEffectif }}</v-icon>
          </div>

          <div class="flex" style="position:absolute; left: 600px; top: 0px; bottom: 0px; right: 9px;">
            <div class='label'>
              Chiffre d'affaire<br/>
              <span style="font-size: 25px" :class="diane.ca_arrow?'':'unknown'">{{ diane.ca }}</span>
            </div>
          </div>
          <div style="position:absolute; left: 780px; top: 42px">
            <v-icon small v-if="diane.ca_arrow">{{ diane.ca_arrow }}</v-icon>
          </div>
          <div class="flex" style="position:absolute; left: 750px; top: 0px; bottom: 0px; right: 9px;">
            <div class='label'>
              Résultat d'exploitation<br/>
              <span style="font-size: 25px" :class='diane.resultat_expl_color'>{{ diane.resultat_expl }}</span>
            </div>
          </div>
        </div>
        <v-dialog  lazy fullscreen v-model="dialog">
          <div style="height: 100%; width: 100%;  font-weight: 800; font-family: 'Oswald', sans;">
            <v-toolbar fixed class="toolbar" height="35px" style="color: #fff; font-size: 22px;">
              <v-spacer/>
                {{ prediction.etablissement.sirene.raison_sociale }}
              <v-spacer/>
              <v-icon @click="dialog=false;log()" style="color: #fff">mdi-close</v-icon>
            </v-toolbar>
            <Etablissement :siret="prediction._id.key"></Etablissement>
          </div>
        </v-dialog>
      </div>
    </v-card>
  </div>
</template>

<script>
import Etablissement from '@/components/Etablissement'
import PredictionWidgetScore from '@/components/PredictionWidgetScore'

export default {
  props: ['prediction'],
  components: {
    PredictionWidgetScore,
    Etablissement
  },
  data () {
    return {
      dialog: false,
      expand: false,
      rowsPerPageItems: [4, 8, 12],
      pagination: {
        rowsPerPage: 4
      }
    }
  },
  mounted () {
    console.log(this.prediction)
  },
  computed: {
    naf () {
      return this.$store.state.naf
    },
    apart () {
      var apart = this.prediction.etablissement.apdemande.filter(d => {
        var end = new Date(d.periode.end)
        return end.getTime() > new Date(new Date().setFullYear(new Date().getFullYear() - 1))
      })
      return (apart.length > 0) ? 'red' : 'gray'
    },
    sirene () {
      return this.prediction.etablissement.sirene
    },
    variationEffectif () {
      var effectifArray = this.prediction.etablissement.effectif
      var l = effectifArray.length
      if (l === 0) { return 'none' }

      var effectif = effectifArray[l - 1].effectif
      var effectifPrecedent = effectifArray[l - Math.min(12, l)].effectif
      if (effectif / effectifPrecedent > 1.10) {
        return 'mdi-arrow-bottom-right'
      }
      if (effectif / effectifPrecedent < 0.90) {
        return 'mdi-arrow-top-right'
      }
      return null
    },
    urssaf () {
      var urssaf = this.prediction.etablissement.debit.map(d => d.part_patronale + d.part_ouvriere)
      var l = urssaf.length
      for (var i = l - 3; i < l; i++) {
        if (urssaf[i] / urssaf[i - 1] > 1.01) {
          return true
        }
      }
      return false
    },
    diane () {
      var entreprise = this.prediction.entreprise || { diane: [] }
      var diane = entreprise.diane.reduce((m, d) => {
        if (d.ca && d.resultat_expl && d.exercice_diane && m.length < 2) {
          m.push({
            exercice: d.exercice_diane,
            ca: d.ca,
            resultat_expl: d.resultat_expl
          })
        }
        return m
      }, [])

      if (diane.length === 2) {
        return {
          ca: diane[0].ca + ' k€',
          resultat_expl: diane[0].resultat_expl + ' k€',
          ca_arrow: (diane[0].ca / diane[1].ca > 1.05) ? 'mdi-arrow-top-right' : (diane[0].ca / diane[1].ca < 0.95) ? 'mdi-arrow-bottom-right' : null,
          resultat_expl_color: diane[0].resultat_expl <= 0 ? 'down' : null
        }
      }
      return {
        ca: 'n/c',
        resultat_expl: 'n/c',
        ca_color: 'unknown',
        resultat_expl_color: 'unknown'
      }
    }
  },
  methods: {
    log () {
      console.log(this.prediction)
    },
    upOrDown (before, after, treshold) {
      if (before == null || after == null) {
        return 'mdi-help-circle'
      }
      if (after / before > 1 + treshold) {
        return 'mdi-arrow-up'
      }
      if (after / before < 1 - treshold) {
        return 'mdi-arrow-down'
      }
      return 'mdi-tilde'
    },
    upOrDownClass (before, after, treshold) {
      if (before == null || after == null) {
        return 'unknown'
      }
      if (after / before > 1 + treshold) {
        return 'high'
      }
      if (after / before < 1 - treshold) {
        return 'down'
      }
      return 'none'
    },
    showEtablissement () {
      this.dialog = true
    }
  }
}
</script>

<style scoped>
div.flex {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  font-size: 16px;
  font-family: "Quicksand"
}
div.label {
  text-align: right;
  font-family: 'Abel';
  width: 180px;
}
div.entete {
  float: left;
  background: linear-gradient(
    270deg,
    rgba(119, 122, 170, 0.219),
    rgba(119, 122, 170, 0)
  );
  border-right: solid 1px #3334;
  width: 80px;
  height: 80px;
  text-align: center;
  padding: 20px;
}
div.corps {
  flex: 1;
  padding: 5px;
  margin-left: 80px;
  height: 80px;
  background: linear-gradient(45deg, rgba(50, 51, 121, 0.212), #0000);
}
.high {
  color: rgb(4, 153, 41);
}
.down {
  color: rgb(244, 67, 54);
}
.unknown {
  color: rgb(150, 150, 150);
}
td {
  width: 80px;
}
.pointer {
  cursor: pointer;
}
</style>
