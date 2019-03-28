<template>
  <v-navigation-drawer class="elevation-6" v-model="drawer" app>
    <v-toolbar flat class="transparent">
      <v-list class="pa-0">
        <v-list-tile avatar>
          <v-list-tile-avatar>
            <img src="/static/logo_signaux_faibles_cercle.svg">
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-title><span class="fblue">Signaux</span>·<span class="fred">Faibles</span></v-list-tile-title>
          </v-list-tile-content>
          <v-list-tile-avatar>
            <v-icon @click="drawer=!drawer">mdi-backburger</v-icon>
          </v-list-tile-avatar>
        </v-list-tile>
      </v-list>
    </v-toolbar>
     <v-list class="pt-0" dense>
        <v-divider></v-divider>
        <!-- <v-list-tile to="/">
          <v-list-tile-action>
            <v-icon>mdi-view-dashboard</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Tableau de Bord</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile> -->

        <v-list-tile to="/detection">
          <v-list-tile-action>
            <v-icon>mdi-target</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Détection</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-list-tile to="/">
          <v-list-tile-action>
            <v-icon>mdi-database-search</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Consultation</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-list-tile to="/data">
          <v-list-tile-action>
            <v-icon>mdi-database</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Données</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <!-- <v-list-tile to="/admin">
          <v-list-tile-action>
            <v-icon>fa-cog</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Administration</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile> -->

        <v-divider></v-divider>

        <v-list-tile @click="logout()">
          <v-list-tile-action>
            <v-icon>logout</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Se déconnecter</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    <v-footer class="elevation-6" style="color: blue; width:100%; position: fixed; bottom: 0px;">
      <v-btn
        flat
        icon
        color="blue"
        href="https://github.com/entrepreneur-interet-general/opensignauxfaibles">
        <v-icon>fab fa-github</v-icon>
      </v-btn>
    </v-footer>
  </v-navigation-drawer>

</template>

<script>
export default {
  methods: {
    logout () {
      this.$store.commit('logout')
    }
  },
  computed: {
    drawer: {
      get () {
        return this.$store.state.appDrawer
      },
      set (val) {
        this.$store.dispatch('setDrawer', val)
      }
    }
  }
}
</script>

<style scoped>
  span.fblue {
    font-family: 'Quicksand', sans-serif;
    color: #20459a
  }
  span.fred {
    font-family: 'Quicksand', sans-serif;
    color: #e9222e
  }
</style>
