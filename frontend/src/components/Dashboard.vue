<template>
  <div>
    <v-toolbar class="toolbar elevation-12" color="#c9aec5" height="35px" app>
      <v-icon
        @click="drawer=!drawer"
        class="fa-rotate-180"
        v-if="!drawer"
        color="#ffffff"
        key="toolbar"
        >mdi-backburger
      </v-icon>
      <div style="width: 100%; text-align: center;" class="titre">
        Tableau de Bord
      </div>
      <v-spacer></v-spacer>
      <v-icon color="#c9aec5" v-if="!rightDrawer" @click="rightDrawer=!rightDrawer">fa-dashboard</v-icon>
    </v-toolbar>
    <DashboardTask v-for="task in tasks" :key="task._id.siret" :task="task"></DashboardTask>
    <v-navigation-drawer
      :class="rightDrawer?'elevation-6':''"
      right app
      v-model="rightDrawer"
    >
      <v-list  two-line class="pt-0">
        <v-toolbar>
          <v-icon @click="rightDrawer=!rightDrawer" color="c9aec5">fa-dashbord</v-icon>
          <v-spacer></v-spacer>
            Suivi d'activité
          <v-divider></v-divider>
        </v-toolbar>
      </v-list>
    </v-navigation-drawer>
  </div>
</template>

<script>
import DashboardTask from '@/components/DashboardTask'

export default {
  components: { DashboardTask },
  data () {
    return {
      tasks: []
    }
  },
  watch: {
    search (val) {
      val && val !== this.select && this.querySelections(val)
    }
  },
  methods: {
    getTasks () {
      this.$axios.get('/api/dashboard/tasks').then(response => {
        this.tasks = response.data
      })
    }
  },
  mounted () {
    this.getTasks()
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
  font-family: 'Abel', cursive;
  color: #fff;
  font-weight: 500;
  font-size: 18px
}
</style>
