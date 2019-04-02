<template>
<div id='#detection'>
  <v-toolbar
    height="35px"
    class="toolbar elevation-12"
    color="#ffffff"
    extension-height="48px"
    app
  >
    <v-icon
      @click="drawer=!drawer"
      class="fa-rotate-180"
      v-if="!drawer"
      color="#ffffff"
      key="toolbar"
    >
      mdi-backburger
    </v-icon>
    
    <div style="width: 100%; text-align: center;" class="toolbar_titre">
      Détection - Bourgogne Franche Comté
    </div>
    <v-spacer></v-spacer>
    <v-icon
    :class="loading?'rotate':''"
    color="#ffffff" v-if="!rightDrawer" @click="rightDrawer=!rightDrawer">mdi-target</v-icon>
  </v-toolbar>
  <span style="visibility: hidden; position:absolute;">{{ detectionLength }} {{ predictionLength }} {{ prediction.length }}</span>
  <div style="width:100%">
    <v-navigation-drawer :class="(rightDrawer?'elevation-6':'') + 'rightDrawer'" v-model="rightDrawer" right app>
      <v-toolbar flat class="transparent" height="40">
        <v-icon :class="loading?'rotate':''" @click="rightDrawer=!rightDrawer">mdi-target</v-icon>
      </v-toolbar>
      <div style="width: 100%; padding: 0 15px;">
          <v-select
            :items="batches"
            v-model="currentBatchKey"
            label="Liste de détection"
            @change="updatePrediction()"
          ></v-select>
      </div>
      <p style="height: 1px; border: 1px solid #eee"/>
      <div style="display: flex; flex-direction: row; vertical-align: middle; padding: 0 15px;" >
        <v-icon style="margin-right: 10px;">fa-industry</v-icon>
        <v-select
          :items="naf1"
          v-model="naf"
          label="Secteur d'activité"
          @change="updatePrediction()"
        ></v-select>
      </div>
      <p style="height: 1px; border: 1px solid #eee"/>
      <div style="display: flex; flex-direction: row; vertical-align: middle; padding: 0 15px;" >
        <v-icon style="margin-right: 10px;">fa-industry</v-icon>
        <v-select
          :items="subzones"
          v-model="zone"
          label="Zone Géographique"
          @change="updatePrediction()"
        ></v-select>
      </div>
      <p style="height: 1px; border: 1px solid #eee"/>
      <div style="display: flex; flex-direction: row; vertical-align: middle; padding: 0 15px;" >
        <v-icon style="margin-right: 10px;">fa-users</v-icon>
        <v-select
          :items="effectifClass"
          v-model="minEffectif"
          label="Effectif minimum"
          @change="updatePrediction()"
        ></v-select>
      </div>
      <p style="height: 1px; border: 1px solid #eee"/>
      <v-radio-group 
        style="padding: 0 15px" 
        v-model="suivi" 
        @change="updatePrediction()"
        :mandatory="false">
        <v-radio label="Non suivies" value="false"></v-radio>
        <v-radio label="Suivies" value="true"></v-radio>
        <v-radio label="Toutes" value="null"></v-radio>
      </v-radio-group> 
      <p style="height: 1px; border: 1px solid #eee"/>
      <div style="display: flex; flex-direction: row; vertical-align: middle; padding: 0 15px;" >
        <v-checkbox
          label="Activité partielle"
          @change="updatePrediction()"
          v-model="activitePartielle">
        </v-checkbox>
      </div>
      <div style="display: flex; flex-direction: row; vertical-align: middle; padding: 0 15px;" >
        <v-checkbox
        label="Sans dette Urssaf"
        @change="updatePrediction()"
          v-model="interetUrssaf">
        </v-checkbox>
      </div>
    </v-navigation-drawer>
  </div>
  <PredictionWidget v-for="p in prediction" :key="p._id.siret" :prediction="p"/>

</div>
</template>

<script>
import PredictionWidget from '@/components/PredictionWidget'
export default {
  data () {
    return {
      effectifClass: [10, 20, 50, 100],
      prediction: [],
      predictionLength: 0,
      naf: 'C',
      minEffectif: 20,
      suivi: 'false',
      horsCCSF: true,
      horsProcol: true,
      loading: false,
      zone: ['21', '25', '39', '58', '70', '71', '89', '90'],
      subzones: [
        {
          text: 'Région complète',
          value:  ['21', '25', '39', '58', '70', '71', '89', '90'],
        },
        {
          text: 'Bourgogne',
          value:  ['21', '58', '71', '89'],
        },
        {
          text: 'Franche Comté',
          value:  ['25', '39', '70', '90'],
        },
        {
          text: 'Côte d\'or',
          value:  ['21'],
        },
        {
          text: 'Doubs',
          value:  ['25'],
        },
        {
          text: 'Jura',
          value:  ['39'],
        },
        {
          text: 'Nièvre',
          value:  ['58'],
        },
        {
          text: 'Haute-Saône',
          value:  ['70'],
        },
        {
          text: 'Saône-et-Loire',
          value:  ['71'],
        },
        {
          text: 'Yonne',
          value:  ['89'],
        },
        {
          text: 'Territoire de Belfort',
          value:  ['90'],
        },
      ]
    }
  },
  mounted () {
    this.$store.dispatch('getNAF')
    this.firstUpdatePrediction()
  },
  methods: {
    firstUpdatePrediction () {
      if (this.prediction.length < 1) {
        this.updatePrediction()
        setTimeout(this.firstUpdatePrediction, 5000)
      }
    },
    updatePrediction () {
      this.loading = true
      var self = this
      var params = {
        batch: this.currentBatchKey,
        naf1: this.naf,
        limit: this.detectionLength,
        offset: 0,
        effectif: this.minEffectif,
        ccsf: this.horsCCSF,
        procol: false,
        suivi: this.suiviQuery,
        zone: this.zone
      }
      
      this.$axios.post('/api/data/prediction', params).then(response => {
        var prediction = response.data
        this.prediction = prediction
        this.predictionLength = this.prediction.length
        self.loading = false
      })
    },
    getPrediction (limit, offset) {
      this.loading = true
      var self = this
      var params = {
        batch: this.currentBatchKey,
        naf1: this.naf,
        limit: limit,
        offset: offset,
        effectif: this.minEffectif,
        ccsf: this.horsCCSF,
        procol: false,
        suivi: this.suiviQuery,
        zone: this.zone
      }
      this.predictionLength = limit + offset
      this.$axios.post('/api/data/prediction', params).then(response => {
        var prediction = response.data
        var newPrediction = this.prediction.concat(prediction)
        this.prediction = newPrediction
        self.loading = false
      })
    }
  },
  computed: {
    suiviQuery () {
      if (this.suivi === 'false') {
        return false
      } else if (this.suivi === 'true') {
        return true
      } else {
        return undefined
      }
    },
    naf1 () {
      return Object.keys(this.$store.state.naf.n1 || {}).sort().map(n => {
        return {
          text: this.$store.state.naf.n1[n].substring(0, 60),
          value: n
        }
      })
    },
    scrollTop () {
      return this.$store.state.scrollTop
    },
    height: {
      get () {
        return this.$store.state.height
      },
      set (height) {
        this.$store.dispatch('setHeight', height)
      }
    },
    drawer: {
      get () {
        return this.$store.state.appDrawer
      },
      set (val) {
        this.$store.dispatch('setDrawer', val)
      }
    },
    rightDrawer: {
      get () {
        return this.$store.state.rightDrawer
      },
      set (val) {
        this.$store.dispatch('setRightDrawer', val)
      }
    },
    currentBatchKey: {
      get () {
        this.$store.dispatch('updateBatches')
        return this.$store.state.currentBatchKey
      },
      set (value) {
        this.$store.commit('setCurrentBatchKey', value)
      }
    },
    batches () {
      return [{
        "value": "1810",
        "text": "Octobre 2019"
      },{
        "value": "1811",
        "text": "Novembre 2018"
      },{
        "value": "1812",
        "text": "Décembre 2018"
      },{
        "value": "1901",
        "text": "Janvier 2018"
      },{
        "value": "1802",
        "text": "Février 2018"
      }]
      return (this.$store.state.batches || []).map(batch => batch.id.key)
    },
    detectionLength () {
      var length = Math.round((this.height + this.scrollTop) / 860 + 5) * 10
      if (length > this.predictionLength) {
        var complement = length - this.predictionLength
        this.getPrediction(complement, this.predictionLength)
      }
      return length
    }
  },
  components: { PredictionWidget },
  name: 'Browse'
}
</script>

<style scoped>

.rotate {
    -webkit-animation:spin 4s linear infinite;
    -moz-animation:spin 4s linear infinite;
    animation:spin 4s linear infinite;
}
@-moz-keyframes spin { 100% { -moz-transform: rotate(360deg); } }
@-webkit-keyframes spin { 100% { -webkit-transform: rotate(360deg); } }
@keyframes spin { 100% { -webkit-transform: rotate(360deg); transform:rotate(360deg); } }
</style>
