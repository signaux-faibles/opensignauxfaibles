<template>
  <v-card class="ma-3 pa-1 task elevation-12">
    <v-card-text class="tasktitle">
      <v-icon style="position: absolute; top: 30px; right: 30px" dark>mdi-eye</v-icon>
      <v-icon 
        style="position: absolute; bottom: 30px; right: 30px" 
        @click="viewDetail=!viewDetail"
        dark>mdi-arrow-{{ viewDetail?'up':'down' }}-bold-box</v-icon>
      <div>
        <div style="vertical-align: 'top'">
          <v-card-title class="title">{{ task.etablissement[0].value.sirene.raisonsociale }}<br/></v-card-title>
        </div>
        <span>SIRET {{ task._id }}</span>
        <div style="vertical-align: 'top'">
          Création: {{ formatDate(task.firstDate) }}<br/>
          Dernier ajout: {{ formatDate(task.lastDate) }}
        </div>
      </div>
    </v-card-text>
    <v-divider></v-divider>
    <div v-if="viewDetail">
    <v-card-text 
      class="tasktext ma-2" 
      dark 
      v-for="s in task.tasks" 
      :key="s.id" 
      style="background-color: transparent;"
    >
        {{ s.event }} le {{ formatDate(s._id.date) }}<br/>
        {{ s.value }}<br/>
    </v-card-text>
    </div>
  </v-card>
</template>

<script>
export default {
  props: ['task'],
  methods: {
    formatDate (dateString) {
      var d = new Date(dateString)
      return ('0' + d.getDate()).slice(-2) + '/' +
        ('0' + d.getMonth()).slice(-2) + '/' + d.getFullYear()
    }
  },
  data () {
    return {
      viewDetail: false
    }
  }
}
</script>
