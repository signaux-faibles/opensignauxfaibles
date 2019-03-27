<template>
  <div style="background: #fffe">
    <div>
      <v-container>
        <v-layout wrap>
          <v-flex
          xs12
          md6
          class="pa-3"
          style="font-size: 18px">
            SIRET <b>{{ siret }}</b> <br/>
            {{ sirene.nature_juridique }}<br/>
            Création: {{ printDate(sirene.debutactivite) }}
            <br/><br/>
            <b>{{ (sirene.adresse || [])[0] }} </b>
            <br
            v-if="(sirene.adresse || [])[0] != ''"
            />
            {{ (sirene.adresse || [])[1] }}
            <br
            v-if="(sirene.adresse || [])[1] != ''"
            />
            {{ (sirene.adresse || [])[2] }}
            <br
            v-if="(sirene.adresse || [])[2] != ''"
            />
            {{ (sirene.adresse || [])[3] }}
            <br
            v-if="(sirene.adresse || [])[3] != ''"
            />
            {{ (sirene.adresse || [])[4] }}
            <br
            v-if="(sirene.adresse || [])[4] != ''"
            />
            {{ (sirene.adresse || [])[5] }}
            <br
            v-if="(sirene.adresse || [])[5] != ''"
            />
            {{ (sirene.adresse || [])[6] }}
            <br/><br/>
            <v-divider/>
            <br/>
            <v-divider/>
            <br/>
            <b>{{ (naf.n1 || {})[((naf.n5to1 || {})[(sirene.ape || '')] || '')] }}</b><br/>
            {{ (naf.n5 || {})[(sirene.ape || '')] }}<br/>
            Code APE: {{ (sirene.ape || '') }}<br/>
          </v-flex>
          <v-flex xs12 md6 class="text-xs-right pa-3">
            <iframe :v-if="sirene.longitude" width="100%" height="360" frameborder="0" scrolling="no" marginheight="0" marginwidth="0" :src="'https://www.openstreetmap.org/export/embed.html?bbox=' + (sirene.longitude - 0.03) + '%2C' + (sirene.lattitude  - 0.03) + '%2C' + (sirene.longitude + 0.03) + '%2C' + (sirene.lattitude + 0.03) + '&amp;layer=mapnik&amp;marker=' + sirene.lattitude + '%2C' + sirene.longitude" style="border: 1px solid black"></iframe><br/><small><a href="https://www.openstreetmap.org/#map=19/47.31581/5.05088">Afficher une carte plus grande</a></small>
          </v-flex>
          <v-flex xs12>
              <v-toolbar
                class="mb-2"
                color="indigo darken-5"
                dark
                flat
              >
                <v-toolbar-title>Informations Financières</v-toolbar-title>
              </v-toolbar>
            <v-data-iterator
            :items="zipDianeBDF"
            :rows-per-page-items="[3]"
            :pagination.sync="pagination"
            row
            wrap>
              <v-flex
                class="pa-1"
                slot="item"
                slot-scope="props"
                xs6
                sm3
                md3
                lg3
              >
                <v-card
                outline
                class="elevation-2">
                  <v-card-title class="subheading font-weight-bold">{{ (props.item.diane[0] || {}).exercice_diane }}</v-card-title>

                  <v-divider></v-divider>

                  <v-list dense>
                    <v-list-tile>
                      <v-list-tile-content>Diane | Chiffre d'Affaire:</v-list-tile-content>
                      <v-list-tile-content
                      :class="((props.item.diane[0]||{})['CA'])?'align-end':'nc align-end'"
                      >
                      {{ ((props.item.diane[0]||{})['CA'])||'n/c ' }}k€</v-list-tile-content>
                    </v-list-tile>
                    <v-list-tile>
                      <v-list-tile-content>Diane | Rentabilité Nette</v-list-tile-content>
                      <v-list-tile-content
                      :class="(props.item.diane[0]||{})['rentabilite_nette_pourcent']?'align-end':'nc align-end'"
                      >
                      {{ (props.item.diane[0]||{})['rentabilite_nette_pourcent'] || 'n/c' }} %</v-list-tile-content>
                    </v-list-tile>
                    <v-list-tile>
                      <v-list-tile-content>Diane | Résultat d'exploitation :</v-list-tile-content>
                      <v-list-tile-content
                      :class="(props.item.diane[0]||{})['resultat_expl']?'align-end':'nc align-end'"
                      >
                      {{ (props.item.diane[0]||{})['resultat_expl'] || 'n/c' }} k€</v-list-tile-content>
                    </v-list-tile>
                    <v-list-tile>
                      <v-list-tile-content>Diane | Résultat Net Consolidé</v-list-tile-content>
                      <v-list-tile-content
                      :class="(props.item.diane[0]||{})['resultat_net_consolide']?'align-end':'nc align-end'"
                      >
                      {{ (props.item.diane[0]||{})['resultat_net_consolide'] || 'n/c' }} k€</v-list-tile-content>
                    </v-list-tile>
                    <v-list-tile>
                      <v-list-tile-content>Diane | Valeur Ajoutée:</v-list-tile-content>
                      <v-list-tile-content
                      :class="((props.item.diane[0]||{})['valeur_ajoutee'])?'align-end':'nc align-end'"
                      >
                      {{ (props.item.diane[0]||{})['valeur_ajoutee'] || 'n/c' }} k€</v-list-tile-content>
                    </v-list-tile>
                  </v-list>
                </v-card>
              </v-flex>
            </v-data-iterator>
          </v-flex>
          <v-flex xs12>
            <v-toolbar
              class="mb-2"
              color="indigo darken-5"
              dark
              flat
            >
              <v-toolbar-title>Effectifs</v-toolbar-title>
            </v-toolbar>
          </v-flex>
          <v-flex xs12 style="height: 350px">
            <IEcharts
              :loading="chart"
              style="height: 350px"
              :option="effectifOptions(effectif)"
            />
          </v-flex>
          <v-flex xs12>
            <v-toolbar
            dark
            color='indigo darken-5'>
              <v-toolbar-title>Débits Urssaf</v-toolbar-title>
            </v-toolbar>
            <IEcharts
              :loading="chart"
              style="height: 350px"
              :option="urssafOptions"
            />
          </v-flex>
          <v-flex xs6 class="pr-1">
            <v-toolbar
            dark
            color='indigo darken-5'>
              <v-toolbar-title>Demandes d'activité partielle</v-toolbar-title>
            </v-toolbar>
            <v-list>
              <v-list-tile>
                <v-list-tile-content class="text-xs-right" style="width: '25%'">
                  Date
                </v-list-tile-content>
                <v-list-tile-content class="text-xs-right" style="width: '25%'">
                  Effectif Autorisé
                </v-list-tile-content>
                <v-list-tile-content class="text-xs-right" style="width: '25%'">
                  Début
                </v-list-tile-content>
                <v-list-tile-content class="text-xs-right" style="width: '25%'">
                  Fin
                </v-list-tile-content>
              </v-list-tile>
              <v-list-tile
                v-for="(d, i) in apdemande"
                :key="i">
                <v-list-tile-content class="text-xs-right" style="width: '25%'">
                  {{ d.date_statut.substring(0,10) }}
                </v-list-tile-content>
                <v-list-tile-content class="text-xs-right" style="width: '25%'">
                  {{ d.effectif_autorise }}
                </v-list-tile-content>
                <v-list-tile-content class="text-xs-right" style="width: '25%'">
                  {{ d.periode.start.substring(0,10) }}
                </v-list-tile-content>
                <v-list-tile-content class="text-xs-right" style="width: '25%'">
                  {{ d.periode.end.substring(0,10) }}
                </v-list-tile-content>
              </v-list-tile>
            </v-list>
          </v-flex>
          <v-flex xs6 class="pl-1">
            <v-toolbar
            dark
            color='indigo darken-5'>
              <v-toolbar-title>Consommations d'activité partielle</v-toolbar-title>
            </v-toolbar>
            <v-list style="width: 100%">
              <v-list-tile>
                <v-list-tile-content class="align-right" style="width: '33%'">
                  Date
                </v-list-tile-content>
                <v-list-tile-content class="align-right" style="width: '33%'">
                  Effectifs
                </v-list-tile-content>
                <v-list-tile-content class="align-right" style="width: '33%'">
                  Heures
                </v-list-tile-content>
              </v-list-tile>
              <v-list-tile
                v-for="(d, i) in apconso"
                :key="i">
                <v-list-tile-content class="align-right" style="width: '25%'">
                  {{ d.periode.substring(0, 10) }}
                </v-list-tile-content>
                <v-list-tile-content class="align-right" style="width: '25%'">
                  {{ d.effectif }}
                </v-list-tile-content>
                <v-list-tile-content class="align-right" style="width: '25%'">
                  {{ d.montant }}
                </v-list-tile-content>
              </v-list-tile>
            </v-list>
          </v-flex>
        </v-layout>
      </v-container>
    </div>
  </div>
</template>

<script>
import IEcharts from 'vue-echarts-v3/src/lite.js'
import 'echarts/lib/chart/line'
import 'echarts/lib/component/title'

export default {
  props: ['siret', 'batch'],
  name: 'Etablissement',
  data () {
    return {
      chart: false,
      bilan: true,
      urssaf: true,
      apart: true,
      etablissement: {},
      entreprise: {},
      pagination: null,
      naf: {}
    }
  },
  methods: {
    close () {
      this.tabs = this.tabs.filter((tab, index) => index !== this.activeTab)
      this.activeTab = this.activeTab - 1
    },
    getEtablissement (val) {
      var params = {
        batch: this.currentBatchKey,
        siret: val
      }
      this.$axios.post('/api/data/etablissement', params ).then(response => {
        this.etablissement = response.data[0]
      })
    },
    printDate (date) {
      return (date || '          ').substring(0, 10)
    },
    round (value, size) {
      return Math.round(value * (10 ^ size)) / (10 ^ size)
    },
    effectifOptions (effectif) {
      return {
        title: {
          text: null
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            label: {
              backgroundColor: '#283b56'
            }
          }
        },
        toolbox: {
          show: true
        },
        xAxis: {
          show: true,
          type: 'category',
          axisTick: false,
          data: this.effectif.map(e => e.periode)
        },
        yAxis: {
          type: 'value',
          show: true
        },
        series: [{
          color: 'indigo',
          smooth: true,
          name: 'taux marge',
          type: 'line',
          data: this.effectif.map(e => e.effectif)
        }]
      }
    }
  },
  mounted () {
    this.$axios.get('/api/data/naf').then(response => { this.naf = response.data })
    this.getEtablissement(this.siret)
  },
  components: {
    IEcharts
  },
  watch: {
    localSiret: function (val) {
      this.getEtablissement(val)
    }
  },  
  computed: {
    localSiret () {
      return this.siret
    },
    apconso () {
      return ((this.etablissement || {}).apconso || []).sort((a, b) => a.periode <= b.periode).slice(0, 10)
    },
    apdemande () {
      return ((this.etablissement || {}).apdemande || []).sort((a, b) => a.periode.start <= b.periode.start).slice(0, 10)
    },
    activeTab: {
      get () { return this.$store.getters.activeTab },
      set (activeTab) { this.$store.dispatch('updateActiveTab', activeTab) }
    },
    tabs: {
      get () { return this.$store.getters.getTabs },
      set (tabs) { this.$store.dispatch('updateTabs', tabs) }
    },
    sirene () {
      return (this.etablissement.value || {'sirene': {}}).sirene
    },
    effectif () {
      return ((this.etablissement.effectif || []) || []).sort((a, b) => a.periode < b.periode).slice(0, 15).reverse()
    },
    bdf () {
      return (this.entreprise.bdf || []).sort((a, b) => a.annee > b.annee).reverse()
    },
    diane () {
      return (Object.keys(this.entreprise.diane || []).map(k => this.entreprise.diane[k]) || [null]).sort((a, b) => a.annee > b.annee).reverse()
    },
    currentBatchKey () {
      return this.$store.state.currentBatchKey
    },
    zipDianeBDF () {
      let annees = new Set(this.bdf.map(b => b.annee).concat(this.diane.map(d => d.annee)))
      return Array.from(annees).sort((a, b) => a < b).map(a => {
        return {
          annee: a,
          bdf: this.bdf.filter(b => b.annee === a),
          diane: this.diane.filter(d => d.annee === a)
        }
      })
    },
    urssafOptions () {
      return {
        title: {
          text: null
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            label: {
              backgroundColor: '#283b56'
            }
          }
        },
        toolbox: {
          show: true
        },
        xAxis: {
          show: true,
          type: 'category',
          axisTick: false,
          data: (this.etablissement.array_debit || []).map(d => d.periode)
        },
        yAxis: {
          type: 'value',
          show: true
        },
        series: [{
          color: 'indigo',
          smooth: true,
          name: 'Cotisation',
          type: 'line',
          data: (this.etablissement.array_debit || []).map(d => d.cotisation)
        }, {
          color: 'red',
          smooth: true,
          name: 'Dette URSSAF',
          type: 'line',
          data: (this.etablissement.array_debit || []).map(d => d.montant_part_ouvriere + d.montant_part_patronale)
        }]
      }
    }
  }
}
</script>

<style scoped>
.echarts {
  width: 400px
}
.widget {
  position: absolute;
  left: 20px;
  top: 20px;
  right: 20px;
}
.nc {
  color: #bbb;
}
</style>
