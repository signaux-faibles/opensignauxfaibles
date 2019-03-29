<template>
<div >
  <v-toolbar class="toolbar" color="#c9aec5" height="35px" app>
    <v-icon
    @click="drawer=!drawer"
    class="fa-rotate-180"
    v-if="!drawer"
    color="#ffffff"
    key="toolbar"
    >mdi-backburger
    </v-icon>
    <div style="width: 100%; text-align: center;" class="toolbar_titre">
      Consultation
    </div>
    <v-spacer></v-spacer>
    <!-- <v-icon color="#ffffff" @click="rightDrawer=!rightDrawer">mdi-database-search</v-icon> -->
  </v-toolbar>
    <v-container bg fill-height grid-list-md text-xs-center>
      <v-layout row wrap align-center>
        <v-flex :class="select?'small':'big'" xs12>
          <span class="fblue">Signaux</span>
          <span class="fblack">·</span>
          <span class="fred">Faibles</span></v-flex>
        <v-flex style="margin: 0 auto" xs10>
          <v-autocomplete
            slot="extension"
            v-model="select"
            :items="items"
            :search-input.sync="search"
            :loading="loading"
            label="Entreprises"
            placeholder="Siret, Raison Sociale..."
            prepend-icon="mdi-database-search"
            cache-items
            class="mx-3"
            solo
            style="border-radius: 25px"
            hide-no-data
            hide-details
          ></v-autocomplete>
        </v-flex>
      </v-layout>
    </v-container>
    <Etablissement v-if="select" :siret="select" batch="1802"></Etablissement>
  </div>
</template>

<script>
import Etablissement from '@/components/Etablissement'

export default {
  components: { Etablissement },
  data () {
    return {
      loading: false,
      items: [],
      search: null,
      select: null
    }
  },
  watch: {
    search (val) {
      val && val !== this.select && this.querySelections(val)
    }
  },
  methods: {
    querySelections (val) {
      this.loading = true
      this.$axios.post('/api/data/search', { 'text': val }).then(r => {
        this.items = r.data.map(e => { return { text: e._id.key + ' ' + e.value.sirene.raison_sociale, value: e._id.key } })
      }).finally(this.loading = false)
    }
  },
  computed: {
    message () {
      return this.$store.getters.reverseLog
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
    }
  }
}
</script>

<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
div.titre {
  font-family: 'Abel', sans-serif;
  color: #ffffff;
  font-weight: 800;
  font-size: 20px
}
span.fblue {
  font-family: 'Quicksand', sans-serif;
  color: #20459a
}
span.fblack {
  font-family: 'Quicksand', sans-serif;
  color: #000000
}
span.fred {
  font-family: 'Quicksand', sans-serif;
  color: #e9222e
}
.small {
  display:none;
  font-size: 15px;
  font-weight: 400;
}
.big {
  margin: 70px;
  text-shadow: 0 0 1px #00000040;
  font-size: 65px;
  font-weight: 800;
}
</style>
