<template>
  <div>
    <v-toolbar height="35px" class="toolbar" color="#ffffff"  app>
      <v-icon
      @click="drawer=!drawer"
      class="fa-rotate-180"
      v-if="!drawer"
      color="#fff"
      key="toolbar"
      >
      mdi-backburger
      </v-icon>
      <div style="width: 100%; text-align: center;" class="toolbar_titre">
        Données
      </div>
      <v-spacer></v-spacer>
      <v-icon color="#fff" v-if="!rightDrawer" @click="rightDrawer=!rightDrawer">fa-database</v-icon>
    </v-toolbar>
    <v-navigation-drawer
      :class="rightDrawer?'elevation-6':''"
      right app
      v-model="rightDrawer"
    >
      <v-list  two-line class="pt-0">
        <v-toolbar>
          <v-icon @click="rightDrawer=!rightDrawer">fa-database</v-icon>
          <v-spacer></v-spacer>
          Lot {{ currentBatchKey }}
        </v-toolbar>
        <v-list-tile three-line>
          <v-select
            :items="batchesKeys"
            v-model="currentBatchKey"
            label="Lot d'intégration"
          ></v-select>
        </v-list-tile>
        <v-list-group v-if="currentBatchKey != null">
          <v-list-tile slot="activator" bgcolor="red">
            <v-list-tile-action>
              <v-icon>fa-cogs</v-icon>
            </v-list-tile-action>
            <v-list-tile-content class="title">
              Paramètres
            </v-list-tile-content>
          </v-list-tile>
          <v-list-tile
          v-for="param in parameters"
          :key="currentBatchKey + param.key"
          ripple
          @click="setCurrentType(param.key)">
            <v-list-tile-content>
              <v-list-tile-title
              :class="(param.key===currentType) ? 'selected': null"
              >{{ param.text }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list-group>
        <v-list-group v-if="currentBatchKey != null">
          <v-list-tile slot="activator">
            <v-list-tile-action>
              <v-icon>fa-copy</v-icon>
            </v-list-tile-action>
            <v-list-tile-title class="title">
              Fichiers
            </v-list-tile-title>
          </v-list-tile>
          <v-list-tile
          v-for="type in types"
          :key="type.text"
          @click="setCurrentType(type.type)"
          >
            <v-list-tile-content>
              <v-list-tile-title
              :class="(type.type==currentType) ? 'selected': null">
              {{ type.text }}
              </v-list-tile-title>
            </v-list-tile-content>
            <v-list-tile-action>
              <v-icon
              @click="toggleComplete(type.type)"
              >{{ currentBatch.complete_types.includes(type.type)?'mdi-square-inc':'mdi-shape-square-plus' }}</v-icon>
            </v-list-tile-action>
          </v-list-tile>
          <v-divider></v-divider>
        </v-list-group>
        <v-list-group v-if="currentBatchKey != null">
          <v-list-tile slot="activator">
            <v-list-tile-action>
              <v-icon>fa-microchip</v-icon>
            </v-list-tile-action>
            <v-list-tile-title class="title">
              Traitements
            </v-list-tile-title>
          </v-list-tile>
          <v-divider></v-divider>
          <v-list-tile
          v-for="process in processes"
          :key="currentBatchKey + process.key"
          @click="setCurrentType(process.key)"
          >
            <v-list-tile-content>
              <v-list-tile-title
              :class="(process.key==currentType) ? 'selected': null">
              {{ process.text }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list-group>
      </v-list>
    </v-navigation-drawer>
    <div style="width:100%">
      <div v-if="currentBatchKey != null">
        <BatchDate
        class="elevation-6"
        :key="currentBatchKey + 'batchDate'"
        :date="currentType"
        :param="parameters.filter(p => p.key === currentType)[0]"
        v-if="parameters.map(p => p.key).includes(currentType)"
        />
        <BatchFile
        :key="currentBatchKey + 'batchFile'"
        :type="currentType"
        v-if="types.map(t => t.type).includes(currentType)"
        />
        <BatchProcess
        :key="currentBatchKey + 'batchProcess'"
        :process="processes.filter(p => p.key === currentType)[0]"
        v-if="processes.map(p => p.key).includes(currentType)"
        />
      </div>
    </div>
  </div>
</template>

<script>
import BatchFile from '@/components/BatchFile'
import BatchDate from '@/components/BatchDate'
import BatchProcess from '@/components/BatchProcess'

// import Batch from '@/components/Batch'

export default {
  mounted () {
    this.$store.dispatch('updateBatches')
    this.$store.dispatch('updateRefs')
  },
  methods: {
    setCurrentBatchKey (batchKey) {
      this.currentBatchKey = batchKey
    },
    setCurrentType (type) {
      this.currentType = type
    }
  },
  computed: {
    types () {
      return this.$store.state.types
    },
    currentType: {
      get () {
        return this.$store.state.currentType
      },
      set (type) {
        this.$store.dispatch('setCurrentType', type)
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
        if (this.currentBatchKey == null) {
          return true
        } else {
          return this.$store.state.rightDrawer
        }
      },
      set (val) {
        this.$store.dispatch('setRightDrawer', val)
      }
    },
    currentBatch: {
      get () {
        if (this.currentBatchKey in this.$store.getters.batchesObject) {
          return this.$store.getters.batchesObject[this.currentBatchKey]
        } else {
          return { 'complete_types': [] }
        }
      },
      set (batch) {
        this.$store.dispatch('saveBatch', batch).then(r => this.$store.dispatch('checkEpoch'))
      }
    },
    currentBatchKey: {
      get () {
        return this.$store.state.currentBatchKey
      },
      set (value) {
        this.$store.commit('setCurrentBatchKey', value)
      }
    },
    batchesKeys () {
      return (this.$store.state.batches || []).map(batch => batch.id.key)
    }
  },
  data () {
    return {
      parameters: [
        { text: 'Date de début',
          key: 'dateDebut',
          prop: 'date_debut',
          desc: 'Début de la période des paramètres' },
        { text: 'Date de fin',
          key: 'dateFin',
          prop: 'date_fin',
          desc: 'Fin de la période des paramètres' },
        { text: 'Date de fin effectifs',
          key: 'dateFinEffectif',
          prop: 'date_fin_effectif',
          desc: 'Dernière date du fichier effectif' }
      ],
      processes: [
        { text: 'Suppression',
          color: 'red',
          key: 'reset',
          img: '/static/poubelle.png',
          description: 'Retour au batch précédent',
          do (self) { self.$axios.get('/api/batch/revert') }
        },
        { text: 'Purger',
          color: 'blue',
          key: 'purge',
          img: '/static/gomme.svg',
          description: 'Retour au paramétrage',
          do (self) { self.$axios.get('/api/batch/purge') }
        },
        { text: 'Calcul Prédictions',
          color: 'green',
          key: 'predict',
          img: '/static/warning.png',
          description: 'Intégration des données et calcul des prédictions.',
          do (self) { self.$axios.get('/api/batch/process') }
        },
        { text: 'Clôture',
          color: 'black',
          key: 'close',
          img: '/static/warning.png',
          description: 'Clôture du batch en cours et création du suivant',
          do (self) { self.$axios.get('/api/batch/next') }
        }
      ]
    }
  },
  components: { BatchDate, BatchFile, BatchProcess },
  name: 'Data'
}
</script>
