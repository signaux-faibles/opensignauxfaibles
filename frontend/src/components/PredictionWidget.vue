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
          <div style="white-space: nowrap; overflow: hidden; max-width: 400px; max-height:30px">
            <span style="font-size: 18px; color: #333; line-height: 10px; display: inline-block; font-family: 'Oswald'; max-width: '100px'">
              {{ prediction.etablissement.sirene.raison_sociale }}
              <br style="line-height: 10px;">
            </span>
          </div>
          <span style="font-size: 12px; color: #333; line-height: 10px;">
            {{ prediction._id.key }}
            <br style="line-height: 10px;">
          </span>
          <v-img
            style="position: absolute; left: 160px; bottom: 10px;"
            width="17"
            src="/static/gray_apart.svg"
          ></v-img>
          <v-img
            style="position: absolute; left: 90px; bottom: 10px;"
            width="57"
            :src="'/static/' + (prediction.etablissement.urssaf?'red':'gray') + '_urssaf.svg'"
          ></v-img>
          <div style="position: absolute; left: 195px; bottom: 4px; color: #333">
            <span
              :class="variationEffectif"
              style="font-size: 20px"
            >{{ prediction.etablissement.dernier_effectif.effectif || 'n/c' }}</span>
          </div>
          <div class="flex" style="position:absolute; left: 400px; top: 0px; bottom: 0px; right: 9px;">
            
            <div class = "label">
              <b>Exercice</b><br/>
              {{ (diane[0] || {'exercice_diane': '-'}).exercice_diane || '-' }}<br/>
              {{ (diane[1] || {'exercice_diane': '-'}).exercice_diane || '-' }}<br/>
              {{ (diane[2] || {'exercice_diane': '-'}).exercice_diane || '-' }}<br/>
            </div>

            <div class = "label">
              <b>Chiffre d'affaire (k€)</b><br/>
              {{ (diane[0] || {'ca': '-'}).ca || '-' }}<br/>
              {{ (diane[1] || {'ca': '-'}).ca || '-' }}<br/>
              {{ (diane[2] || {'ca': '-'}).ca || '-' }}<br/>
            </div>

            <div class = "label">
              <b>Dont export (k€)</b><br/>
              {{ (diane[0] || {'ca_exportation': '-'}).ca_exportation || '-' }}<br/>
              {{ (diane[1] || {'ca_exportation': '-'}).ca_exportation || '-' }}<br/>
              {{ (diane[2] || {'ca_exportation': '-'}).ca_exportation || '-' }}<br/>
            </div>

            <div class = "label">
              <b>Résultat net (k€)</b><br/>
              {{ (diane[0] || {'benefice_ou_perte': '-'}).benefice_ou_perte || '-' }}<br/>
              {{ (diane[1] || {'benefice_ou_perte': '-'}).benefice_ou_perte || '-' }}<br/>
              {{ (diane[2] || {'benefice_ou_perte': '-'}).benefice_ou_perte || '-' }}<br/>
            </div>
          </div>
        </div>
        <v-dialog attach="#detection" lazy fullscreen v-model="dialog">
          <div style="height: 100%; width: 100%;  font-weight: 800; font-family: 'Abel', sans;">
            <v-toolbar fixed class="toolbar" height="35px" style="color: #fff; font-size: 22px;">
              <v-spacer/>
              {{ prediction.etablissement.sirene.raison_sociale }}
              <v-spacer/>
              <v-icon @click="dialog=false" style="color: #fff">mdi-close</v-icon>
            </v-toolbar>
            <Etablissement :siret="prediction._id.key"></Etablissement>
          </div>
        </v-dialog>
      </div>
    </v-card>
  </div>
</template>

<script>
import Etablissement from "@/components/Etablissement";
import PredictionWidgetScore from "@/components/PredictionWidgetScore";

export default {
  props: ["prediction"],
  components: {
    PredictionWidgetScore,
    Etablissement
  },
  data() {
    return {
      dialog: false,
      expand: false,
      rowsPerPageItems: [4, 8, 12],
      pagination: {
        rowsPerPage: 4
      },
    };
  },
  computed: {
    variationEffectif() {
      var effectif = this.prediction.etablissement.effectif[
        this.prediction.etablissement.effectif.length - 1
      ];
      var effectif_precedent = this.prediction.etablissement.effectif[
        this.prediction.etablissement.effectif.length - 12
      ];
      if (effectif / effectif_precedent > 1.05) {
        return "high";
      }
      if (effectif / effectif_precedent < 0.95) {
        return "down";
      }
      return "none";
    },
    urssaf() {
      this.prediction.etablissement.dette.reduce((m, d) => {
        return m + d.part_patronale + d.part_ouvriere;
      }, 0);
    },
    diane() {
      return (this.prediction.entreprise || {diane: []}).diane
    }
  },
  methods: {
    upOrDown(before, after, treshold) {
      if (before == null || after == null) {
        return "mdi-help-circle";
      }
      if (after / before > 1 + treshold) {
        return "mdi-arrow-up";
      }
      if (after / before < 1 - treshold) {
        return "mdi-arrow-down";
      }
      return "mdi-tilde";
    },
    upOrDownClass(before, after, treshold) {
      if (before == null || after == null) {
        return "unknown";
      }
      if (after / before > 1 + treshold) {
        return "high";
      }
      if (after / before < 1 - treshold) {
        return "down";
      }
      return "none";
    },
    showEtablissement() {
      this.dialog = true;
    }
  }
};
</script>

<style scoped>
div.flex {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  font-size: 11px;
}
div.label {
  text-align: right;
  font-family: 'Quicksand';
  width: 130px;
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
.pointer {
  cursor: pointer;
}
</style>
