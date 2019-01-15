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
      Détection
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
        <v-list class="pa-0">
          <v-list-tile avatar>
            <v-list-tile-avatar>
              <v-icon :class="loading?'rotate':''" @click="rightDrawer=!rightDrawer">mdi-target</v-icon>
            </v-list-tile-avatar>
            <v-spacer></v-spacer>
            <v-img src="/static/regions/PDL.svg"></v-img>
          </v-list-tile>
        </v-list>
      </v-toolbar>
      <v-list two-line>
        <v-list-tile three-line>
          <v-select
            :items="batches"
            v-model="currentBatchKey"
            label="Lot d'intégration"
            @change="updatePrediction()"
          ></v-select>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
            <v-icon>fa-industry</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-select
              :items="naf1"
              v-model="naf"
              label="Secteur d'activité"
              @change="updatePrediction()"
            ></v-select>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
            <v-icon>fa-users</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-select
              :items="effectifClass"
              v-model="minEffectif"
              label="Effectif minimum"
              @change="updatePrediction()"
            ></v-select>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
          <v-checkbox
          @change="updatePrediction()"
            v-model="entrepriseConnue">
          </v-checkbox>
          </v-list-tile-action>
          <v-list-tile-content>
            Entreprise non suivie
          </v-list-tile-content>
        </v-list-tile>
        <!-- <v-list-tile>
          <v-list-tile-action>
          <v-checkbox
          @change="updatePrediction()"
            v-model="horsCCSF">
          </v-checkbox>
          </v-list-tile-action>
          <v-list-tile-content>
            hors CCSF
          </v-list-tile-content>
        </v-list-tile> -->
        <!-- <v-list-tile>
          <v-list-tile-action>
          <v-checkbox
          @change="updatePrediction()"
            v-model="horsProcol">
          </v-checkbox>
          </v-list-tile-action>
          <v-list-tile-content>
            hors Procédure Collective
          </v-list-tile-content>
        </v-list-tile> -->
      </v-list>
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
      entrepriseConnue: true,
      horsCCSF: true,
      horsProcol: true,
      loading: false
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
        suivi: this.entrepriseConnue
      }
      this.$axios.post('/api/data/prediction', params).then(response => {
        var prediction = response.data
        prediction.forEach(p => {
          p.bdf = Object.keys(p.bdf || {})
            .map(b => p.bdf[b])
            .sort((a, b) => a.annee < b.annee)
        })
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
        suivi: this.entrepriseConnue
      }
      this.predictionLength = limit + offset
      this.$axios.post('/api/data/prediction', params).then(response => {
        var prediction = response.data
        prediction.forEach(p => {
          p.bdf = Object.keys(p.bdf || {})
            .map(b => p.bdf[b])
            .sort((a, b) => a.annee < b.annee)
          p.entreprise = {diane: []}
        })
        this.prediction = this.prediction.concat(prediction)
        self.loading = false
      })
    }
  },
  computed: {
    naf1 () {
      return Object.keys(this.$store.state.naf.n1 || {}).sort().map(n => {return {
        text: this.$store.state.naf.n1[n].substring(0, 60),
        value: n
      }})
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

<style scoper>
.rotate {
    -webkit-animation:spin 4s linear infinite;
    -moz-animation:spin 4s linear infinite;
    animation:spin 4s linear infinite;
}
@-moz-keyframes spin { 100% { -moz-transform: rotate(360deg); } }
@-webkit-keyframes spin { 100% { -webkit-transform: rotate(360deg); } }
@keyframes spin { 100% { -webkit-transform: rotate(360deg); transform:rotate(360deg); } }
</style>
