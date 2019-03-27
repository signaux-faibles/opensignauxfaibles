<template>
<div>
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
    <v-icon color="#ffffff" @click="rightDrawer=!rightDrawer">mdi-database-search</v-icon>
  </v-toolbar>
    <v-container>
      <v-layout>
        <v-flex>
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
            flat
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
      this.$axios.post('/api/data/search', {'text': val}).then(r => {
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
</style>
