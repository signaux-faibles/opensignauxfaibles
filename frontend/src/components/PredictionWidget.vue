<template>
  <div >
    <v-card
      @click="showEtablissement()"
      style="height: 80px; background: linear-gradient(#fff, #eee 45%, #ccc);"
      class="elevation-2 ma-2 pointer"
    >
      <div style="height: 100%; width: 100%; overflow: hidden;" >
        <div class="entete pointer" >
          <PredictionWidgetScore id="widget" :prob="prediction.prob" :diff="prediction.diff"/>
        </div>
        <div class="corps">
          <div 
          style="left: 250px; position: absolute;"
          :id="'marge_' + prediction._id.siret"></div>
          <div style="white-space: nowrap; overflow: hidden; max-width: 400px; max-height:30px">
            <span style="font-size: 18px; color: #333; line-height: 10px; font-family: 'Oswald';">{{ prediction.etablissement.sirene.raisonsociale }}<br style="line-height: 10px;"/></span>
          </div>
          <div style="left: 450px; position: absolute; top: 3px;  padding: 2px">
            <!-- <b>Chiffre d'affaire (k€)</b><br/> -->
            <!-- <table>
              <tr>
                <td>{{ prediction.entreprise.diane[0].exercice_diane || '' }}</td>
                <td>{{ prediction.entreprise.diane[1].exercice_diane || '' }}</td>
              </tr>
              <tr>
                <td>
                  {{ prediction.entreprise.diane[0].ca || '' }} 
                  <v-icon 
                  :class="upOrDownClass(prediction.entreprise.diane[1].ca, prediction.entreprise.diane[0].ca, 0.04)" small
                  >
                    {{ upOrDown(prediction.entreprise.diane[1].ca, prediction.entreprise.diane[0].ca, 0.04) }}
                  </v-icon>
                </td>
                <td>
                  {{ prediction.entreprise.diane[1].ca || '' }} 
                  <v-icon 
                  :class="upOrDownClass(prediction.entreprise.diane[2].ca, prediction.entreprise.diane[1].ca, 0.04)" small
                  >
                    {{ upOrDown(prediction.entreprise.diane[2].ca, prediction.entreprise.diane[1].ca, 0.04) }}</v-icon>
                <td/>
              </tr>
            </table>
          </div>

          <div style="left: 620px; position: absolute; top: 3px;  padding: 2px">
            <b>EBE (k€)</b><br/>
            <table>
              <tr>
                <td>{{ prediction.entreprise.diane[0].exercice_diane || 'n/c' }}</td>
                <td>{{ prediction.entreprise.diane[1].exercice_diane || 'n/c' }}</td>
              </tr>
              <tr>
                <td>{{ prediction.entreprise.diane[0].excedent_brut_d_exploitation || 'n/c' }}</td>
                <td>{{ prediction.entreprise.diane[1].excedent_brut_d_exploitation || 'n/c' }}</td>
              </tr>
            </table>
          </div>

          <div style="left: 790px; position: absolute; top: 3px;  padding: 2px">
            <b>Bénéfice ou perte (k€)</b><br/>
            <table>
              <tr>
                <td>{{ prediction.entreprise.diane[0].exercice_diane || 'n/c' }}</td>
                <td>{{ prediction.entreprise.diane[1].exercice_diane || 'n/c' }}</td>
              </tr>
              <tr>
                <td>{{ prediction.entreprise.diane[0].benefice_ou_perte || 'n/c' }}</td>
                <td>{{ prediction.entreprise.diane[1].benefice_ou_perte || 'n/c' }}</td>
              </tr>
            </table>
          </div>

          <div style="left: 960px; position: absolute; top: 3px;  padding: 2px">
            <b>Liquidité réduite (%)</b><br/>
            <table>
              <tr>
                <td>{{ prediction.entreprise.diane[0].exercice_diane || 'n/c' }}</td>
                <td>{{ prediction.entreprise.diane[1].exercice_diane || 'n/c' }}</td>
              </tr>
              <tr>
                <td>{{ prediction.entreprise.diane[0].liquidite_reduite || 'n/c' }}</td>
                <td>{{ prediction.entreprise.diane[1].liquidite_reduite || 'n/c' }}</td>
              </tr>
            </table> -->
          </div>

          <span style="font-size: 12px; color: #333; line-height: 10px;">{{ prediction._id.siret }}<br style="line-height: 10px;"/></span>
          <v-img style="position: absolute; left: 160px; bottom: 10px;" width="17" src="/static/gray_apart.svg"></v-img>
          <v-img style="position: absolute; left: 90px; bottom: 10px;" width="57" :src="'/static/' + (prediction.etablissement.urssaf?'red':'gray') + '_urssaf.svg'"></v-img>
          <div style="position: absolute; left: 195px; bottom: 4px; color: #333">
            <span :class="variationEffectif" style="font-size: 20px">{{ prediction.etablissement.effectif || 'n/c' }}</span>
          </div>

        </div>
        <v-dialog 
        attach="#detection"
        lazy
        fullscreen
        v-model="dialog">
          <div style="height: 100%; width: 100%;  font-weight: 800; font-family: 'Abel', sans;">
            <v-toolbar fixed class="toolbar" height="35px" style="color: #fff; font-size: 22px;">
              <v-spacer/>
                {{ prediction.etablissement.sirene.raisonsociale }}
              <v-spacer/>
              <v-icon @click="dialog=false"  style="color: #fff">mdi-close</v-icon>
            </v-toolbar>
          <Etablissement :siret="prediction._id.siret"></Etablissement>
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
      dialog: false
    }
  },
  computed: {
    variationEffectif () {
      if (this.prediction.etablissement.effectif / this.prediction.etablissement.effectif_precedent > 1.05) {
        return 'high'
      }
      if (this.prediction.etablissement.effectif / this.prediction.etablissement.effectif_precedent < 0.95) {
        return 'down'
      }
      return 'none'
    }
  },
  methods: {
    upOrDown(before, after, treshold) {
      if (before == null || after == null) {
        return 'mdi-help-circle'
      } 
      if (after/before > 1+treshold) {
        return 'mdi-arrow-up'
      }
      if (after/before < 1-treshold) {
        return 'mdi-arrow-down'
      }
      return 'mdi-tilde'
    },
    upOrDownClass(before, after, treshold) {
      if (before == null || after == null) {
        return 'unknown'
      } 
      if (after/before > 1+treshold) {
        return 'high'
      }
      if (after/before < 1-treshold) {
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
div.entete {
  float: left;
  background: linear-gradient(270deg, rgba(119, 122, 170, 0.219), rgba(119, 122, 170, 0));
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
  color: rgb(16, 114, 16);
}
.down {
  color: rgb(139, 19, 19);
}
.unknown {
  color: rgb(150, 150, 150);
}
td {
  width: 80px;
}
.pointer {cursor: pointer;}
</style>
