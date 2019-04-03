<template>
  <div style="background: #fffe">
    <div>
      <v-container>
        <v-layout wrap>
          <v-flex
          style="vertical-align: middle; text-align: center; margin-top: 30px; font-size: 24px"
          xs12>
            {{ sirene.raison_sociale }}
          </v-flex>
          <v-flex
          xs12
          md6
          class="pa-3"
          style="font-size: 18px">
            <!-- {{ Object.keys(etablissement) }} -->
            SIRET <b>{{ siret }} – {{ sirene.nature_juridique }}</b> <br/>

            <v-btn :color="suivi?'error':'success'" @click="suivi=!suivi">
              {{ suivi?'Ne plus suivre cet établissement':'Suivre cet établissement' }}
            </v-btn>            

            <!-- Création: {{ printDate(sirene.debutactivite) }} -->
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
            <iframe 
            :v-if="sirene.longitude" 
            width="100%" 
            height="360" 
            frameborder="0" 
            scrolling="no" 
            marginheight="0" 
            marginwidth="0" 
            :src="'https://www.openstreetmap.org/export/embed.html?bbox=' + (sirene.longitude - 0.03) + '%2C' + (sirene.lattitude  - 0.03) + '%2C' + (sirene.longitude + 0.03) + '%2C' + (sirene.lattitude + 0.03) + '&amp;layer=mapnik&amp;marker=' + sirene.lattitude + '%2C' + sirene.longitude" style="border: 1px solid black"></iframe><br/>
            <small><a href="https://www.openstreetmap.org/#map=19/47.31581/5.05088">Afficher une carte plus grande</a></small>
          </v-flex>
          

          <v-flex 
            xs12 md12 class="pa-3"
            v-for="(c, i) in comments"
            :key="comment + i">
            <v-textarea
              box
              :label="c.author + ', le ' + c.date"
              :value="c.comment"
            ></v-textarea>
          </v-flex>
          <v-flex xs6 class="pr-1" style="height: 200px">
            <v-toolbar
              dark
              color='indigo darken-5'>
              <v-toolbar-title class="localtoolbar">Effectifs</v-toolbar-title>
            </v-toolbar>
            <IEcharts
              :loading="chart"
              style="height: 350px"
              :option="effectifOptions"
            />
          </v-flex>

          <v-flex xs6 class="pr-1">
            <v-toolbar
              dark
              color='indigo darken-5'>
              <v-toolbar-title class="localtoolbar">Débits Urssaf</v-toolbar-title>
            </v-toolbar>
            <IEcharts
              :loading="chart"
              style="height: 350px"
              :option="urssafOptions"
            />
          </v-flex>

          <v-flex xs6 class="pr-1" v-if="apdemande.length + apconso.length > 0">
            <v-toolbar
              dark
              color='indigo darken-5'>
              <v-toolbar-title class="localtoolbar">Demandes d'activité partielle</v-toolbar-title>
            </v-toolbar>
            <table style="width: 100%; padding: 5px; text-align: center">
              <tr>
                <th>
                  Date
                </th>
                <th>
                  Effectif Autorisé
                </th>
                <th>
                  Début
                </th>
                <th>
                  Fin
                </th>
              </tr>
              <tr v-for="d in apdemande" :key="d.date_statut" style="font-weight: 400">
                <td>
                  {{ d.date_statut.substring(0,10) }}
                </td>
                <td>
                  {{ d.effectif_autorise }}
                </td>
                <td>
                  {{ d.periode.start.substring(0,10) }}
                </td>
                <td>
                  {{ d.periode.end.substring(0,10) }}
                </td>
              </tr>
            </table>

          </v-flex>

          <v-flex xs6 class="pr-1" v-if="apdemande.length + apconso.length > 0">
            <v-toolbar
            dark
            color='indigo darken-5'>
              <v-toolbar-title class="localtoolbar">Consommations d'activité partielle</v-toolbar-title>
            </v-toolbar>
            <table style="width: 100%; padding: 5px; text-align: center">
              <tr>
                <th>
                  Date
                </th>
                <th>
                  Effectifs
                </th>
                <th>
                  Heures
                </th>
              </tr>
              <tr v-for="d in apconso" :key="d.periode" style="font-weight: 400">
                <td>
                  {{ d.periode.substring(0, 10) }}
                </td>
                <td>
                  {{ d.effectif }}
                </td>
                <td>
                  {{ d.heure_consomme }}
                </td>
              </tr>
            </table>
          </v-flex>

          <v-flex xs12 class="pr-1">
            <v-toolbar
              class="mb-2"
              color="indigo darken-5"
              dark
            >
              <v-toolbar-title class="localtoolbar">Informations Financières</v-toolbar-title>
            </v-toolbar>
          </v-flex>
          <v-flex
            v-for="f in finance.slice(0,3)"
            :key="f.annee"
            pa-1
            xs12
            sm6
            md4
            lg4
          >
            <v-card
            outline
            class="elevation-2">
              <v-card-title class="subheading font-weight-bold">{{ f.annee }}</v-card-title>

              <v-divider></v-divider>

              <v-list dense>
                <v-list-tile>
                  <v-list-tile-content>Chiffre d'Affaire:</v-list-tile-content>
                  <v-list-tile-content class="text-xs-right">
                    <div style="width: 100%;" :class="f.ca_class">
                    {{ f.ca }}
                    </div>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>

              <v-list dense>
                <v-list-tile>
                  <v-list-tile-content>Résultat d'exploitation:</v-list-tile-content>
                  <v-list-tile-content class="text-xs-right">
                    <div style="width: 100%;" :class="f.marge_ope_class">
                      {{ f.resultat_expl }} ({{ (f.marge_ope) }})
                    </div>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>

              <v-list dense>
                <v-list-tile>
                  <v-list-tile-content>Résultat net:</v-list-tile-content>
                  <v-list-tile-content class="text-xs-right">
                    <div style="width: 100%;" :class="f.marge_nette_class">
                      {{ f.benefice_ou_perte }} ({{ (f.marge_nette) }})
                    </div>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>

              <v-list dense>
                <v-list-tile>
                  <v-list-tile-content>Délai Fournisseur:</v-list-tile-content>
                  <v-list-tile-content class="text-xs-right">
                    <div style="width: 100%;" :class="f.delai_fournisseur_class">
                      {{ f.delai_fournisseur }}
                    </div>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>

              <v-list dense>
                <v-list-tile>
                  <v-list-tile-content>Délai Client:</v-list-tile-content>
                  <v-list-tile-content class="text-xs-right">
                    <div style="width: 100%;" :class="f.delai_client_class">
                      {{ f.delai_client }}
                    </div>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>

              <v-list dense>
                <v-list-tile>
                  <v-list-tile-content>Poids FRNG:</v-list-tile-content>
                  <v-list-tile-content class="text-xs-right">
                    <div style="width: 100%;" :class="f.poids_frng_class">
                      {{ f.poids_frng }}
                    </div>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>

              <v-list dense>
                <v-list-tile>
                  <v-list-tile-content>Financement court terme:</v-list-tile-content>
                  <v-list-tile-content class="text-xs-right">
                    <div style="width: 100%;" :class="f.financier_court_terme_class">
                      {{ f.financier_court_terme }}
                    </div>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>

            </v-card>
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
import 'echarts/lib/component/legend'
import 'echarts/lib/component/tooltip'

export default {
  props: ['siret', 'batch'],
  name: 'Etablissement',
  data () {
    return {
      suivi: false,
      chart: false,
      bilan: true,
      urssaf: true,
      apart: true,
      etablissement: {value: {}},
      pagination: null,
      comments: []
    }
  },
  methods: {
    computeFinance(f) {
      var annee = f.annee
      var ca = f.diane.ca?f.diane.ca + ' k€':'n/c'
      var ca_class = (!f.diane.ca)?'gray':''
      var resultat_expl = f.diane.resultat_expl?f.diane.resultat_expl + ' k€':'n/c'
      var marge_ope = f.diane.resultat_expl / f.diane.ca 
      var marge_ope_class = (!marge_ope)?'gray':(marge_ope<0)?'down':''
      marge_ope = marge_ope?(Math.floor(marge_ope*1000)/10) + ' %':'n/c'

      var benefice_ou_perte = f.diane.benefice_ou_perte?f.diane.benefice_ou_perte + ' k€':'n/c'
      var marge_nette = f.diane.benefice_ou_perte / f.diane.ca 
      var marge_nette_class = !(marge_nette)?'gray':(marge_nette<0)?'down':''
      marge_nette = marge_nette?(Math.floor(marge_nette*1000)/10) + ' %':'n/c'

      var delai_fournisseur = f.bdf.delai_fournisseur?Math.round(f.bdf.delai_fournisseur) + ' jours':'n/c'
      var delai_fournisseur_class = !(f.bdf.delai_fournisseur)?'gray':''
      var delai_client = Math.round(f.diane.credit_client / f.diane.ca * 360)
      var delai_client_class = !(delai_client)?'gray':''

      var delai_client = delai_client?delai_client + ' jours':'n/c' 

      var poids_frng = f.bdf.poids_frng?Math.round(f.bdf.poids_frng*10)/10 + ' %':'n/c'
      var poids_frng_class = !(f.bdf.poids_frng)?'gray':''

      var financier_court_terme = f.bdf.financier_court_terme?Math.round(f.bdf.financier_court_terme*10)/10 + ' %':'n/c'
      var financier_court_terme_class = !(f.bdf.financier_court_terme)?'gray':''

      return {
        annee,
        ca,
        ca_class,
        resultat_expl,
        marge_ope,
        marge_ope_class,
        benefice_ou_perte,
        marge_nette,
        marge_nette_class,
        delai_fournisseur,
        delai_fournisseur_class,
        delai_client,
        delai_client_class,
        poids_frng,
        poids_frng_class,
        financier_court_terme,
        financier_court_terme_class
      }
    },
    addComment() {
      this.comments.push({
        'comment': '',
        'author': 'C.Ninucci',
        'date': this.formattedDate(new Date)
        })
    },
    formattedDate (d) {
      let month = String(d.getMonth() + 1);
      let day = String(d.getDate());
      const year = String(d.getFullYear());

      if (month.length < 2) month = '0' + month;
      if (day.length < 2) day = '0' + day;

      return `${day}/${month}/${year}`
    },
    close () {
      this.tabs = this.tabs.filter((tab, index) => index !== this.activeTab)
      this.activeTab = this.activeTab - 1
    },
    getEtablissement (val) {
      var params = {
        batch: '1802',
        siret: val
      }
      this.$axios.post('/api/data/etablissement', params).then(response => {
        this.etablissement = response.data[0]
      })
    },
    printDate (date) {
      return (date || '          ').substring(0, 10)
    },
    round (value, size) {
      return Math.round(value * (10 ^ size)) / (10 ^ size)
    }
  },
  mounted () {
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
    finance () {
      return this.zipDianeBDF.filter(z => z.annee).map(z => this.computeFinance(z))
    },
    naf () {
      return this.$store.state.naf
    },
    localSiret () {
      return this.siret
    },
    apconso () {
      return ((this.etablissement.value || {}).apconso || []).sort((a, b) => a.periode <= b.periode).slice(0, 10)
    },
    apdemande () {
      return ((this.etablissement.value || {}).apdemande || []).sort((a, b) => a.periode.start <= b.periode.start).slice(0, 10)
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
      return this.etablissement.value.sirene || {}
    },
    debit () {
      return this.etablissement.value.debit || []
    },
    cotisation () {
      return this.etablissement.value.cotisation || []
    },
    effectif () {
      return this.etablissement.value.effectif || []
    },
    bdf () {
      if (this.etablissement.entreprise) {
        return this.etablissement.entreprise.value.bdf
      } else {
        return []
      }
    },
    diane () {
      if (this.etablissement.entreprise) {
        return this.etablissement.entreprise.value.diane
      } else {
        return []
      }
    },
    currentBatchKey () {
      return this.$store.state.currentBatchKey
    },
    zipDianeBDF () {
      let annees = new Set(this.bdf.map(b => b.annee_bdf).concat(this.diane.map(d => d.exercice_diane)))
      return Array.from(annees).sort((a, b) => a < b).map(a => {
        return {
          annee: a,
          bdf: this.bdf.filter(b => b.annee_bdf === a)[0] || {},
          diane: this.diane.filter(d => d.exercice_diane === a)[0] || {}
        }
      })
    },
    effectifOptions () {
      return {
        title: {
          text: null
        },
        legend: {
            data:['effectif'],
            y: 'bottom'
        },
        toolbox: {
          show: true
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
        xAxis: {
          show: true,
          type: 'category',
          data: this.effectif.map(e => e.periode.slice(0,10))
        },
        yAxis: {
          type: 'value',
          show: true
        },
        series: [{
          name: 'effectif',
          color: 'indigo',
          step: 'end',
          type: 'line',
          data: this.effectif.map(e => e.effectif)
        }]
      }
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
        legend: {
          data: ['Cotisation', 'Débit'],
          y: 'bottom'
        },
        toolbox: {
          show: true
        },
        xAxis: {
          show: true,
          type: 'category',
          axisTick: false,
          data: this.debit.map(d => (d.periode||'').slice(0,10))
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
          data: this.cotisation,
        }, {
          color: 'red',
          smooth: true,
          name: 'Débit',
          type: 'line',
          data: this.debit.map(d => d.part_ouvriere + d.part_patronale)
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
.gray {
  color: #aaa;
}
.down {
  color: rgb(244, 67, 54);
}
.widget {
  position: absolute;
  left: 20px;
  top: 20px;
  right: 20px;
}
.localtoolbar {
  font-weight: 600;
  text-align: center;
  width: 100%;
}
.nc {
  color: #bbb;
}
.align-right > * {
  text-align: center;
}
</style>
